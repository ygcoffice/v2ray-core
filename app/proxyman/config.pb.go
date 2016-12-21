// Code generated by protoc-gen-go.
// source: v2ray.com/core/app/proxyman/config.proto
// DO NOT EDIT!

/*
Package proxyman is a generated protocol buffer package.

It is generated from these files:
	v2ray.com/core/app/proxyman/config.proto

It has these top-level messages:
	InboundConfig
	OutboundConfig
*/
package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type InboundConfig struct {
}

func (m *InboundConfig) Reset()                    { *m = InboundConfig{} }
func (m *InboundConfig) String() string            { return proto.CompactTextString(m) }
func (*InboundConfig) ProtoMessage()               {}
func (*InboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type OutboundConfig struct {
}

func (m *OutboundConfig) Reset()                    { *m = OutboundConfig{} }
func (m *OutboundConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundConfig) ProtoMessage()               {}
func (*OutboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*InboundConfig)(nil), "v2ray.core.app.proxyman.InboundConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 129 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xd2, 0x28, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0x2c, 0x28, 0xd0, 0x2f,
	0x28, 0xca, 0xaf, 0xa8, 0xcc, 0x4d, 0xcc, 0xd3, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x87, 0xa9, 0x2c, 0x4a, 0xd5, 0x4b, 0x2c, 0x28, 0xd0, 0x83,
	0xa9, 0x52, 0xe2, 0xe7, 0xe2, 0xf5, 0xcc, 0x4b, 0xca, 0x2f, 0xcd, 0x4b, 0x71, 0x06, 0xab, 0x57,
	0x12, 0xe0, 0xe2, 0xf3, 0x2f, 0x2d, 0x41, 0x12, 0x71, 0x32, 0xe1, 0x92, 0x4e, 0xce, 0xcf, 0xd5,
	0xc3, 0x61, 0x82, 0x13, 0x37, 0x44, 0x59, 0x00, 0xc8, 0x9e, 0x28, 0x0e, 0x98, 0x70, 0x12, 0x1b,
	0xd8, 0x62, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x41, 0x39, 0xa0, 0x25, 0xa4, 0x00, 0x00,
	0x00,
}