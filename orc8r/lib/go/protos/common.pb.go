// Code generated by protoc-gen-go. DO NOT EDIT.
// source: orc8r/protos/common.proto

package protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// --------------------------------------------------------------------------
// Logging levels
// --------------------------------------------------------------------------
type LogLevel int32

const (
	LogLevel_DEBUG   LogLevel = 0
	LogLevel_INFO    LogLevel = 1
	LogLevel_WARNING LogLevel = 2
	LogLevel_ERROR   LogLevel = 3
	LogLevel_FATAL   LogLevel = 4
)

var LogLevel_name = map[int32]string{
	0: "DEBUG",
	1: "INFO",
	2: "WARNING",
	3: "ERROR",
	4: "FATAL",
}

var LogLevel_value = map[string]int32{
	"DEBUG":   0,
	"INFO":    1,
	"WARNING": 2,
	"ERROR":   3,
	"FATAL":   4,
}

func (x LogLevel) String() string {
	return proto.EnumName(LogLevel_name, int32(x))
}

func (LogLevel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1f19606b164b0edc, []int{0}
}

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f19606b164b0edc, []int{0}
}

func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (m *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(m, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

// -------------------------------------------------------------------------------
// Bytes is a special message type used to marshal & unmarshal unknown types as is
// -------------------------------------------------------------------------------
type Bytes struct {
	Val                  []byte   `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f19606b164b0edc, []int{1}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetVal() []byte {
	if m != nil {
		return m.Val
	}
	return nil
}

// --------------------------------------------------------------------------
// NetworkID uniquely identifies the network
// --------------------------------------------------------------------------
type NetworkID struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkID) Reset()         { *m = NetworkID{} }
func (m *NetworkID) String() string { return proto.CompactTextString(m) }
func (*NetworkID) ProtoMessage()    {}
func (*NetworkID) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f19606b164b0edc, []int{2}
}

func (m *NetworkID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkID.Unmarshal(m, b)
}
func (m *NetworkID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkID.Marshal(b, m, deterministic)
}
func (m *NetworkID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkID.Merge(m, src)
}
func (m *NetworkID) XXX_Size() int {
	return xxx_messageInfo_NetworkID.Size(m)
}
func (m *NetworkID) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkID.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkID proto.InternalMessageInfo

func (m *NetworkID) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// --------------------------------------------------------------------------
// IDList is a generic definition of an array of IDs (network, gateway, etc.)
// --------------------------------------------------------------------------
type IDList struct {
	Ids                  []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDList) Reset()         { *m = IDList{} }
func (m *IDList) String() string { return proto.CompactTextString(m) }
func (*IDList) ProtoMessage()    {}
func (*IDList) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f19606b164b0edc, []int{3}
}

func (m *IDList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDList.Unmarshal(m, b)
}
func (m *IDList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDList.Marshal(b, m, deterministic)
}
func (m *IDList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDList.Merge(m, src)
}
func (m *IDList) XXX_Size() int {
	return xxx_messageInfo_IDList.Size(m)
}
func (m *IDList) XXX_DiscardUnknown() {
	xxx_messageInfo_IDList.DiscardUnknown(m)
}

var xxx_messageInfo_IDList proto.InternalMessageInfo

func (m *IDList) GetIds() []string {
	if m != nil {
		return m.Ids
	}
	return nil
}

func init() {
	proto.RegisterEnum("magma.orc8r.LogLevel", LogLevel_name, LogLevel_value)
	proto.RegisterType((*Void)(nil), "magma.orc8r.Void")
	proto.RegisterType((*Bytes)(nil), "magma.orc8r.Bytes")
	proto.RegisterType((*NetworkID)(nil), "magma.orc8r.NetworkID")
	proto.RegisterType((*IDList)(nil), "magma.orc8r.IDList")
}

func init() { proto.RegisterFile("orc8r/protos/common.proto", fileDescriptor_1f19606b164b0edc) }

var fileDescriptor_1f19606b164b0edc = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x2c, 0x8f, 0x3d, 0x6b, 0xc3, 0x30,
	0x10, 0x86, 0xeb, 0x8f, 0xb8, 0xf1, 0xa5, 0x14, 0xa1, 0x29, 0x6e, 0x96, 0xe0, 0x29, 0x74, 0x88,
	0x87, 0x2e, 0x5d, 0x6d, 0x9c, 0x04, 0x83, 0x51, 0x40, 0xf4, 0x03, 0xba, 0x39, 0x91, 0x30, 0xa2,
	0x56, 0xaf, 0x48, 0x22, 0xa5, 0xff, 0xbe, 0x48, 0xcd, 0x74, 0xcf, 0x71, 0xef, 0xc1, 0xfb, 0x40,
	0x81, 0xe6, 0xfc, 0x6c, 0xaa, 0x6f, 0x83, 0x0e, 0x6d, 0x75, 0x46, 0xad, 0xf1, 0x6b, 0x1b, 0x36,
	0xba, 0xd0, 0xc3, 0xa8, 0x87, 0x6d, 0x08, 0x94, 0x19, 0xa4, 0x6f, 0xa8, 0x44, 0x59, 0xc0, 0xac,
	0xf9, 0x75, 0xd2, 0x52, 0x02, 0xc9, 0x65, 0x98, 0x96, 0xd1, 0x3a, 0xda, 0xdc, 0x71, 0x8f, 0xe5,
	0x0a, 0x72, 0x26, 0xdd, 0x0f, 0x9a, 0xcf, 0xae, 0xa5, 0xf7, 0x10, 0x2b, 0x11, 0xae, 0x39, 0x8f,
	0x95, 0x28, 0x1f, 0x20, 0xeb, 0xda, 0x5e, 0x59, 0xe7, 0x1f, 0x95, 0xb0, 0xcb, 0x68, 0x9d, 0x6c,
	0x72, 0xee, 0xf1, 0xb1, 0x81, 0x79, 0x8f, 0x63, 0x2f, 0x2f, 0x72, 0xa2, 0x39, 0xcc, 0xda, 0x5d,
	0xf3, 0x7a, 0x20, 0x37, 0x74, 0x0e, 0x69, 0xc7, 0xf6, 0x47, 0x12, 0xd1, 0x05, 0xdc, 0xbe, 0xd7,
	0x9c, 0x75, 0xec, 0x40, 0x62, 0x9f, 0xd8, 0x71, 0x7e, 0xe4, 0x24, 0xf1, 0xb8, 0xaf, 0x5f, 0xea,
	0x9e, 0xa4, 0xcd, 0xea, 0xa3, 0x08, 0x75, 0xab, 0x7f, 0x9f, 0x49, 0x9d, 0xaa, 0x11, 0xaf, 0x5a,
	0xa7, 0x2c, 0xcc, 0xa7, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x98, 0xbf, 0xd0, 0xed, 0x00,
	0x00, 0x00,
}
