// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// protoc-gen-go-vtproto version: v0.3.1-0.20220817155510-0ae748fd2007
// source: api/v1/backup_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	storage "github.com/stackrox/rox/generated/storage"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	proto "google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *GetExternalBackupsResponse) CloneVT() *GetExternalBackupsResponse {
	if m == nil {
		return (*GetExternalBackupsResponse)(nil)
	}
	r := &GetExternalBackupsResponse{}
	if rhs := m.ExternalBackups; rhs != nil {
		tmpContainer := make([]*storage.ExternalBackup, len(rhs))
		for k, v := range rhs {
			if vtpb, ok := interface{}(v).(interface {
				CloneVT() *storage.ExternalBackup
			}); ok {
				tmpContainer[k] = vtpb.CloneVT()
			} else {
				tmpContainer[k] = proto.Clone(v).(*storage.ExternalBackup)
			}
		}
		r.ExternalBackups = tmpContainer
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *GetExternalBackupsResponse) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (m *UpdateExternalBackupRequest) CloneVT() *UpdateExternalBackupRequest {
	if m == nil {
		return (*UpdateExternalBackupRequest)(nil)
	}
	r := &UpdateExternalBackupRequest{
		UpdatePassword: m.UpdatePassword,
	}
	if rhs := m.ExternalBackup; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface {
			CloneVT() *storage.ExternalBackup
		}); ok {
			r.ExternalBackup = vtpb.CloneVT()
		} else {
			r.ExternalBackup = proto.Clone(rhs).(*storage.ExternalBackup)
		}
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *UpdateExternalBackupRequest) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (this *GetExternalBackupsResponse) EqualVT(that *GetExternalBackupsResponse) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if len(this.ExternalBackups) != len(that.ExternalBackups) {
		return false
	}
	for i, vx := range this.ExternalBackups {
		vy := that.ExternalBackups[i]
		if p, q := vx, vy; p != q {
			if p == nil {
				p = &storage.ExternalBackup{}
			}
			if q == nil {
				q = &storage.ExternalBackup{}
			}
			if equal, ok := interface{}(p).(interface {
				EqualVT(*storage.ExternalBackup) bool
			}); ok {
				if !equal.EqualVT(q) {
					return false
				}
			} else if !proto.Equal(p, q) {
				return false
			}
		}
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *UpdateExternalBackupRequest) EqualVT(that *UpdateExternalBackupRequest) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if equal, ok := interface{}(this.ExternalBackup).(interface {
		EqualVT(*storage.ExternalBackup) bool
	}); ok {
		if !equal.EqualVT(that.ExternalBackup) {
			return false
		}
	} else if !proto.Equal(this.ExternalBackup, that.ExternalBackup) {
		return false
	}
	if this.UpdatePassword != that.UpdatePassword {
		return false
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExternalBackupServiceClient is the client API for ExternalBackupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExternalBackupServiceClient interface {
	// GetExternalBackup returns the external backup configuration given its ID.
	GetExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*storage.ExternalBackup, error)
	// GetExternalBackups returns all external backup configurations.
	GetExternalBackups(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetExternalBackupsResponse, error)
	// PostExternalBackup creates an external backup configuration.
	PostExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*storage.ExternalBackup, error)
	// PutExternalBackup modifies a given external backup, without using stored credential reconciliation.
	PutExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*storage.ExternalBackup, error)
	// TestExternalBackup tests an external backup configuration.
	TestExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*Empty, error)
	// DeleteExternalBackup removes an external backup configuration given its ID.
	DeleteExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*Empty, error)
	// TriggerExternalBackup initiates an external backup for the given configuration.
	TriggerExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*Empty, error)
	// UpdateExternalBackup modifies a given external backup, with optional stored credential reconciliation.
	UpdateExternalBackup(ctx context.Context, in *UpdateExternalBackupRequest, opts ...grpc.CallOption) (*storage.ExternalBackup, error)
	// TestUpdatedExternalBackup checks if the given external backup is correctly configured, with optional stored credential reconciliation.
	TestUpdatedExternalBackup(ctx context.Context, in *UpdateExternalBackupRequest, opts ...grpc.CallOption) (*Empty, error)
}

type externalBackupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExternalBackupServiceClient(cc grpc.ClientConnInterface) ExternalBackupServiceClient {
	return &externalBackupServiceClient{cc}
}

