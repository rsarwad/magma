// Code generated by protoc-gen-go. DO NOT EDIT.
// source: feg/protos/envoy_controller.proto

package protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protos "magma/lte/cloud/go/protos"
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

type Header struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_c190610d29559b01, []int{0}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Header) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type AddUEHeaderEnrichmentRequest struct {
	UeIp                 *protos.IPAddress `protobuf:"bytes,1,opt,name=ue_ip,json=ueIp,proto3" json:"ue_ip,omitempty"`
	Websites             []string          `protobuf:"bytes,2,rep,name=websites,proto3" json:"websites,omitempty"`
	Headers              []*Header         `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AddUEHeaderEnrichmentRequest) Reset()         { *m = AddUEHeaderEnrichmentRequest{} }
func (m *AddUEHeaderEnrichmentRequest) String() string { return proto.CompactTextString(m) }
func (*AddUEHeaderEnrichmentRequest) ProtoMessage()    {}
func (*AddUEHeaderEnrichmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c190610d29559b01, []int{1}
}

func (m *AddUEHeaderEnrichmentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddUEHeaderEnrichmentRequest.Unmarshal(m, b)
}
func (m *AddUEHeaderEnrichmentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddUEHeaderEnrichmentRequest.Marshal(b, m, deterministic)
}
func (m *AddUEHeaderEnrichmentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddUEHeaderEnrichmentRequest.Merge(m, src)
}
func (m *AddUEHeaderEnrichmentRequest) XXX_Size() int {
	return xxx_messageInfo_AddUEHeaderEnrichmentRequest.Size(m)
}
func (m *AddUEHeaderEnrichmentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddUEHeaderEnrichmentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddUEHeaderEnrichmentRequest proto.InternalMessageInfo

func (m *AddUEHeaderEnrichmentRequest) GetUeIp() *protos.IPAddress {
	if m != nil {
		return m.UeIp
	}
	return nil
}

func (m *AddUEHeaderEnrichmentRequest) GetWebsites() []string {
	if m != nil {
		return m.Websites
	}
	return nil
}

func (m *AddUEHeaderEnrichmentRequest) GetHeaders() []*Header {
	if m != nil {
		return m.Headers
	}
	return nil
}

type AddUEHeaderEnrichmentResult struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddUEHeaderEnrichmentResult) Reset()         { *m = AddUEHeaderEnrichmentResult{} }
func (m *AddUEHeaderEnrichmentResult) String() string { return proto.CompactTextString(m) }
func (*AddUEHeaderEnrichmentResult) ProtoMessage()    {}
func (*AddUEHeaderEnrichmentResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_c190610d29559b01, []int{2}
}

func (m *AddUEHeaderEnrichmentResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddUEHeaderEnrichmentResult.Unmarshal(m, b)
}
func (m *AddUEHeaderEnrichmentResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddUEHeaderEnrichmentResult.Marshal(b, m, deterministic)
}
func (m *AddUEHeaderEnrichmentResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddUEHeaderEnrichmentResult.Merge(m, src)
}
func (m *AddUEHeaderEnrichmentResult) XXX_Size() int {
	return xxx_messageInfo_AddUEHeaderEnrichmentResult.Size(m)
}
func (m *AddUEHeaderEnrichmentResult) XXX_DiscardUnknown() {
	xxx_messageInfo_AddUEHeaderEnrichmentResult.DiscardUnknown(m)
}

var xxx_messageInfo_AddUEHeaderEnrichmentResult proto.InternalMessageInfo

type DeactivateUEHeaderEnrichmentRequest struct {
	UeIp                 *protos.IPAddress `protobuf:"bytes,1,opt,name=ue_ip,json=ueIp,proto3" json:"ue_ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DeactivateUEHeaderEnrichmentRequest) Reset()         { *m = DeactivateUEHeaderEnrichmentRequest{} }
func (m *DeactivateUEHeaderEnrichmentRequest) String() string { return proto.CompactTextString(m) }
func (*DeactivateUEHeaderEnrichmentRequest) ProtoMessage()    {}
func (*DeactivateUEHeaderEnrichmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c190610d29559b01, []int{3}
}

func (m *DeactivateUEHeaderEnrichmentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest.Unmarshal(m, b)
}
func (m *DeactivateUEHeaderEnrichmentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest.Marshal(b, m, deterministic)
}
func (m *DeactivateUEHeaderEnrichmentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest.Merge(m, src)
}
func (m *DeactivateUEHeaderEnrichmentRequest) XXX_Size() int {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest.Size(m)
}
func (m *DeactivateUEHeaderEnrichmentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeactivateUEHeaderEnrichmentRequest proto.InternalMessageInfo

func (m *DeactivateUEHeaderEnrichmentRequest) GetUeIp() *protos.IPAddress {
	if m != nil {
		return m.UeIp
	}
	return nil
}

type DeactivateUEHeaderEnrichmentResult struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeactivateUEHeaderEnrichmentResult) Reset()         { *m = DeactivateUEHeaderEnrichmentResult{} }
func (m *DeactivateUEHeaderEnrichmentResult) String() string { return proto.CompactTextString(m) }
func (*DeactivateUEHeaderEnrichmentResult) ProtoMessage()    {}
func (*DeactivateUEHeaderEnrichmentResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_c190610d29559b01, []int{4}
}

func (m *DeactivateUEHeaderEnrichmentResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentResult.Unmarshal(m, b)
}
func (m *DeactivateUEHeaderEnrichmentResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentResult.Marshal(b, m, deterministic)
}
func (m *DeactivateUEHeaderEnrichmentResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeactivateUEHeaderEnrichmentResult.Merge(m, src)
}
func (m *DeactivateUEHeaderEnrichmentResult) XXX_Size() int {
	return xxx_messageInfo_DeactivateUEHeaderEnrichmentResult.Size(m)
}
func (m *DeactivateUEHeaderEnrichmentResult) XXX_DiscardUnknown() {
	xxx_messageInfo_DeactivateUEHeaderEnrichmentResult.DiscardUnknown(m)
}

var xxx_messageInfo_DeactivateUEHeaderEnrichmentResult proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Header)(nil), "magma.feg.Header")
	proto.RegisterType((*AddUEHeaderEnrichmentRequest)(nil), "magma.feg.AddUEHeaderEnrichmentRequest")
	proto.RegisterType((*AddUEHeaderEnrichmentResult)(nil), "magma.feg.AddUEHeaderEnrichmentResult")
	proto.RegisterType((*DeactivateUEHeaderEnrichmentRequest)(nil), "magma.feg.DeactivateUEHeaderEnrichmentRequest")
	proto.RegisterType((*DeactivateUEHeaderEnrichmentResult)(nil), "magma.feg.DeactivateUEHeaderEnrichmentResult")
}

