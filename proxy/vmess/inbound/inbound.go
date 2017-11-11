package inbound

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg inbound -path Proxy,VMess,Inbound

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type userByEmail struct {
	sync.RWMutex
	cache           map[string]*protocol.User
	defaultLevel    uint32
	defaultAlterIDs uint16
}

func NewUserByEmail(users []*protocol.User, config *DefaultConfig) *userByEmail {
	cache := make(map[string]*protocol.User)
	for _, user := range users {
		cache[user.Email] = user
	}
	return &userByEmail{
		cache:           cache,
		defaultLevel:    config.Level,
		defaultAlterIDs: uint16(config.AlterId),
	}
}

func (v *userByEmail) Get(email string) (*protocol.User, bool) {
	var user *protocol.User
	var found bool
	v.RLock()
	user, found = v.cache[email]
	v.RUnlock()
	if !found {
		v.Lock()
		user, found = v.cache[email]
		if !found {
			account := &vmess.Account{
				Id:      uuid.New().String(),
				AlterId: uint32(v.defaultAlterIDs),
			}
			user = &protocol.User{
				Level:   v.defaultLevel,
				Email:   email,
				Account: serial.ToTypedMessage(account),
			}
			v.cache[email] = user
		}
		v.Unlock()
	}
	return user, found
}

// Handler is an inbound connection handler that handles messages in VMess protocol.
type Handler struct {
	inboundHandlerManager proxyman.InboundHandlerManager
	clients               protocol.UserValidator
	usersByEmail          *userByEmail
	detours               *DetourConfig
	sessionHistory        *encoding.SessionHistory
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}

	allowedClients := vmess.NewTimedUserValidator(ctx, protocol.DefaultIDHash)
	for _, user := range config.User {
		if err := allowedClients.Add(user); err != nil {
			return nil, newError("failed to initiate user").Base(err)
		}
	}

	handler := &Handler{
		clients:        allowedClients,
		detours:        config.Detour,
		usersByEmail:   NewUserByEmail(config.User, config.GetDefaultValue()),
		sessionHistory: encoding.NewSessionHistory(ctx),
	}

	space.OnInitialize(func() error {
		handler.inboundHandlerManager = proxyman.InboundHandlerManagerFromSpace(space)
		if handler.inboundHandlerManager == nil {
			return newError("InboundHandlerManager is not found is space")
		}
		return nil
	})

	return handler, nil
}

// Network implements proxy.Inbound.Network().
func (*Handler) Network() net.NetworkList {
	return net.NetworkList{
		Network: []net.Network{net.Network_TCP},
	}
}

func (v *Handler) GetUser(email string) *protocol.User {
	user, existing := v.usersByEmail.Get(email)
	if !existing {
		v.clients.Add(user)
	}
	return user
}

func transferRequest(timer signal.ActivityUpdater, session *encoding.ServerSession, request *protocol.RequestHeader, input io.Reader, output ray.OutputStream) error {
	defer output.Close()

	bodyReader := session.DecodeRequestBody(request, input)
	if err := buf.Copy(bodyReader, output, buf.UpdateActivity(timer)); err != nil {
		return newError("failed to transfer request").Base(err)
	}
	return nil
}

func transferResponse(timer signal.ActivityUpdater, session *encoding.ServerSession, request *protocol.RequestHeader, response *protocol.ResponseHeader, input buf.Reader, output io.Writer) error {
	session.EncodeResponseHeader(response, output)

	bodyWriter := session.EncodeResponseBody(request, output)

	// Optimize for small response packet
	data, err := input.ReadMultiBuffer()
	if err != nil {
		return err
	}

	if err := bodyWriter.WriteMultiBuffer(data); err != nil {
		return err
	}
	data.Release()

	if bufferedWriter, ok := output.(*buf.BufferedWriter); ok {
		if err := bufferedWriter.SetBuffered(false); err != nil {
			return err
		}
	}

	if err := buf.Copy(input, bodyWriter, buf.UpdateActivity(timer)); err != nil {
		return err
	}

	if request.Option.Has(protocol.RequestOptionChunkStream) {
		if err := bodyWriter.WriteMultiBuffer(buf.MultiBuffer{}); err != nil {
			return err
		}
	}

	return nil
}

// Process implements proxy.Inbound.Process().
func (v *Handler) Process(ctx context.Context, network net.Network, connection internet.Connection, dispatcher dispatcher.Interface) error {
	if err := connection.SetReadDeadline(time.Now().Add(time.Second * 8)); err != nil {
		return err
	}

	reader := buf.NewBufferedReader(buf.NewReader(connection))

	session := encoding.NewServerSession(v.clients, v.sessionHistory)
	request, err := session.DecodeRequestHeader(reader)

	if err != nil {
		if errors.Cause(err) != io.EOF {
			log.Access(connection.RemoteAddr(), "", log.AccessRejected, err)
			log.Trace(newError("invalid request from ", connection.RemoteAddr(), ": ", err).AtInfo())
		}
		return err
	}

	if request.Command == protocol.RequestCommandMux {
		request.Address = net.DomainAddress("v1.mux.com")
		request.Port = net.Port(0)
	}

	log.Access(connection.RemoteAddr(), request.Destination(), log.AccessAccepted, "")
	log.Trace(newError("received request for ", request.Destination()))

	common.Must(connection.SetReadDeadline(time.Time{}))

	userSettings := request.User.GetSettings()

	ctx = protocol.ContextWithUser(ctx, request.User)

	ctx, timer := signal.CancelAfterInactivity(ctx, userSettings.PayloadTimeout)
	ray, err := dispatcher.Dispatch(ctx, request.Destination())
	if err != nil {
		return newError("failed to dispatch request to ", request.Destination()).Base(err)
	}

	input := ray.InboundInput()
	output := ray.InboundOutput()

	requestDone := signal.ExecuteAsync(func() error {
		return transferRequest(timer, session, request, reader, input)
	})

	responseDone := signal.ExecuteAsync(func() error {
		writer := buf.NewBufferedWriter(buf.NewWriter(connection))
		defer writer.Flush()

		response := &protocol.ResponseHeader{
			Command: v.generateCommand(ctx, request),
		}
		return transferResponse(timer, session, request, response, output, writer)
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		input.CloseError()
		output.CloseError()
		return newError("connection ends").Base(err)
	}

	return nil
}

func (v *Handler) generateCommand(ctx context.Context, request *protocol.RequestHeader) protocol.ResponseCommand {
	if v.detours != nil {
		tag := v.detours.To
		if v.inboundHandlerManager != nil {
			handler, err := v.inboundHandlerManager.GetHandler(ctx, tag)
			if err != nil {
				log.Trace(newError("failed to get detour handler: ", tag, err).AtWarning())
				return nil
			}
			proxyHandler, port, availableMin := handler.GetRandomInboundProxy()
			inboundHandler, ok := proxyHandler.(*Handler)
			if ok && inboundHandler != nil {
				if availableMin > 255 {
					availableMin = 255
				}

				log.Trace(newError("pick detour handler for port ", port, " for ", availableMin, " minutes.").AtDebug())
				user := inboundHandler.GetUser(request.User.Email)
				if user == nil {
					return nil
				}
				account, _ := user.GetTypedAccount()
				return &protocol.CommandSwitchAccount{
					Port:     port,
					ID:       account.(*vmess.InternalAccount).ID.UUID(),
					AlterIds: uint16(len(account.(*vmess.InternalAccount).AlterIDs)),
					Level:    user.Level,
					ValidMin: byte(availableMin),
				}
			}
		}
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
