// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// protoc-gen-go-vtproto version: v0.3.1-0.20220817155510-0ae748fd2007
// source: internalapi/sensor/signal_iservice.proto

package sensor

import (
	context "context"
	fmt "fmt"
	v1 "github.com/stackrox/rox/generated/api/v1"
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

func (m *SignalStreamMessage) CloneVT() *SignalStreamMessage {
	if m == nil {
		return (*SignalStreamMessage)(nil)
	}
	r := &SignalStreamMessage{}
	if m.Msg != nil {
		r.Msg = m.Msg.(interface {
			CloneVT() isSignalStreamMessage_Msg
		}).CloneVT()
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *SignalStreamMessage) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (m *SignalStreamMessage_CollectorRegisterRequest) CloneVT() isSignalStreamMessage_Msg {
	if m == nil {
		return (*SignalStreamMessage_CollectorRegisterRequest)(nil)
	}
	r := &SignalStreamMessage_CollectorRegisterRequest{
		CollectorRegisterRequest: m.CollectorRegisterRequest.CloneVT(),
	}
	return r
}

func (m *SignalStreamMessage_Signal) CloneVT() isSignalStreamMessage_Msg {
	if m == nil {
		return (*SignalStreamMessage_Signal)(nil)
	}
	r := &SignalStreamMessage_Signal{}
	if rhs := m.Signal; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface{ CloneVT() *v1.Signal }); ok {
			r.Signal = vtpb.CloneVT()
		} else {
			r.Signal = proto.Clone(rhs).(*v1.Signal)
		}
	}
	return r
}

func (this *SignalStreamMessage) EqualVT(that *SignalStreamMessage) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if this.Msg == nil && that.Msg != nil {
		return false
	} else if this.Msg != nil {
		if that.Msg == nil {
			return false
		}
		if !this.Msg.(interface {
			EqualVT(isSignalStreamMessage_Msg) bool
		}).EqualVT(that.Msg) {
			return false
		}
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *SignalStreamMessage_CollectorRegisterRequest) EqualVT(thatIface isSignalStreamMessage_Msg) bool {
	that, ok := thatIface.(*SignalStreamMessage_CollectorRegisterRequest)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.CollectorRegisterRequest, that.CollectorRegisterRequest; p != q {
		if p == nil {
			p = &CollectorRegisterRequest{}
		}
		if q == nil {
			q = &CollectorRegisterRequest{}
		}
		if !p.EqualVT(q) {
			return false
		}
	}
	return true
}

func (this *SignalStreamMessage_Signal) EqualVT(thatIface isSignalStreamMessage_Msg) bool {
	that, ok := thatIface.(*SignalStreamMessage_Signal)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.Signal, that.Signal; p != q {
		if p == nil {
			p = &v1.Signal{}
		}
		if q == nil {
			q = &v1.Signal{}
		}
		if equal, ok := interface{}(p).(interface{ EqualVT(*v1.Signal) bool }); ok {
			if !equal.EqualVT(q) {
				return false
			}
		} else if !proto.Equal(p, q) {
			return false
		}
	}
	return true
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SignalServiceClient is the client API for SignalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SignalServiceClient interface {
	// Note: the response is a stream due to a bug in the C++ GRPC client library. The server is not expected to
	// send anything via this stream.
	PushSignals(ctx context.Context, opts ...grpc.CallOption) (SignalService_PushSignalsClient, error)
}

type signalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSignalServiceClient(cc grpc.ClientConnInterface) SignalServiceClient {
	return &signalServiceClient{cc}
}

func (c *signalServiceClient) PushSignals(ctx context.Context, opts ...grpc.CallOption) (SignalService_PushSignalsClient, error) {
	stream, err := c.cc.NewStream(ctx, &SignalService_ServiceDesc.Streams[0], "/sensor.SignalService/PushSignals", opts...)
	if err != nil {
		return nil, err
	}
	x := &signalServicePushSignalsClient{stream}
	return x, nil
}

type SignalService_PushSignalsClient interface {
	Send(*SignalStreamMessage) error
	Recv() (*v1.Empty, error)
	grpc.ClientStream
}

type signalServicePushSignalsClient struct {
	grpc.ClientStream
}

func (x *signalServicePushSignalsClient) Send(m *SignalStreamMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *signalServicePushSignalsClient) Recv() (*v1.Empty, error) {
	m := new(v1.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SignalServiceServer is the server API for SignalService service.
// All implementations must embed UnimplementedSignalServiceServer
// for forward compatibility
type SignalServiceServer interface {
	// Note: the response is a stream due to a bug in the C++ GRPC client library. The server is not expected to
	// send anything via this stream.
	PushSignals(SignalService_PushSignalsServer) error
	mustEmbedUnimplementedSignalServiceServer()
}

// UnimplementedSignalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSignalServiceServer struct {
}

func (UnimplementedSignalServiceServer) PushSignals(SignalService_PushSignalsServer) error {
	return status.Errorf(codes.Unimplemented, "method PushSignals not implemented")
}
func (UnimplementedSignalServiceServer) mustEmbedUnimplementedSignalServiceServer() {}

// UnsafeSignalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SignalServiceServer will
// result in compilation errors.
type UnsafeSignalServiceServer interface {
	mustEmbedUnimplementedSignalServiceServer()
}

func RegisterSignalServiceServer(s grpc.ServiceRegistrar, srv SignalServiceServer) {
	s.RegisterService(&SignalService_ServiceDesc, srv)
}

func _SignalService_PushSignals_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SignalServiceServer).PushSignals(&signalServicePushSignalsServer{stream})
}

