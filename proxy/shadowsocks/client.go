package shadowsocks

import (
	"context"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// Client is a inbound handler for Shadowsocks protocol
type Client struct {
	serverPicker protocol.ServerPicker
}

// NewClient create a new Shadowsocks client.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	if serverList.Size() == 0 {
		return nil, newError("0 server")
	}
	client := &Client{
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
	}

	return client, nil
}

// Process implements OutboundHandler.Process().
func (v *Client) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	destination, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return newError("target not specified")
	}
	network := destination.Network

	var server *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = v.serverPicker.PickServer()
		dest := server.Destination()
		dest.Network = network
		rawConn, err := dialer.Dial(ctx, dest)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return newError("failed to find an available destination").AtWarning().Base(err)
	}
	log.Trace(newError("tunneling request to ", destination, " via ", server.Destination()))

	defer conn.Close()

	request := &protocol.RequestHeader{
		Version: Version,
		Address: destination.Address,
		Port:    destination.Port,
	}
	if destination.Network == net.Network_TCP {
		request.Command = protocol.RequestCommandTCP
	} else {
		request.Command = protocol.RequestCommandUDP
	}

	user := server.PickUser()
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return newError("failed to get a valid user account").AtWarning().Base(err)
	}
	account := rawAccount.(*ShadowsocksAccount)
	request.User = user

	if account.OneTimeAuth == Account_Auto || account.OneTimeAuth == Account_Enabled {
		request.Option |= RequestOptionOneTimeAuth
	}

	ctx, timer := signal.CancelAfterInactivity(ctx, time.Minute*5)

	if request.Command == protocol.RequestCommandTCP {
		bufferedWriter := buf.NewBufferedWriter(buf.NewWriter(conn))
		bodyWriter, err := WriteTCPRequest(request, bufferedWriter)
		if err != nil {
			return newError("failed to write request").Base(err)
		}

		if err := bufferedWriter.SetBuffered(false); err != nil {
			return err
		}

		requestDone := signal.ExecuteAsync(func() error {
			if err := buf.Copy(outboundRay.OutboundInput(), bodyWriter, buf.UpdateActivity(timer)); err != nil {
				return err
			}
			return nil
		})

		responseDone := signal.ExecuteAsync(func() error {
			defer outboundRay.OutboundOutput().Close()

			responseReader, err := ReadTCPResponse(user, conn)
			if err != nil {
				return err
			}

			if err := buf.Copy(responseReader, outboundRay.OutboundOutput(), buf.UpdateActivity(timer)); err != nil {
				return err
			}

			return nil
		})

		if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
			return newError("connection ends").Base(err)
		}

		return nil
	}

	if request.Command == protocol.RequestCommandUDP {

		writer := buf.NewSequentialWriter(&UDPWriter{
			Writer:  conn,
			Request: request,
		})

		requestDone := signal.ExecuteAsync(func() error {
			if err := buf.Copy(outboundRay.OutboundInput(), writer, buf.UpdateActivity(timer)); err != nil {
				return newError("failed to transport all UDP request").Base(err)
			}
			return nil
		})

		responseDone := signal.ExecuteAsync(func() error {
			defer outboundRay.OutboundOutput().Close()

			reader := &UDPReader{
				Reader: conn,
				User:   user,
			}

			if err := buf.Copy(reader, outboundRay.OutboundOutput(), buf.UpdateActivity(timer)); err != nil {
				return newError("failed to transport all UDP response").Base(err)
			}
			return nil
		})

		if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
			return newError("connection ends").Base(err)
		}

		return nil
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