func (c *externalBackupServiceClient) GetExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*storage.ExternalBackup, error) {
	out := new(storage.ExternalBackup)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/GetExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) GetExternalBackups(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetExternalBackupsResponse, error) {
	out := new(GetExternalBackupsResponse)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/GetExternalBackups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) PostExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*storage.ExternalBackup, error) {
	out := new(storage.ExternalBackup)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/PostExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) PutExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*storage.ExternalBackup, error) {
	out := new(storage.ExternalBackup)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/PutExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) TestExternalBackup(ctx context.Context, in *storage.ExternalBackup, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/TestExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) DeleteExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/DeleteExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) TriggerExternalBackup(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/TriggerExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) UpdateExternalBackup(ctx context.Context, in *UpdateExternalBackupRequest, opts ...grpc.CallOption) (*storage.ExternalBackup, error) {
	out := new(storage.ExternalBackup)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/UpdateExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *externalBackupServiceClient) TestUpdatedExternalBackup(ctx context.Context, in *UpdateExternalBackupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.ExternalBackupService/TestUpdatedExternalBackup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExternalBackupServiceServer is the server API for ExternalBackupService service.
// All implementations must embed UnimplementedExternalBackupServiceServer
// for forward compatibility
type ExternalBackupServiceServer interface {
	// GetExternalBackup returns the external backup configuration given its ID.
	GetExternalBackup(context.Context, *ResourceByID) (*storage.ExternalBackup, error)
	// GetExternalBackups returns all external backup configurations.
	GetExternalBackups(context.Context, *Empty) (*GetExternalBackupsResponse, error)
	// PostExternalBackup creates an external backup configuration.
	PostExternalBackup(context.Context, *storage.ExternalBackup) (*storage.ExternalBackup, error)
	// PutExternalBackup modifies a given external backup, without using stored credential reconciliation.
	PutExternalBackup(context.Context, *storage.ExternalBackup) (*storage.ExternalBackup, error)
	// TestExternalBackup tests an external backup configuration.
	TestExternalBackup(context.Context, *storage.ExternalBackup) (*Empty, error)
	// DeleteExternalBackup removes an external backup configuration given its ID.
	DeleteExternalBackup(context.Context, *ResourceByID) (*Empty, error)
	// TriggerExternalBackup initiates an external backup for the given configuration.
	TriggerExternalBackup(context.Context, *ResourceByID) (*Empty, error)
	// UpdateExternalBackup modifies a given external backup, with optional stored credential reconciliation.
	UpdateExternalBackup(context.Context, *UpdateExternalBackupRequest) (*storage.ExternalBackup, error)
	// TestUpdatedExternalBackup checks if the given external backup is correctly configured, with optional stored credential reconciliation.
	TestUpdatedExternalBackup(context.Context, *UpdateExternalBackupRequest) (*Empty, error)
	mustEmbedUnimplementedExternalBackupServiceServer()
}

// UnimplementedExternalBackupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExternalBackupServiceServer struct {
}

func (UnimplementedExternalBackupServiceServer) GetExternalBackup(context.Context, *ResourceByID) (*storage.ExternalBackup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) GetExternalBackups(context.Context, *Empty) (*GetExternalBackupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExternalBackups not implemented")
}
func (UnimplementedExternalBackupServiceServer) PostExternalBackup(context.Context, *storage.ExternalBackup) (*storage.ExternalBackup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) PutExternalBackup(context.Context, *storage.ExternalBackup) (*storage.ExternalBackup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) TestExternalBackup(context.Context, *storage.ExternalBackup) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) DeleteExternalBackup(context.Context, *ResourceByID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) TriggerExternalBackup(context.Context, *ResourceByID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) UpdateExternalBackup(context.Context, *UpdateExternalBackupRequest) (*storage.ExternalBackup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) TestUpdatedExternalBackup(context.Context, *UpdateExternalBackupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestUpdatedExternalBackup not implemented")
}
func (UnimplementedExternalBackupServiceServer) mustEmbedUnimplementedExternalBackupServiceServer() {}

// UnsafeExternalBackupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExternalBackupServiceServer will
// result in compilation errors.
type UnsafeExternalBackupServiceServer interface {
	mustEmbedUnimplementedExternalBackupServiceServer()
}

func RegisterExternalBackupServiceServer(s grpc.ServiceRegistrar, srv ExternalBackupServiceServer) {
	s.RegisterService(&ExternalBackupService_ServiceDesc, srv)
}

