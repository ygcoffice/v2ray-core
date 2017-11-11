package scenarios

import (
	"testing"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/proxy/vmess/outbound"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/http"
	tcptransport "v2ray.com/core/transport/internet/tcp"
	. "v2ray.com/ext/assert"
)

func TestHttpConnectionHeader(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
					StreamSettings: &internet.StreamConfig{
						TransportSettings: []*internet.TransportConfig{
							{
								Protocol: internet.TransportProtocol_TCP,
								Settings: serial.ToTypedMessage(&tcptransport.Config{
									HeaderSettings: serial.ToTypedMessage(&http.Config{}),
								}),
							},
						},
					},
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id: userID.String(),
							}),
						},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(clientPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: userID.String(),
									}),
								},
							},
						},
					},
				}),
				SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
					StreamSettings: &internet.StreamConfig{
						TransportSettings: []*internet.TransportConfig{
							{
								Protocol: internet.TransportProtocol_TCP,
								Settings: serial.ToTypedMessage(&tcptransport.Config{
									HeaderSettings: serial.ToTypedMessage(&http.Config{}),
								}),
							},
						},
					},
				}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig, clientConfig)
	assert(err, IsNil)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert(err, IsNil)

	payload := "dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert(err, IsNil)
	assert(nBytes, Equals, len(payload))

	response := readFrom(conn, time.Second*2, len(payload))
	assert(response, Equals, xor([]byte(payload)))
	assert(conn.Close(), IsNil)

	CloseAllServers(servers)
}
