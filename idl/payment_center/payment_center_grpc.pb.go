// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.17.3
// source: idl/payment_center/payment_center.proto

package payment_center

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PaymentCenter_AddPaymentApply_FullMethodName    = "/idl.PaymentCenter/AddPaymentApply"
	PaymentCenter_UpdatePaymentApply_FullMethodName = "/idl.PaymentCenter/UpdatePaymentApply"
)

// PaymentCenterClient is the client API for PaymentCenter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentCenterClient interface {
	// 付款申请单-创建
	AddPaymentApply(ctx context.Context, in *AddPAReq, opts ...grpc.CallOption) (*AddPAResp, error)
	// 付款申请单-创建
	UpdatePaymentApply(ctx context.Context, in *UpdatePAReq, opts ...grpc.CallOption) (*UpdatePAResp, error)
}

type paymentCenterClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentCenterClient(cc grpc.ClientConnInterface) PaymentCenterClient {
	return &paymentCenterClient{cc}
}

func (c *paymentCenterClient) AddPaymentApply(ctx context.Context, in *AddPAReq, opts ...grpc.CallOption) (*AddPAResp, error) {
	out := new(AddPAResp)
	err := c.cc.Invoke(ctx, PaymentCenter_AddPaymentApply_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentCenterClient) UpdatePaymentApply(ctx context.Context, in *UpdatePAReq, opts ...grpc.CallOption) (*UpdatePAResp, error) {
	out := new(UpdatePAResp)
	err := c.cc.Invoke(ctx, PaymentCenter_UpdatePaymentApply_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentCenterServer is the server API for PaymentCenter service.
// All implementations must embed UnimplementedPaymentCenterServer
// for forward compatibility
type PaymentCenterServer interface {
	// 付款申请单-创建
	AddPaymentApply(context.Context, *AddPAReq) (*AddPAResp, error)
	// 付款申请单-创建
	UpdatePaymentApply(context.Context, *UpdatePAReq) (*UpdatePAResp, error)
	mustEmbedUnimplementedPaymentCenterServer()
}

// UnimplementedPaymentCenterServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentCenterServer struct {
}

func (UnimplementedPaymentCenterServer) AddPaymentApply(context.Context, *AddPAReq) (*AddPAResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPaymentApply not implemented")
}
func (UnimplementedPaymentCenterServer) UpdatePaymentApply(context.Context, *UpdatePAReq) (*UpdatePAResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentApply not implemented")
}
func (UnimplementedPaymentCenterServer) mustEmbedUnimplementedPaymentCenterServer() {}

// UnsafePaymentCenterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentCenterServer will
// result in compilation errors.
type UnsafePaymentCenterServer interface {
	mustEmbedUnimplementedPaymentCenterServer()
}

func RegisterPaymentCenterServer(s grpc.ServiceRegistrar, srv PaymentCenterServer) {
	s.RegisterService(&PaymentCenter_ServiceDesc, srv)
}

func _PaymentCenter_AddPaymentApply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPAReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentCenterServer).AddPaymentApply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentCenter_AddPaymentApply_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentCenterServer).AddPaymentApply(ctx, req.(*AddPAReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentCenter_UpdatePaymentApply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePAReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentCenterServer).UpdatePaymentApply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentCenter_UpdatePaymentApply_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentCenterServer).UpdatePaymentApply(ctx, req.(*UpdatePAReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentCenter_ServiceDesc is the grpc.ServiceDesc for PaymentCenter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentCenter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "idl.PaymentCenter",
	HandlerType: (*PaymentCenterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPaymentApply",
			Handler:    _PaymentCenter_AddPaymentApply_Handler,
		},
		{
			MethodName: "UpdatePaymentApply",
			Handler:    _PaymentCenter_UpdatePaymentApply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/payment_center/payment_center.proto",
}