func init() { proto.RegisterFile("feg/protos/envoy_controller.proto", fileDescriptor_c190610d29559b01) }

var fileDescriptor_c190610d29559b01 = []byte{
	// 353 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xcf, 0x4f, 0xea, 0x40,
	0x10, 0xc7, 0x5f, 0xf9, 0xf5, 0x1e, 0xc3, 0xe9, 0x6d, 0x30, 0xa9, 0x05, 0x12, 0xac, 0x46, 0x31,
	0xc6, 0x36, 0xa9, 0x7f, 0x01, 0x46, 0x12, 0xb9, 0x91, 0x26, 0x5e, 0xbc, 0x90, 0xb6, 0x3b, 0x94,
	0x9a, 0x6d, 0x17, 0x77, 0xb7, 0x18, 0x4e, 0xfe, 0x0f, 0xfe, 0xaf, 0xde, 0x0d, 0xbb, 0x42, 0x38,
	0x28, 0x72, 0xf0, 0xd4, 0xce, 0xcc, 0x77, 0x66, 0x3e, 0x33, 0x3b, 0x70, 0x32, 0xc3, 0xd4, 0x5f,
	0x08, 0xae, 0xb8, 0xf4, 0xb1, 0x58, 0xf2, 0xd5, 0x34, 0xe1, 0x85, 0x12, 0x9c, 0x31, 0x14, 0x9e,
	0xf6, 0x93, 0x66, 0x1e, 0xa5, 0x79, 0xe4, 0xcd, 0x30, 0x75, 0x1c, 0xa6, 0x70, 0xa3, 0xce, 0x79,
	0x9c, 0xb1, 0x4c, 0xad, 0xa8, 0x91, 0x39, 0xbd, 0x9d, 0x98, 0x2c, 0x63, 0x99, 0x88, 0x2c, 0x46,
	0x41, 0x63, 0x13, 0x76, 0x03, 0x68, 0xdc, 0x63, 0x44, 0x51, 0x10, 0x02, 0xb5, 0x22, 0xca, 0xd1,
	0xb6, 0xfa, 0xd6, 0xa0, 0x19, 0xea, 0x7f, 0xd2, 0x86, 0xfa, 0x32, 0x62, 0x25, 0xda, 0x15, 0xed,
	0x34, 0x86, 0xfb, 0x66, 0x41, 0x77, 0x48, 0xe9, 0xc3, 0xc8, 0x64, 0x8e, 0x0a, 0x91, 0x25, 0xf3,
	0x1c, 0x0b, 0x15, 0xe2, 0x73, 0x89, 0x52, 0x91, 0x4b, 0xa8, 0x97, 0x38, 0xcd, 0x16, 0xba, 0x56,
	0x2b, 0x68, 0x7b, 0x06, 0x95, 0x29, 0xf4, 0xc6, 0x93, 0x21, 0xa5, 0x02, 0xa5, 0x0c, 0x6b, 0x25,
	0x8e, 0x17, 0xc4, 0x81, 0x7f, 0x2f, 0x18, 0xcb, 0x4c, 0xa1, 0xb4, 0x2b, 0xfd, 0xea, 0xa0, 0x19,
	0x6e, 0x6d, 0x72, 0x05, 0x7f, 0xe7, 0xba, 0x83, 0xb4, 0xab, 0xfd, 0xea, 0xa0, 0x15, 0xfc, 0xf7,
	0xb6, 0x33, 0x7b, 0xa6, 0x77, 0xb8, 0x51, 0xb8, 0x3d, 0xe8, 0x7c, 0xc3, 0x24, 0x4b, 0xa6, 0xdc,
	0x09, 0x9c, 0xde, 0x61, 0x94, 0xa8, 0x6c, 0x19, 0x29, 0xfc, 0x0d, 0x72, 0xf7, 0x0c, 0xdc, 0xfd,
	0x15, 0xd7, 0x7d, 0x83, 0x77, 0x0b, 0x1a, 0xa3, 0xf5, 0x03, 0x52, 0xf2, 0x04, 0x47, 0x5f, 0x12,
	0x92, 0x8b, 0x9d, 0xb1, 0xf6, 0xed, 0xd5, 0x39, 0xff, 0x59, 0xa8, 0x87, 0xfd, 0x43, 0x5e, 0xa1,
	0xbb, 0x0f, 0x8e, 0x78, 0x3b, 0x95, 0x0e, 0xd8, 0x8b, 0x73, 0x7d, 0xb0, 0xde, 0x00, 0xdc, 0x76,
	0x1e, 0x8f, 0x75, 0x86, 0xbf, 0x3e, 0xe4, 0x84, 0xf1, 0x92, 0xfa, 0x29, 0xff, 0xbc, 0xc3, 0xb8,
	0xa1, 0xbf, 0x37, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x12, 0xf1, 0x57, 0x1e, 0xe6, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// EnvoydClient is the client API for Envoyd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EnvoydClient interface {
	// Add UE header enrichment configuration
	AddUEHeaderEnrichment(ctx context.Context, in *AddUEHeaderEnrichmentRequest, opts ...grpc.CallOption) (*AddUEHeaderEnrichmentResult, error)
	DeactivateUEHeaderEnrichment(ctx context.Context, in *DeactivateUEHeaderEnrichmentRequest, opts ...grpc.CallOption) (*DeactivateUEHeaderEnrichmentResult, error)
}