type SignalService_PushSignalsServer interface {
	Send(*v1.Empty) error
	Recv() (*SignalStreamMessage, error)
	grpc.ServerStream
}

type signalServicePushSignalsServer struct {
	grpc.ServerStream
}

func (x *signalServicePushSignalsServer) Send(m *v1.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *signalServicePushSignalsServer) Recv() (*SignalStreamMessage, error) {
	m := new(SignalStreamMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SignalService_ServiceDesc is the grpc.ServiceDesc for SignalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SignalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sensor.SignalService",
	HandlerType: (*SignalServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushSignals",
			Handler:       _SignalService_PushSignals_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internalapi/sensor/signal_iservice.proto",
}

func (m *SignalStreamMessage) MarshalVT() (dAtA []byte, err error) {
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

func (m *SignalStreamMessage) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *SignalStreamMessage) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
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
	if vtmsg, ok := m.Msg.(interface {
		MarshalToSizedBufferVT([]byte) (int, error)
	}); ok {
		size, err := vtmsg.MarshalToSizedBufferVT(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	return len(dAtA) - i, nil
}

func (m *SignalStreamMessage_CollectorRegisterRequest) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *SignalStreamMessage_CollectorRegisterRequest) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CollectorRegisterRequest != nil {
		size, err := m.CollectorRegisterRequest.MarshalToSizedBufferVT(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *SignalStreamMessage_Signal) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *SignalStreamMessage_Signal) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Signal != nil {
		if vtmsg, ok := interface{}(m.Signal).(interface {
			MarshalToSizedBufferVT([]byte) (int, error)
		}); ok {
			size, err := vtmsg.MarshalToSizedBufferVT(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarint(dAtA, i, uint64(size))
		} else {
			encoded, err := proto.Marshal(m.Signal)
			if err != nil {
				return 0, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = encodeVarint(dAtA, i, uint64(len(encoded)))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *SignalStreamMessage) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if vtmsg, ok := m.Msg.(interface{ SizeVT() int }); ok {
		n += vtmsg.SizeVT()
	}
	n += len(m.unknownFields)
	return n
}

func (m *SignalStreamMessage_CollectorRegisterRequest) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CollectorRegisterRequest != nil {
		l = m.CollectorRegisterRequest.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *SignalStreamMessage_Signal) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Signal != nil {
		if size, ok := interface{}(m.Signal).(interface {
			SizeVT() int
		}); ok {
			l = size.SizeVT()
		} else {
			l = proto.Size(m.Signal)
		}
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *SignalStreamMessage) UnmarshalVT(dAtA []byte) error {
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
			return fmt.Errorf("proto: SignalStreamMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignalStreamMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollectorRegisterRequest", wireType)
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
			if oneof, ok := m.Msg.(*SignalStreamMessage_CollectorRegisterRequest); ok {
				if err := oneof.CollectorRegisterRequest.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				v := &CollectorRegisterRequest{}
				if err := v.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
				m.Msg = &SignalStreamMessage_CollectorRegisterRequest{CollectorRegisterRequest: v}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signal", wireType)
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
			if oneof, ok := m.Msg.(*SignalStreamMessage_Signal); ok {
				if unmarshal, ok := interface{}(oneof.Signal).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], oneof.Signal); err != nil {
						return err
					}
				}
			} else {
				v := &v1.Signal{}
				if unmarshal, ok := interface{}(v).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], v); err != nil {
						return err
					}
				}
				m.Msg = &SignalStreamMessage_Signal{Signal: v}
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