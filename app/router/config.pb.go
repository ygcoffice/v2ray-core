package router

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net1 "v2ray.com/core/common/net"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Type of domain value.
type Domain_Type int32

const (
	// The value is used as is.
	Domain_Plain Domain_Type = 0
	// The value is used as a regular expression.
	Domain_Regex Domain_Type = 1
	// The value is a domain.
	Domain_Domain Domain_Type = 2
)

var Domain_Type_name = map[int32]string{
	0: "Plain",
	1: "Regex",
	2: "Domain",
}
var Domain_Type_value = map[string]int32{
	"Plain":  0,
	"Regex":  1,
	"Domain": 2,
}

func (x Domain_Type) String() string {
	return proto.EnumName(Domain_Type_name, int32(x))
}
func (Domain_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Config_DomainStrategy int32

const (
	// Use domain as is.
	Config_AsIs Config_DomainStrategy = 0
	// Always resolve IP for domains.
	Config_UseIp Config_DomainStrategy = 1
	// Resolve to IP if the domain doesn't match any rules.
	Config_IpIfNonMatch Config_DomainStrategy = 2
)

var Config_DomainStrategy_name = map[int32]string{
	0: "AsIs",
	1: "UseIp",
	2: "IpIfNonMatch",
}
var Config_DomainStrategy_value = map[string]int32{
	"AsIs":         0,
	"UseIp":        1,
	"IpIfNonMatch": 2,
}

func (x Config_DomainStrategy) String() string {
	return proto.EnumName(Config_DomainStrategy_name, int32(x))
}
func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

// Domain for routing decision.
type Domain struct {
	// Domain matching type.
	Type Domain_Type `protobuf:"varint,1,opt,name=type,enum=v2ray.core.app.router.Domain_Type" json:"type,omitempty"`
	// Domain value.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Domain) Reset()                    { *m = Domain{} }
func (m *Domain) String() string            { return proto.CompactTextString(m) }
func (*Domain) ProtoMessage()               {}
func (*Domain) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Domain) GetType() Domain_Type {
	if m != nil {
		return m.Type
	}
	return Domain_Plain
}

func (m *Domain) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// IP for routing decision, in CIDR form.
type CIDR struct {
	// IP address, should be either 4 or 16 bytes.
	Ip []byte `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	// Number of leading ones in the network mask.
	Prefix uint32 `protobuf:"varint,2,opt,name=prefix" json:"prefix,omitempty"`
}

func (m *CIDR) Reset()                    { *m = CIDR{} }
func (m *CIDR) String() string            { return proto.CompactTextString(m) }
func (*CIDR) ProtoMessage()               {}
func (*CIDR) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CIDR) GetIp() []byte {
	if m != nil {
		return m.Ip
	}
	return nil
}

func (m *CIDR) GetPrefix() uint32 {
	if m != nil {
		return m.Prefix
	}
	return 0
}

type GeoIP struct {
	CountryCode string  `protobuf:"bytes,1,opt,name=country_code,json=countryCode" json:"country_code,omitempty"`
	Cidr        []*CIDR `protobuf:"bytes,2,rep,name=cidr" json:"cidr,omitempty"`
}

func (m *GeoIP) Reset()                    { *m = GeoIP{} }
func (m *GeoIP) String() string            { return proto.CompactTextString(m) }
func (*GeoIP) ProtoMessage()               {}
func (*GeoIP) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GeoIP) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *GeoIP) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
	}
	return nil
}

type GeoIPList struct {
	Entry []*GeoIP `protobuf:"bytes,1,rep,name=entry" json:"entry,omitempty"`
}

func (m *GeoIPList) Reset()                    { *m = GeoIPList{} }
func (m *GeoIPList) String() string            { return proto.CompactTextString(m) }
func (*GeoIPList) ProtoMessage()               {}
func (*GeoIPList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GeoIPList) GetEntry() []*GeoIP {
	if m != nil {
		return m.Entry
	}
	return nil
}

type GeoSite struct {
	CountryCode string    `protobuf:"bytes,1,opt,name=country_code,json=countryCode" json:"country_code,omitempty"`
	Domain      []*Domain `protobuf:"bytes,2,rep,name=domain" json:"domain,omitempty"`
}

func (m *GeoSite) Reset()                    { *m = GeoSite{} }
func (m *GeoSite) String() string            { return proto.CompactTextString(m) }
func (*GeoSite) ProtoMessage()               {}
func (*GeoSite) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GeoSite) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *GeoSite) GetDomain() []*Domain {
	if m != nil {
		return m.Domain
	}
	return nil
}

type GeoSiteList struct {
	Entry []*GeoSite `protobuf:"bytes,1,rep,name=entry" json:"entry,omitempty"`
}

func (m *GeoSiteList) Reset()                    { *m = GeoSiteList{} }
func (m *GeoSiteList) String() string            { return proto.CompactTextString(m) }
func (*GeoSiteList) ProtoMessage()               {}
func (*GeoSiteList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GeoSiteList) GetEntry() []*GeoSite {
	if m != nil {
		return m.Entry
	}
	return nil
}

type RoutingRule struct {
	Tag         string                              `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	Domain      []*Domain                           `protobuf:"bytes,2,rep,name=domain" json:"domain,omitempty"`
	Cidr        []*CIDR                             `protobuf:"bytes,3,rep,name=cidr" json:"cidr,omitempty"`
	PortRange   *v2ray_core_common_net.PortRange    `protobuf:"bytes,4,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	NetworkList *v2ray_core_common_net1.NetworkList `protobuf:"bytes,5,opt,name=network_list,json=networkList" json:"network_list,omitempty"`
	SourceCidr  []*CIDR                             `protobuf:"bytes,6,rep,name=source_cidr,json=sourceCidr" json:"source_cidr,omitempty"`
	UserEmail   []string                            `protobuf:"bytes,7,rep,name=user_email,json=userEmail" json:"user_email,omitempty"`
	InboundTag  []string                            `protobuf:"bytes,8,rep,name=inbound_tag,json=inboundTag" json:"inbound_tag,omitempty"`
}

func (m *RoutingRule) Reset()                    { *m = RoutingRule{} }
func (m *RoutingRule) String() string            { return proto.CompactTextString(m) }
func (*RoutingRule) ProtoMessage()               {}
func (*RoutingRule) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RoutingRule) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *RoutingRule) GetDomain() []*Domain {
	if m != nil {
		return m.Domain
	}
	return nil
}

func (m *RoutingRule) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
	}
	return nil
}

func (m *RoutingRule) GetPortRange() *v2ray_core_common_net.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *RoutingRule) GetNetworkList() *v2ray_core_common_net1.NetworkList {
	if m != nil {
		return m.NetworkList
	}
	return nil
}

func (m *RoutingRule) GetSourceCidr() []*CIDR {
	if m != nil {
		return m.SourceCidr
	}
	return nil
}

func (m *RoutingRule) GetUserEmail() []string {
	if m != nil {
		return m.UserEmail
	}
	return nil
}

func (m *RoutingRule) GetInboundTag() []string {
	if m != nil {
		return m.InboundTag
	}
	return nil
}

type Config struct {
	DomainStrategy Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,enum=v2ray.core.app.router.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Rule           []*RoutingRule        `protobuf:"bytes,2,rep,name=rule" json:"rule,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Config) GetDomainStrategy() Config_DomainStrategy {
	if m != nil {
		return m.DomainStrategy
	}
	return Config_AsIs
}

func (m *Config) GetRule() []*RoutingRule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func init() {
	proto.RegisterType((*Domain)(nil), "v2ray.core.app.router.Domain")
	proto.RegisterType((*CIDR)(nil), "v2ray.core.app.router.CIDR")
	proto.RegisterType((*GeoIP)(nil), "v2ray.core.app.router.GeoIP")
	proto.RegisterType((*GeoIPList)(nil), "v2ray.core.app.router.GeoIPList")
	proto.RegisterType((*GeoSite)(nil), "v2ray.core.app.router.GeoSite")
	proto.RegisterType((*GeoSiteList)(nil), "v2ray.core.app.router.GeoSiteList")
	proto.RegisterType((*RoutingRule)(nil), "v2ray.core.app.router.RoutingRule")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.router.Config")
	proto.RegisterEnum("v2ray.core.app.router.Domain_Type", Domain_Type_name, Domain_Type_value)
	proto.RegisterEnum("v2ray.core.app.router.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/router/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 626 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0x49, 0xda, 0x66, 0xcb, 0x49, 0x29, 0x91, 0xc5, 0x50, 0x18, 0x0c, 0x4a, 0x84, 0xa0,
	0x17, 0x28, 0x91, 0xca, 0xbf, 0x1b, 0xd0, 0x34, 0xba, 0x69, 0xaa, 0x04, 0x53, 0xe5, 0x6d, 0x5c,
	0xc0, 0x45, 0x94, 0xa5, 0x5e, 0x88, 0x68, 0x6d, 0xcb, 0x71, 0xc6, 0x7a, 0xc7, 0x0b, 0xf0, 0x22,
	0x3c, 0x0d, 0x8f, 0x84, 0x6c, 0xa7, 0xb0, 0xa2, 0x05, 0x26, 0xee, 0x6c, 0xe7, 0xf7, 0x9d, 0xf3,
	0xe5, 0xf8, 0x1c, 0xc3, 0xa3, 0xb3, 0xa1, 0x48, 0x17, 0x51, 0xc6, 0xe6, 0x71, 0xc6, 0x04, 0x89,
	0x53, 0xce, 0x63, 0xc1, 0x2a, 0x49, 0x44, 0x9c, 0x31, 0x7a, 0x5a, 0xe4, 0x11, 0x17, 0x4c, 0x32,
	0xb4, 0xb1, 0xe4, 0x04, 0x89, 0x52, 0xce, 0x23, 0xc3, 0x6c, 0x3e, 0xfc, 0x43, 0x9e, 0xb1, 0xf9,
	0x9c, 0xd1, 0x98, 0x12, 0x19, 0x73, 0x26, 0xa4, 0x11, 0x6f, 0x3e, 0x6e, 0xa6, 0x28, 0x91, 0x5f,
	0x98, 0xf8, 0x6c, 0xc0, 0xf0, 0xab, 0x05, 0xce, 0x2e, 0x9b, 0xa7, 0x05, 0x45, 0x2f, 0xa0, 0x2d,
	0x17, 0x9c, 0x04, 0x56, 0xdf, 0x1a, 0xf4, 0x86, 0x61, 0x74, 0x69, 0xfe, 0xc8, 0xc0, 0xd1, 0xd1,
	0x82, 0x13, 0xac, 0x79, 0x74, 0x13, 0x3a, 0x67, 0xe9, 0xac, 0x22, 0x81, 0xdd, 0xb7, 0x06, 0x2e,
	0x36, 0x9b, 0x70, 0x00, 0x6d, 0xc5, 0x20, 0x17, 0x3a, 0x93, 0x59, 0x5a, 0x50, 0xff, 0x9a, 0x5a,
	0x62, 0x92, 0x93, 0x73, 0xdf, 0x42, 0xb0, 0xcc, 0xea, 0xdb, 0x61, 0x04, 0xed, 0xd1, 0x78, 0x17,
	0xa3, 0x1e, 0xd8, 0x05, 0xd7, 0xd9, 0xbb, 0xd8, 0x2e, 0x38, 0xba, 0x05, 0x0e, 0x17, 0xe4, 0xb4,
	0x38, 0xd7, 0x81, 0xaf, 0xe3, 0x7a, 0x17, 0x7e, 0x84, 0xce, 0x3e, 0x61, 0xe3, 0x09, 0x7a, 0x00,
	0xdd, 0x8c, 0x55, 0x54, 0x8a, 0x45, 0x92, 0xb1, 0xa9, 0x31, 0xee, 0x62, 0xaf, 0x3e, 0x1b, 0xb1,
	0x29, 0x41, 0x31, 0xb4, 0xb3, 0x62, 0x2a, 0x02, 0xbb, 0xdf, 0x1a, 0x78, 0xc3, 0x3b, 0x0d, 0xff,
	0xa4, 0xd2, 0x63, 0x0d, 0x86, 0xdb, 0xe0, 0xea, 0xe0, 0x6f, 0x8b, 0x52, 0xa2, 0x21, 0x74, 0x88,
	0x0a, 0x15, 0x58, 0x5a, 0x7e, 0xb7, 0x41, 0xae, 0x05, 0xd8, 0xa0, 0x61, 0x06, 0x6b, 0xfb, 0x84,
	0x1d, 0x16, 0x92, 0x5c, 0xc5, 0xdf, 0x73, 0x70, 0xa6, 0xba, 0x0e, 0xb5, 0xc3, 0xad, 0xbf, 0x56,
	0x1d, 0xd7, 0x70, 0x38, 0x02, 0xaf, 0x4e, 0xa2, 0x7d, 0x3e, 0x5b, 0xf5, 0x79, 0xaf, 0xd9, 0xa7,
	0x92, 0x2c, 0x9d, 0x7e, 0x6b, 0x81, 0x87, 0x59, 0x25, 0x0b, 0x9a, 0xe3, 0x6a, 0x46, 0x90, 0x0f,
	0x2d, 0x99, 0xe6, 0xb5, 0x4b, 0xb5, 0xfc, 0x4f, 0x77, 0xbf, 0x8a, 0xde, 0xba, 0x62, 0xd1, 0xd1,
	0x36, 0x80, 0xea, 0xdd, 0x44, 0xa4, 0x34, 0x27, 0x41, 0xbb, 0x6f, 0x0d, 0xbc, 0x61, 0xff, 0xa2,
	0xcc, 0xb4, 0x6f, 0x44, 0x89, 0x8c, 0x26, 0x4c, 0x48, 0xac, 0x38, 0xec, 0xf2, 0xe5, 0x12, 0xed,
	0x41, 0xb7, 0x6e, 0xeb, 0x64, 0x56, 0x94, 0x32, 0xe8, 0xe8, 0x10, 0x61, 0x43, 0x88, 0x03, 0x83,
	0xaa, 0xd2, 0x61, 0x8f, 0xfe, 0xde, 0xa0, 0x57, 0xe0, 0x95, 0xac, 0x12, 0x19, 0x49, 0xb4, 0x7f,
	0xe7, 0xdf, 0xfe, 0xc1, 0xf0, 0x23, 0xf5, 0x17, 0x5b, 0x00, 0x55, 0x49, 0x44, 0x42, 0xe6, 0x69,
	0x31, 0x0b, 0xd6, 0xfa, 0xad, 0x81, 0x8b, 0x5d, 0x75, 0xb2, 0xa7, 0x0e, 0xd0, 0x7d, 0xf0, 0x0a,
	0x7a, 0xc2, 0x2a, 0x3a, 0x4d, 0x54, 0x99, 0xd7, 0xf5, 0x77, 0xa8, 0x8f, 0x8e, 0xd2, 0x3c, 0xfc,
	0x61, 0x81, 0x33, 0xd2, 0x2f, 0x00, 0x3a, 0x86, 0x1b, 0xa6, 0x96, 0x49, 0x29, 0x45, 0x2a, 0x49,
	0xbe, 0xa8, 0xa7, 0xf2, 0x49, 0x93, 0x19, 0xf3, 0x72, 0x98, 0x8b, 0x38, 0xac, 0x35, 0xb8, 0x37,
	0x5d, 0xd9, 0xab, 0x09, 0x17, 0xd5, 0x8c, 0xd4, 0xb7, 0xd9, 0x34, 0xe1, 0x17, 0x7a, 0x02, 0x6b,
	0x3e, 0x7c, 0x09, 0xbd, 0xd5, 0xc8, 0x68, 0x1d, 0xda, 0x3b, 0xe5, 0xb8, 0x34, 0x43, 0x7d, 0x5c,
	0x92, 0x31, 0xf7, 0x2d, 0xe4, 0x43, 0x77, 0xcc, 0xc7, 0xa7, 0x07, 0x8c, 0xbe, 0x4b, 0x65, 0xf6,
	0xc9, 0xb7, 0xdf, 0xbc, 0x86, 0xdb, 0x19, 0x9b, 0x5f, 0x9e, 0x67, 0x62, 0x7d, 0x70, 0xcc, 0xea,
	0xbb, 0xbd, 0xf1, 0x7e, 0x88, 0xd3, 0x45, 0x34, 0x52, 0xc4, 0x0e, 0xe7, 0xda, 0x02, 0x11, 0x27,
	0x8e, 0x7e, 0xa3, 0x9e, 0xfe, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x6d, 0xa3, 0xd9, 0xd0, 0x33, 0x05,
	0x00, 0x00,
}
