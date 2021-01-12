// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ScoreServiceClient is the client API for ScoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScoreServiceClient interface {
	// Register the tests to be expected for this test run.
	Register(ctx context.Context, in *Tests, opts ...grpc.CallOption) (*Void, error)
	// Notify sends one score for each test.
	Notify(ctx context.Context, opts ...grpc.CallOption) (ScoreService_NotifyClient, error)
}

type scoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScoreServiceClient(cc grpc.ClientConnInterface) ScoreServiceClient {
	return &scoreServiceClient{cc}
}

func (c *scoreServiceClient) Register(ctx context.Context, in *Tests, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/proto.ScoreService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scoreServiceClient) Notify(ctx context.Context, opts ...grpc.CallOption) (ScoreService_NotifyClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ScoreService_serviceDesc.Streams[0], "/proto.ScoreService/Notify", opts...)
	if err != nil {
		return nil, err
	}
	x := &scoreServiceNotifyClient{stream}
	return x, nil
}

type ScoreService_NotifyClient interface {
	Send(*Score) error
	CloseAndRecv() (*Void, error)
	grpc.ClientStream
}

type scoreServiceNotifyClient struct {
	grpc.ClientStream
}

func (x *scoreServiceNotifyClient) Send(m *Score) error {
	return x.ClientStream.SendMsg(m)
}

func (x *scoreServiceNotifyClient) CloseAndRecv() (*Void, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Void)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ScoreServiceServer is the server API for ScoreService service.
// All implementations must embed UnimplementedScoreServiceServer
// for forward compatibility
type ScoreServiceServer interface {
	// Register the tests to be expected for this test run.
	Register(context.Context, *Tests) (*Void, error)
	// Notify sends one score for each test.
	Notify(ScoreService_NotifyServer) error
	mustEmbedUnimplementedScoreServiceServer()
}

// UnimplementedScoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedScoreServiceServer struct {
}

func (UnimplementedScoreServiceServer) Register(context.Context, *Tests) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedScoreServiceServer) Notify(ScoreService_NotifyServer) error {
	return status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedScoreServiceServer) mustEmbedUnimplementedScoreServiceServer() {}

// UnsafeScoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScoreServiceServer will
// result in compilation errors.
type UnsafeScoreServiceServer interface {
	mustEmbedUnimplementedScoreServiceServer()
}

func RegisterScoreServiceServer(s grpc.ServiceRegistrar, srv ScoreServiceServer) {
	s.RegisterService(&_ScoreService_serviceDesc, srv)
}

func _ScoreService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tests)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScoreServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ScoreService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScoreServiceServer).Register(ctx, req.(*Tests))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScoreService_Notify_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ScoreServiceServer).Notify(&scoreServiceNotifyServer{stream})
}

type ScoreService_NotifyServer interface {
	SendAndClose(*Void) error
	Recv() (*Score, error)
	grpc.ServerStream
}

type scoreServiceNotifyServer struct {
	grpc.ServerStream
}

func (x *scoreServiceNotifyServer) SendAndClose(m *Void) error {
	return x.ServerStream.SendMsg(m)
}

func (x *scoreServiceNotifyServer) Recv() (*Score, error) {
	m := new(Score)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ScoreService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ScoreService",
	HandlerType: (*ScoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ScoreService_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Notify",
			Handler:       _ScoreService_Notify_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "score.proto",
}