type envoydClient struct {
	cc grpc.ClientConnInterface
}

func NewEnvoydClient(cc grpc.ClientConnInterface) EnvoydClient {
	return &envoydClient{cc}
}

func (c *envoydClient) AddUEHeaderEnrichment(ctx context.Context, in *AddUEHeaderEnrichmentRequest, opts ...grpc.CallOption) (*AddUEHeaderEnrichmentResult, error) {
	out := new(AddUEHeaderEnrichmentResult)
	err := c.cc.Invoke(ctx, "/magma.feg.Envoyd/AddUEHeaderEnrichment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envoydClient) DeactivateUEHeaderEnrichment(ctx context.Context, in *DeactivateUEHeaderEnrichmentRequest, opts ...grpc.CallOption) (*DeactivateUEHeaderEnrichmentResult, error) {
	out := new(DeactivateUEHeaderEnrichmentResult)
	err := c.cc.Invoke(ctx, "/magma.feg.Envoyd/DeactivateUEHeaderEnrichment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnvoydServer is the server API for Envoyd service.
type EnvoydServer interface {
	// Add UE header enrichment configuration
	AddUEHeaderEnrichment(context.Context, *AddUEHeaderEnrichmentRequest) (*AddUEHeaderEnrichmentResult, error)
	DeactivateUEHeaderEnrichment(context.Context, *DeactivateUEHeaderEnrichmentRequest) (*DeactivateUEHeaderEnrichmentResult, error)
}

// UnimplementedEnvoydServer can be embedded to have forward compatible implementations.
type UnimplementedEnvoydServer struct {
}

func (*UnimplementedEnvoydServer) AddUEHeaderEnrichment(ctx context.Context, req *AddUEHeaderEnrichmentRequest) (*AddUEHeaderEnrichmentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUEHeaderEnrichment not implemented")
}
func (*UnimplementedEnvoydServer) DeactivateUEHeaderEnrichment(ctx context.Context, req *DeactivateUEHeaderEnrichmentRequest) (*DeactivateUEHeaderEnrichmentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateUEHeaderEnrichment not implemented")
}

func RegisterEnvoydServer(s *grpc.Server, srv EnvoydServer) {
	s.RegisterService(&_Envoyd_serviceDesc, srv)
}

func _Envoyd_AddUEHeaderEnrichment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUEHeaderEnrichmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvoydServer).AddUEHeaderEnrichment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.feg.Envoyd/AddUEHeaderEnrichment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvoydServer).AddUEHeaderEnrichment(ctx, req.(*AddUEHeaderEnrichmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Envoyd_DeactivateUEHeaderEnrichment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateUEHeaderEnrichmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvoydServer).DeactivateUEHeaderEnrichment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.feg.Envoyd/DeactivateUEHeaderEnrichment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvoydServer).DeactivateUEHeaderEnrichment(ctx, req.(*DeactivateUEHeaderEnrichmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Envoyd_serviceDesc = grpc.ServiceDesc{
	ServiceName: "magma.feg.Envoyd",
	HandlerType: (*EnvoydServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUEHeaderEnrichment",
			Handler:    _Envoyd_AddUEHeaderEnrichment_Handler,
		},
		{
			MethodName: "DeactivateUEHeaderEnrichment",
			Handler:    _Envoyd_DeactivateUEHeaderEnrichment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feg/protos/envoy_controller.proto",
}
