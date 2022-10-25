// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// protoc-gen-go-vtproto version: v0.3.1-0.20220817155510-0ae748fd2007
// source: internalapi/sensor/network_connection_iservice.proto

package sensor

import (
	context "context"
	binary "encoding/binary"
	fmt "fmt"
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

func (m *NetworkConnectionInfoMessage) CloneVT() *NetworkConnectionInfoMessage {
	if m == nil {
		return (*NetworkConnectionInfoMessage)(nil)
	}
	r := &NetworkConnectionInfoMessage{}
	if m.Msg != nil {
		r.Msg = m.Msg.(interface {
			CloneVT() isNetworkConnectionInfoMessage_Msg
		}).CloneVT()
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *NetworkConnectionInfoMessage) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (m *NetworkConnectionInfoMessage_Register) CloneVT() isNetworkConnectionInfoMessage_Msg {
	if m == nil {
		return (*NetworkConnectionInfoMessage_Register)(nil)
	}
	r := &NetworkConnectionInfoMessage_Register{
		Register: m.Register.CloneVT(),
	}
	return r
}

func (m *NetworkConnectionInfoMessage_Info) CloneVT() isNetworkConnectionInfoMessage_Msg {
	if m == nil {
		return (*NetworkConnectionInfoMessage_Info)(nil)
	}
	r := &NetworkConnectionInfoMessage_Info{
		Info: m.Info.CloneVT(),
	}
	return r
}

func (m *NetworkFlowsControlMessage) CloneVT() *NetworkFlowsControlMessage {
	if m == nil {
		return (*NetworkFlowsControlMessage)(nil)
	}
	r := &NetworkFlowsControlMessage{
		PublicIpAddresses: m.PublicIpAddresses.CloneVT(),
		IpNetworks:        m.IpNetworks.CloneVT(),
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *NetworkFlowsControlMessage) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (m *IPAddressList) CloneVT() *IPAddressList {
	if m == nil {
		return (*IPAddressList)(nil)
	}
	r := &IPAddressList{}
	if rhs := m.Ipv4Addresses; rhs != nil {
		tmpContainer := make([]uint32, len(rhs))
		copy(tmpContainer, rhs)
		r.Ipv4Addresses = tmpContainer
	}
	if rhs := m.Ipv6Addresses; rhs != nil {
		tmpContainer := make([]uint64, len(rhs))
		copy(tmpContainer, rhs)
		r.Ipv6Addresses = tmpContainer
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *IPAddressList) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (m *IPNetworkList) CloneVT() *IPNetworkList {
	if m == nil {
		return (*IPNetworkList)(nil)
	}
	r := &IPNetworkList{}
	if rhs := m.Ipv4Networks; rhs != nil {
		tmpBytes := make([]byte, len(rhs))
		copy(tmpBytes, rhs)
		r.Ipv4Networks = tmpBytes
	}
	if rhs := m.Ipv6Networks; rhs != nil {
		tmpBytes := make([]byte, len(rhs))
		copy(tmpBytes, rhs)
		r.Ipv6Networks = tmpBytes
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *IPNetworkList) CloneGenericVT() proto.Message {
	return m.CloneVT()
}

func (this *NetworkConnectionInfoMessage) EqualVT(that *NetworkConnectionInfoMessage) bool {
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
			EqualVT(isNetworkConnectionInfoMessage_Msg) bool
		}).EqualVT(that.Msg) {
			return false
		}
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *NetworkConnectionInfoMessage_Register) EqualVT(thatIface isNetworkConnectionInfoMessage_Msg) bool {
	that, ok := thatIface.(*NetworkConnectionInfoMessage_Register)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.Register, that.Register; p != q {
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

func (this *NetworkConnectionInfoMessage_Info) EqualVT(thatIface isNetworkConnectionInfoMessage_Msg) bool {
	that, ok := thatIface.(*NetworkConnectionInfoMessage_Info)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.Info, that.Info; p != q {
		if p == nil {
			p = &NetworkConnectionInfo{}
		}
		if q == nil {
			q = &NetworkConnectionInfo{}
		}
		if !p.EqualVT(q) {
			return false
		}
	}
	return true
}

func (this *NetworkFlowsControlMessage) EqualVT(that *NetworkFlowsControlMessage) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if !this.PublicIpAddresses.EqualVT(that.PublicIpAddresses) {
		return false
	}
	if !this.IpNetworks.EqualVT(that.IpNetworks) {
		return false
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *IPAddressList) EqualVT(that *IPAddressList) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if len(this.Ipv4Addresses) != len(that.Ipv4Addresses) {
		return false
	}
	for i, vx := range this.Ipv4Addresses {
		vy := that.Ipv4Addresses[i]
		if vx != vy {
			return false
		}
	}
	if len(this.Ipv6Addresses) != len(that.Ipv6Addresses) {
		return false
	}
	for i, vx := range this.Ipv6Addresses {
		vy := that.Ipv6Addresses[i]
		if vx != vy {
			return false
		}
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *IPNetworkList) EqualVT(that *IPNetworkList) bool {
	if this == nil {
		return that == nil
	} else if that == nil {
		return false
	}
	if string(this.Ipv4Networks) != string(that.Ipv4Networks) {
		return false
	}
	if string(this.Ipv6Networks) != string(that.Ipv6Networks) {
		return false
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NetworkConnectionInfoServiceClient is the client API for NetworkConnectionInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkConnectionInfoServiceClient interface {
	// Note: the response is a stream due to a bug in the C++ GRPC client library. The server is not expected to
	// send anything via this stream.
	PushNetworkConnectionInfo(ctx context.Context, opts ...grpc.CallOption) (NetworkConnectionInfoService_PushNetworkConnectionInfoClient, error)
}

type networkConnectionInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkConnectionInfoServiceClient(cc grpc.ClientConnInterface) NetworkConnectionInfoServiceClient {
	return &networkConnectionInfoServiceClient{cc}
}

func (c *networkConnectionInfoServiceClient) PushNetworkConnectionInfo(ctx context.Context, opts ...grpc.CallOption) (NetworkConnectionInfoService_PushNetworkConnectionInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &NetworkConnectionInfoService_ServiceDesc.Streams[0], "/sensor.NetworkConnectionInfoService/PushNetworkConnectionInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &networkConnectionInfoServicePushNetworkConnectionInfoClient{stream}
	return x, nil
}

type NetworkConnectionInfoService_PushNetworkConnectionInfoClient interface {
	Send(*NetworkConnectionInfoMessage) error
	Recv() (*NetworkFlowsControlMessage, error)
	grpc.ClientStream
}

type networkConnectionInfoServicePushNetworkConnectionInfoClient struct {
	grpc.ClientStream
}

func (x *networkConnectionInfoServicePushNetworkConnectionInfoClient) Send(m *NetworkConnectionInfoMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *networkConnectionInfoServicePushNetworkConnectionInfoClient) Recv() (*NetworkFlowsControlMessage, error) {
	m := new(NetworkFlowsControlMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NetworkConnectionInfoServiceServer is the server API for NetworkConnectionInfoService service.
// All implementations must embed UnimplementedNetworkConnectionInfoServiceServer
// for forward compatibility
type NetworkConnectionInfoServiceServer interface {
	// Note: the response is a stream due to a bug in the C++ GRPC client library. The server is not expected to
	// send anything via this stream.
	PushNetworkConnectionInfo(NetworkConnectionInfoService_PushNetworkConnectionInfoServer) error
	mustEmbedUnimplementedNetworkConnectionInfoServiceServer()
}

// UnimplementedNetworkConnectionInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNetworkConnectionInfoServiceServer struct {
}

func (UnimplementedNetworkConnectionInfoServiceServer) PushNetworkConnectionInfo(NetworkConnectionInfoService_PushNetworkConnectionInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method PushNetworkConnectionInfo not implemented")
}
func (UnimplementedNetworkConnectionInfoServiceServer) mustEmbedUnimplementedNetworkConnectionInfoServiceServer() {
}

// UnsafeNetworkConnectionInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkConnectionInfoServiceServer will
// result in compilation errors.
type UnsafeNetworkConnectionInfoServiceServer interface {
	mustEmbedUnimplementedNetworkConnectionInfoServiceServer()
}

func RegisterNetworkConnectionInfoServiceServer(s grpc.ServiceRegistrar, srv NetworkConnectionInfoServiceServer) {
	s.RegisterService(&NetworkConnectionInfoService_ServiceDesc, srv)
}

func _NetworkConnectionInfoService_PushNetworkConnectionInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NetworkConnectionInfoServiceServer).PushNetworkConnectionInfo(&networkConnectionInfoServicePushNetworkConnectionInfoServer{stream})
}

type NetworkConnectionInfoService_PushNetworkConnectionInfoServer interface {
	Send(*NetworkFlowsControlMessage) error
	Recv() (*NetworkConnectionInfoMessage, error)
	grpc.ServerStream
}

type networkConnectionInfoServicePushNetworkConnectionInfoServer struct {
	grpc.ServerStream
}

func (x *networkConnectionInfoServicePushNetworkConnectionInfoServer) Send(m *NetworkFlowsControlMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *networkConnectionInfoServicePushNetworkConnectionInfoServer) Recv() (*NetworkConnectionInfoMessage, error) {
	m := new(NetworkConnectionInfoMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NetworkConnectionInfoService_ServiceDesc is the grpc.ServiceDesc for NetworkConnectionInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkConnectionInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sensor.NetworkConnectionInfoService",
	HandlerType: (*NetworkConnectionInfoServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushNetworkConnectionInfo",
			Handler:       _NetworkConnectionInfoService_PushNetworkConnectionInfo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internalapi/sensor/network_connection_iservice.proto",
}

func (m *NetworkConnectionInfoMessage) MarshalVT() (dAtA []byte, err error) {
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

func (m *NetworkConnectionInfoMessage) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *NetworkConnectionInfoMessage) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
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

func (m *NetworkConnectionInfoMessage_Register) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *NetworkConnectionInfoMessage_Register) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Register != nil {
		size, err := m.Register.MarshalToSizedBufferVT(dAtA[:i])
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
func (m *NetworkConnectionInfoMessage_Info) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *NetworkConnectionInfoMessage_Info) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Info != nil {
		size, err := m.Info.MarshalToSizedBufferVT(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *NetworkFlowsControlMessage) MarshalVT() (dAtA []byte, err error) {
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

func (m *NetworkFlowsControlMessage) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *NetworkFlowsControlMessage) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
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
	if m.IpNetworks != nil {
		size, err := m.IpNetworks.MarshalToSizedBufferVT(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	if m.PublicIpAddresses != nil {
		size, err := m.PublicIpAddresses.MarshalToSizedBufferVT(dAtA[:i])
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

func (m *IPAddressList) MarshalVT() (dAtA []byte, err error) {
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

func (m *IPAddressList) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *IPAddressList) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
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
	if len(m.Ipv6Addresses) > 0 {
		for iNdEx := len(m.Ipv6Addresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= 8
			binary.LittleEndian.PutUint64(dAtA[i:], uint64(m.Ipv6Addresses[iNdEx]))
		}
		i = encodeVarint(dAtA, i, uint64(len(m.Ipv6Addresses)*8))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Ipv4Addresses) > 0 {
		for iNdEx := len(m.Ipv4Addresses) - 1; iNdEx >= 0; iNdEx-- {
			i -= 4
			binary.LittleEndian.PutUint32(dAtA[i:], uint32(m.Ipv4Addresses[iNdEx]))
		}
		i = encodeVarint(dAtA, i, uint64(len(m.Ipv4Addresses)*4))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IPNetworkList) MarshalVT() (dAtA []byte, err error) {
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

func (m *IPNetworkList) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *IPNetworkList) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
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
	if len(m.Ipv6Networks) > 0 {
		i -= len(m.Ipv6Networks)
		copy(dAtA[i:], m.Ipv6Networks)
		i = encodeVarint(dAtA, i, uint64(len(m.Ipv6Networks)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Ipv4Networks) > 0 {
		i -= len(m.Ipv4Networks)
		copy(dAtA[i:], m.Ipv4Networks)
		i = encodeVarint(dAtA, i, uint64(len(m.Ipv4Networks)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *NetworkConnectionInfoMessage) SizeVT() (n int) {
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

func (m *NetworkConnectionInfoMessage_Register) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Register != nil {
		l = m.Register.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *NetworkConnectionInfoMessage_Info) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Info != nil {
		l = m.Info.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *NetworkFlowsControlMessage) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PublicIpAddresses != nil {
		l = m.PublicIpAddresses.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	if m.IpNetworks != nil {
		l = m.IpNetworks.SizeVT()
		n += 1 + l + sov(uint64(l))
	}
	n += len(m.unknownFields)
	return n
}

func (m *IPAddressList) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ipv4Addresses) > 0 {
		n += 1 + sov(uint64(len(m.Ipv4Addresses)*4)) + len(m.Ipv4Addresses)*4
	}
	if len(m.Ipv6Addresses) > 0 {
		n += 1 + sov(uint64(len(m.Ipv6Addresses)*8)) + len(m.Ipv6Addresses)*8
	}
	n += len(m.unknownFields)
	return n
}

func (m *IPNetworkList) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Ipv4Networks)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Ipv6Networks)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	n += len(m.unknownFields)
	return n
}

func (m *NetworkConnectionInfoMessage) UnmarshalVT(dAtA []byte) error {
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
			return fmt.Errorf("proto: NetworkConnectionInfoMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NetworkConnectionInfoMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Register", wireType)
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
			if oneof, ok := m.Msg.(*NetworkConnectionInfoMessage_Register); ok {
				if err := oneof.Register.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				v := &CollectorRegisterRequest{}
				if err := v.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
				m.Msg = &NetworkConnectionInfoMessage_Register{Register: v}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Info", wireType)
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
			if oneof, ok := m.Msg.(*NetworkConnectionInfoMessage_Info); ok {
				if err := oneof.Info.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				v := &NetworkConnectionInfo{}
				if err := v.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
				m.Msg = &NetworkConnectionInfoMessage_Info{Info: v}
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
func (m *NetworkFlowsControlMessage) UnmarshalVT(dAtA []byte) error {
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
			return fmt.Errorf("proto: NetworkFlowsControlMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NetworkFlowsControlMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublicIpAddresses", wireType)
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
			if m.PublicIpAddresses == nil {
				m.PublicIpAddresses = &IPAddressList{}
			}
			if err := m.PublicIpAddresses.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IpNetworks", wireType)
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
			if m.IpNetworks == nil {
				m.IpNetworks = &IPNetworkList{}
			}
			if err := m.IpNetworks.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
				return err
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
func (m *IPAddressList) UnmarshalVT(dAtA []byte) error {
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
			return fmt.Errorf("proto: IPAddressList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IPAddressList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 5 {
				var v uint32
				if (iNdEx + 4) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint32(binary.LittleEndian.Uint32(dAtA[iNdEx:]))
				iNdEx += 4
				m.Ipv4Addresses = append(m.Ipv4Addresses, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflow
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLength
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLength
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 4
				if elementCount != 0 && len(m.Ipv4Addresses) == 0 {
					m.Ipv4Addresses = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					if (iNdEx + 4) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint32(binary.LittleEndian.Uint32(dAtA[iNdEx:]))
					iNdEx += 4
					m.Ipv4Addresses = append(m.Ipv4Addresses, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Ipv4Addresses", wireType)
			}
		case 2:
			if wireType == 1 {
				var v uint64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
				iNdEx += 8
				m.Ipv6Addresses = append(m.Ipv6Addresses, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflow
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLength
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLength
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 8
				if elementCount != 0 && len(m.Ipv6Addresses) == 0 {
					m.Ipv6Addresses = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
					iNdEx += 8
					m.Ipv6Addresses = append(m.Ipv6Addresses, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Ipv6Addresses", wireType)
			}
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
func (m *IPNetworkList) UnmarshalVT(dAtA []byte) error {
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
			return fmt.Errorf("proto: IPNetworkList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IPNetworkList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ipv4Networks", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ipv4Networks = append(m.Ipv4Networks[:0], dAtA[iNdEx:postIndex]...)
			if m.Ipv4Networks == nil {
				m.Ipv4Networks = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ipv6Networks", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ipv6Networks = append(m.Ipv6Networks[:0], dAtA[iNdEx:postIndex]...)
			if m.Ipv6Networks == nil {
				m.Ipv6Networks = []byte{}
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