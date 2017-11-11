package scenarios

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/freedom"
	v2http "v2ray.com/core/proxy/http"
	v2httptest "v2ray.com/core/testing/servers/http"
	. "v2ray.com/ext/assert"
)

func TestHttpConformance(t *testing.T) {
	assert := With(t)

	httpServerPort := pickPort()
	httpServer := &v2httptest.Server{
		Port:        httpServerPort,
		PathHandler: make(map[string]http.HandlerFunc),
	}
	_, err := httpServer.Start()
	assert(err, IsNil)
	defer httpServer.Close()

	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&v2http.ServerConfig{}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig)
	assert(err, IsNil)

	{
		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:" + serverPort.String())
			},
		}

		client := &http.Client{
			Transport: transport,
		}

		resp, err := client.Get("http://127.0.0.1:" + httpServerPort.String())
		assert(err, IsNil)
		assert(resp.StatusCode, Equals, 200)

		content, err := ioutil.ReadAll(resp.Body)
		assert(err, IsNil)
		assert(string(content), Equals, "Home")

	}

	CloseAllServers(servers)
}

func setProxyBasicAuth(req *http.Request, user, pass string) {
	req.SetBasicAuth(user, pass)
	req.Header.Set("Proxy-Authorization", req.Header.Get("Authorization"))
	req.Header.Del("Authorization")
}

func TestHttpBasicAuth(t *testing.T) {
	assert := With(t)

	httpServerPort := pickPort()
	httpServer := &v2httptest.Server{
		Port:        httpServerPort,
		PathHandler: make(map[string]http.HandlerFunc),
	}
	_, err := httpServer.Start()
	assert(err, IsNil)
	defer httpServer.Close()

	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&v2http.ServerConfig{
					Accounts: map[string]string{
						"a": "b",
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

	servers, err := InitializeServerConfigs(serverConfig)
	assert(err, IsNil)

	{
		transport := &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:" + serverPort.String())
			},
		}

		client := &http.Client{
			Transport: transport,
		}

		{
			resp, err := client.Get("http://127.0.0.1:" + httpServerPort.String())
			assert(err, IsNil)
			assert(resp.StatusCode, Equals, 401)
		}

		{
			req, err := http.NewRequest("GET", "http://127.0.0.1:"+httpServerPort.String(), nil)
			assert(err, IsNil)

			setProxyBasicAuth(req, "a", "c")
			resp, err := client.Do(req)
			assert(err, IsNil)
			assert(resp.StatusCode, Equals, 401)
		}

		{
			req, err := http.NewRequest("GET", "http://127.0.0.1:"+httpServerPort.String(), nil)
			assert(err, IsNil)

			setProxyBasicAuth(req, "a", "b")
			resp, err := client.Do(req)
			assert(err, IsNil)
			assert(resp.StatusCode, Equals, 200)

			content, err := ioutil.ReadAll(resp.Body)
			assert(err, IsNil)
			assert(string(content), Equals, "Home")
		}
	}

	CloseAllServers(servers)
}