func _ExternalBackupService_GetExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceByID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).GetExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/GetExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).GetExternalBackup(ctx, req.(*ResourceByID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_GetExternalBackups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).GetExternalBackups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/GetExternalBackups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).GetExternalBackups(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_PostExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(storage.ExternalBackup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).PostExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/PostExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).PostExternalBackup(ctx, req.(*storage.ExternalBackup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_PutExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(storage.ExternalBackup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).PutExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/PutExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).PutExternalBackup(ctx, req.(*storage.ExternalBackup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_TestExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(storage.ExternalBackup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).TestExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/TestExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).TestExternalBackup(ctx, req.(*storage.ExternalBackup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_DeleteExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceByID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).DeleteExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/DeleteExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).DeleteExternalBackup(ctx, req.(*ResourceByID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_TriggerExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceByID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).TriggerExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/TriggerExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).TriggerExternalBackup(ctx, req.(*ResourceByID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_UpdateExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateExternalBackupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).UpdateExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/UpdateExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).UpdateExternalBackup(ctx, req.(*UpdateExternalBackupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExternalBackupService_TestUpdatedExternalBackup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateExternalBackupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalBackupServiceServer).TestUpdatedExternalBackup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ExternalBackupService/TestUpdatedExternalBackup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalBackupServiceServer).TestUpdatedExternalBackup(ctx, req.(*UpdateExternalBackupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExternalBackupService_ServiceDesc is the grpc.ServiceDesc for ExternalBackupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExternalBackupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ExternalBackupService",
	HandlerType: (*ExternalBackupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetExternalBackup",
			Handler:    _ExternalBackupService_GetExternalBackup_Handler,
		},
		{
			MethodName: "GetExternalBackups",
			Handler:    _ExternalBackupService_GetExternalBackups_Handler,
		},
		{
			MethodName: "PostExternalBackup",
			Handler:    _ExternalBackupService_PostExternalBackup_Handler,
		},
		{
			MethodName: "PutExternalBackup",
			Handler:    _ExternalBackupService_PutExternalBackup_Handler,
		},
		{
			MethodName: "TestExternalBackup",
			Handler:    _ExternalBackupService_TestExternalBackup_Handler,
		},
		{
			MethodName: "DeleteExternalBackup",
			Handler:    _ExternalBackupService_DeleteExternalBackup_Handler,
		},
		{
			MethodName: "TriggerExternalBackup",
			Handler:    _ExternalBackupService_TriggerExternalBackup_Handler,
		},
		{
			MethodName: "UpdateExternalBackup",
			Handler:    _ExternalBackupService_UpdateExternalBackup_Handler,
		},
		{
			MethodName: "TestUpdatedExternalBackup",
			Handler:    _ExternalBackupService_TestUpdatedExternalBackup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/backup_service.proto",
}

func (m *GetExternalBackupsResponse) MarshalVT() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVT(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetExternalBackupsResponse) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *GetExternalBackupsResponse) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if len(m.ExternalBackups) > 0 {
		for iNdEx := len(m.ExternalBackups) - 1; iNdEx >= 0; iNdEx-- {
			if vtmsg, ok := interface{}(m.ExternalBackups[iNdEx]).(interface {
				MarshalToSizedBufferVT([]byte) (int, error)
			}); ok {
				size, err := vtmsg.MarshalToSizedBufferVT(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarint(dAtA, i, uint64(size))
			} else {
				encoded, err := proto.Marshal(m.ExternalBackups[iNdEx])
				if err != nil {
					return 0, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = encodeVarint(dAtA, i, uint64(len(encoded)))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *UpdateExternalBackupRequest) MarshalVT() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVT(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateExternalBackupRequest) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *UpdateExternalBackupRequest) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if m.UpdatePassword {
		i--
		if m.UpdatePassword {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.ExternalBackup != nil {
		if vtmsg, ok := interface{}(m.ExternalBackup).(interface {
			MarshalToSizedBufferVT([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVT(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.ExternalBackup)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GetExternalBackupsResponse) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ExternalBackups) > 0 {
		for _, e := range m.ExternalBackups {
			if size, ok := interface{}(e).(interface {
				SizeVT() int
			}); ok {
				l = size.SizeVT()
			} else {
				l = proto.Size(e)
			}
			n += 1 + l + sov(uint64(l))
		}
	}
	n += len(m.unknownFields)
	return n
}

func (m *UpdateExternalBackupRequest) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ExternalBackup != nil {
		if size, ok := interface{}(m.ExternalBackup).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.ExternalBackup)
		}
		n += 1 + l + sov(uint64(l))
	}
	if m.UpdatePassword {
		n += 2
	}
	n += len(m.unknownFields)
	return n
}

func (m *GetExternalBackupsResponse) UnmarshalVT(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetExternalBackupsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetExternalBackupsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalBackups", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExternalBackups = append(m.ExternalBackups, &storage.ExternalBackup{})
			if unmarshal, ok := interface{}(m.ExternalBackups[len(m.ExternalBackups)-1]).(interface {
				UnmarshalVT([]byte) error
			}); ok {
				if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				if err := proto.Unmarshal(dAtA[iNdEx:postIndex], m.ExternalBackups[len(m.ExternalBackups)-1]); err != nil {
					return err
				}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateExternalBackupRequest) UnmarshalVT(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpdateExternalBackupRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateExternalBackupRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalBackup", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ExternalBackup == nil {
				m.ExternalBackup = &storage.ExternalBackup{}
			}
			if unmarshal, ok := interface{}(m.ExternalBackup).(interface {
				UnmarshalVT([]byte) error
			}); ok {
				if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				if err := proto.Unmarshal(dAtA[iNdEx:postIndex], m.ExternalBackup); err != nil {
					return err
				}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatePassword", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.UpdatePassword = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}