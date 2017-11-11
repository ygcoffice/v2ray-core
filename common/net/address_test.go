package net_test

import (
	"net"
	"testing"

	. "v2ray.com/core/common/net"
	. "v2ray.com/core/common/net/testing"
	. "v2ray.com/ext/assert"
)

func TestIPv4Address(t *testing.T) {
	assert := With(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	addr := IPAddress(ip)

	assert(addr, IsIPv4)
	assert(addr, Not(IsIPv6))
	assert(addr, Not(IsDomain))
	assert([]byte(addr.IP()), Equals, ip)
	assert(addr.String(), Equals, "1.2.3.4")
}

func TestIPv6Address(t *testing.T) {
	assert := With(t)

	ip := []byte{
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := IPAddress(ip)

	assert(addr, IsIPv6)
	assert(addr, Not(IsIPv4))
	assert(addr, Not(IsDomain))
	assert(addr.IP(), Equals, net.IP(ip))
	assert(addr.String(), Equals, "[102:304:102:304:102:304:102:304]")
}

func TestIPv4Asv6(t *testing.T) {
	assert := With(t)
	ip := []byte{
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(255), byte(255),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := IPAddress(ip)
	assert(addr.String(), Equals, "1.2.3.4")
}

func TestDomainAddress(t *testing.T) {
	assert := With(t)

	domain := "v2ray.com"
	addr := DomainAddress(domain)

	assert(addr, IsDomain)
	assert(addr, Not(IsIPv6))
	assert(addr, Not(IsIPv4))
	assert(addr.Domain(), Equals, domain)
	assert(addr.String(), Equals, "v2ray.com")
}

func TestNetIPv4Address(t *testing.T) {
	assert := With(t)

	ip := net.IPv4(1, 2, 3, 4)
	addr := IPAddress(ip)
	assert(addr, IsIPv4)
	assert(addr.String(), Equals, "1.2.3.4")
}
