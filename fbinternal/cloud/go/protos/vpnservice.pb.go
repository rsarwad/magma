// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vpnservice.proto

package protos

import (
	context "context"
	fmt "fmt"
	math "math"

	protos "magma/orc8r/lib/go/protos"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type VPNCertRequest struct {
	// Represents an x509 certificate request (.csr file)
	Request              []byte   `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VPNCertRequest) Reset()         { *m = VPNCertRequest{} }
func (m *VPNCertRequest) String() string { return proto.CompactTextString(m) }
func (*VPNCertRequest) ProtoMessage()    {}
func (*VPNCertRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_473742cf0fb06c4c, []int{0}
}

func (m *VPNCertRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VPNCertRequest.Unmarshal(m, b)
}
func (m *VPNCertRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VPNCertRequest.Marshal(b, m, deterministic)
}
func (m *VPNCertRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VPNCertRequest.Merge(m, src)
}
func (m *VPNCertRequest) XXX_Size() int {
	return xxx_messageInfo_VPNCertRequest.Size(m)
}
func (m *VPNCertRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VPNCertRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VPNCertRequest proto.InternalMessageInfo

func (m *VPNCertRequest) GetRequest() []byte {
	if m != nil {
		return m.Request
	}
	return nil
}

type VPNCertificate struct {
	// Represents an x509 certificate used to connect with OpenVPN (.crt file)
	Serial               string   `protobuf:"bytes,1,opt,name=serial,proto3" json:"serial,omitempty"`
	Cert                 []byte   `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VPNCertificate) Reset()         { *m = VPNCertificate{} }
func (m *VPNCertificate) String() string { return proto.CompactTextString(m) }
func (*VPNCertificate) ProtoMessage()    {}
func (*VPNCertificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_473742cf0fb06c4c, []int{1}
}

func (m *VPNCertificate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VPNCertificate.Unmarshal(m, b)
}
func (m *VPNCertificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VPNCertificate.Marshal(b, m, deterministic)
}
func (m *VPNCertificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VPNCertificate.Merge(m, src)
}
func (m *VPNCertificate) XXX_Size() int {
	return xxx_messageInfo_VPNCertificate.Size(m)
}
func (m *VPNCertificate) XXX_DiscardUnknown() {
	xxx_messageInfo_VPNCertificate.DiscardUnknown(m)
}

var xxx_messageInfo_VPNCertificate proto.InternalMessageInfo

func (m *VPNCertificate) GetSerial() string {
	if m != nil {
		return m.Serial
	}
	return ""
}

func (m *VPNCertificate) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

type PSK struct {
	TaKey                []byte   `protobuf:"bytes,1,opt,name=ta_key,json=taKey,proto3" json:"ta_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PSK) Reset()         { *m = PSK{} }
func (m *PSK) String() string { return proto.CompactTextString(m) }
func (*PSK) ProtoMessage()    {}
func (*PSK) Descriptor() ([]byte, []int) {
	return fileDescriptor_473742cf0fb06c4c, []int{2}
}

func (m *PSK) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PSK.Unmarshal(m, b)
}
func (m *PSK) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PSK.Marshal(b, m, deterministic)
}
func (m *PSK) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PSK.Merge(m, src)
}
func (m *PSK) XXX_Size() int {
	return xxx_messageInfo_PSK.Size(m)
}
func (m *PSK) XXX_DiscardUnknown() {
	xxx_messageInfo_PSK.DiscardUnknown(m)
}

var xxx_messageInfo_PSK proto.InternalMessageInfo

func (m *PSK) GetTaKey() []byte {
	if m != nil {
		return m.TaKey
	}
	return nil
}

func init() {
	proto.RegisterType((*VPNCertRequest)(nil), "magma.fbinternal.VPNCertRequest")
	proto.RegisterType((*VPNCertificate)(nil), "magma.fbinternal.VPNCertificate")
	proto.RegisterType((*PSK)(nil), "magma.fbinternal.PSK")
}

func init() { proto.RegisterFile("vpnservice.proto", fileDescriptor_473742cf0fb06c4c) }

var fileDescriptor_473742cf0fb06c4c = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0xc6, 0xdb, 0xf7, 0xb5, 0x15, 0x47, 0x91, 0xba, 0x52, 0xa9, 0xa1, 0x87, 0xb0, 0x27, 0xf1,
	0xb0, 0x41, 0xbd, 0xf4, 0xe0, 0xa5, 0xe6, 0xe0, 0x21, 0x50, 0x42, 0x02, 0x39, 0x78, 0x91, 0x6d,
	0x3a, 0x29, 0x8b, 0x49, 0xb6, 0x6e, 0xb6, 0x85, 0x7e, 0x49, 0x3f, 0x93, 0x64, 0x77, 0x55, 0xe2,
	0x9f, 0xd3, 0xce, 0xcc, 0xf3, 0x30, 0x3c, 0xbf, 0x59, 0x18, 0xed, 0x36, 0x75, 0x83, 0x6a, 0x27,
	0x72, 0x64, 0x1b, 0x25, 0xb5, 0x24, 0xa3, 0x8a, 0xaf, 0x2b, 0xce, 0x8a, 0xa5, 0xa8, 0x35, 0xaa,
	0x9a, 0x97, 0xde, 0x54, 0xaa, 0x7c, 0xa6, 0x02, 0x23, 0x37, 0x41, 0x8e, 0x4a, 0x8b, 0x42, 0xa0,
	0xb2, 0x7e, 0xef, 0xb2, 0xab, 0xca, 0xaa, 0x92, 0xb5, 0x95, 0xe8, 0x35, 0x9c, 0x66, 0xf1, 0x22,
	0x44, 0xa5, 0x13, 0x7c, 0xdd, 0x62, 0xa3, 0xc9, 0x04, 0x0e, 0x95, 0x2d, 0x27, 0x7d, 0xbf, 0x7f,
	0x75, 0x92, 0x7c, 0xb4, 0xf4, 0xfe, 0xd3, 0x2b, 0x0a, 0x91, 0x73, 0x8d, 0xe4, 0x02, 0x86, 0x0d,
	0x2a, 0xc1, 0x4b, 0x63, 0x3d, 0x4a, 0x5c, 0x47, 0x08, 0x1c, 0xb4, 0x19, 0x26, 0xff, 0xcc, 0x02,
	0x53, 0xd3, 0x29, 0xfc, 0x8f, 0xd3, 0x88, 0x8c, 0x61, 0xa8, 0xf9, 0xf3, 0x0b, 0xee, 0xdd, 0xf6,
	0x81, 0xe6, 0x11, 0xee, 0x6f, 0xdf, 0xfa, 0x00, 0x59, 0xbc, 0x48, 0x2d, 0x27, 0xb9, 0x81, 0xc1,
	0x23, 0xea, 0x70, 0x4e, 0xce, 0x98, 0x65, 0x35, 0x04, 0x2c, 0x93, 0x62, 0xe5, 0x9d, 0x77, 0x46,
	0xe1, 0xbc, 0x0d, 0x44, 0x7b, 0x24, 0x85, 0x63, 0x87, 0xd0, 0x0e, 0x88, 0xcf, 0xbe, 0x1f, 0x89,
	0x75, 0x41, 0xbd, 0xbf, 0x1d, 0x0e, 0x8f, 0xf6, 0xc8, 0x0c, 0xc0, 0xd9, 0xdb, 0xec, 0xbf, 0x84,
	0x19, 0xff, 0x5c, 0x12, 0xa7, 0x11, 0xed, 0x3d, 0xd0, 0x27, 0xdf, 0x5e, 0xfd, 0x4b, 0x09, 0xf2,
	0x52, 0x6e, 0x57, 0xc1, 0x5a, 0xba, 0x9f, 0x58, 0x0e, 0xcd, 0x7b, 0xf7, 0x1e, 0x00, 0x00, 0xff,
	0xff, 0x16, 0x3c, 0xd5, 0x6a, 0xe2, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// VPNServiceClient is the client API for VPNService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VPNServiceClient interface {
	// Return the CA (ca.crt)
	GetCA(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*protos.CACert, error)
	// Given a request (client.csr), return a signed certificate (client.crt)
	RequestCert(ctx context.Context, in *VPNCertRequest, opts ...grpc.CallOption) (*VPNCertificate, error)
	// Request for the PSK (preshared key, i.e. tls-auth key used in openvpn)
	// See https://community.openvpn.net/openvpn/wiki/Hardening#Useof--tls-auth
	// for detail
	RequestPSK(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*PSK, error)
}

type vPNServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVPNServiceClient(cc grpc.ClientConnInterface) VPNServiceClient {
	return &vPNServiceClient{cc}
}

func (c *vPNServiceClient) GetCA(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*protos.CACert, error) {
	out := new(protos.CACert)
	err := c.cc.Invoke(ctx, "/magma.fbinternal.VPNService/GetCA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vPNServiceClient) RequestCert(ctx context.Context, in *VPNCertRequest, opts ...grpc.CallOption) (*VPNCertificate, error) {
	out := new(VPNCertificate)
	err := c.cc.Invoke(ctx, "/magma.fbinternal.VPNService/RequestCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vPNServiceClient) RequestPSK(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*PSK, error) {
	out := new(PSK)
	err := c.cc.Invoke(ctx, "/magma.fbinternal.VPNService/RequestPSK", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VPNServiceServer is the server API for VPNService service.
type VPNServiceServer interface {
	// Return the CA (ca.crt)
	GetCA(context.Context, *protos.Void) (*protos.CACert, error)
	// Given a request (client.csr), return a signed certificate (client.crt)
	RequestCert(context.Context, *VPNCertRequest) (*VPNCertificate, error)
	// Request for the PSK (preshared key, i.e. tls-auth key used in openvpn)
	// See https://community.openvpn.net/openvpn/wiki/Hardening#Useof--tls-auth
	// for detail
	RequestPSK(context.Context, *protos.Void) (*PSK, error)
}

// UnimplementedVPNServiceServer can be embedded to have forward compatible implementations.
type UnimplementedVPNServiceServer struct {
}

func (*UnimplementedVPNServiceServer) GetCA(ctx context.Context, req *protos.Void) (*protos.CACert, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCA not implemented")
}
func (*UnimplementedVPNServiceServer) RequestCert(ctx context.Context, req *VPNCertRequest) (*VPNCertificate, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestCert not implemented")
}
func (*UnimplementedVPNServiceServer) RequestPSK(ctx context.Context, req *protos.Void) (*PSK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestPSK not implemented")
}

func RegisterVPNServiceServer(s *grpc.Server, srv VPNServiceServer) {
	s.RegisterService(&_VPNService_serviceDesc, srv)
}

func _VPNService_GetCA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protos.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VPNServiceServer).GetCA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.fbinternal.VPNService/GetCA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VPNServiceServer).GetCA(ctx, req.(*protos.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VPNService_RequestCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VPNCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VPNServiceServer).RequestCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.fbinternal.VPNService/RequestCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VPNServiceServer).RequestCert(ctx, req.(*VPNCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VPNService_RequestPSK_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protos.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VPNServiceServer).RequestPSK(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.fbinternal.VPNService/RequestPSK",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VPNServiceServer).RequestPSK(ctx, req.(*protos.Void))
	}
	return interceptor(ctx, in, info, handler)
}

var _VPNService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "magma.fbinternal.VPNService",
	HandlerType: (*VPNServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCA",
			Handler:    _VPNService_GetCA_Handler,
		},
		{
			MethodName: "RequestCert",
			Handler:    _VPNService_RequestCert_Handler,
		},
		{
			MethodName: "RequestPSK",
			Handler:    _VPNService_RequestPSK_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vpnservice.proto",
}
