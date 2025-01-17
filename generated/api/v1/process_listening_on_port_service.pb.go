// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/process_listening_on_port_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	storage "github.com/stackrox/rox/generated/storage"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
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

type GetProcessesListeningOnPortsRequest struct {
	DeploymentId         string   `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProcessesListeningOnPortsRequest) Reset()         { *m = GetProcessesListeningOnPortsRequest{} }
func (m *GetProcessesListeningOnPortsRequest) String() string { return proto.CompactTextString(m) }
func (*GetProcessesListeningOnPortsRequest) ProtoMessage()    {}
func (*GetProcessesListeningOnPortsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8bf40985da37317, []int{0}
}
func (m *GetProcessesListeningOnPortsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetProcessesListeningOnPortsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetProcessesListeningOnPortsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetProcessesListeningOnPortsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProcessesListeningOnPortsRequest.Merge(m, src)
}
func (m *GetProcessesListeningOnPortsRequest) XXX_Size() int {
	return m.Size()
}
func (m *GetProcessesListeningOnPortsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProcessesListeningOnPortsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetProcessesListeningOnPortsRequest proto.InternalMessageInfo

func (m *GetProcessesListeningOnPortsRequest) GetDeploymentId() string {
	if m != nil {
		return m.DeploymentId
	}
	return ""
}

func (m *GetProcessesListeningOnPortsRequest) MessageClone() proto.Message {
	return m.Clone()
}
func (m *GetProcessesListeningOnPortsRequest) Clone() *GetProcessesListeningOnPortsRequest {
	if m == nil {
		return nil
	}
	cloned := new(GetProcessesListeningOnPortsRequest)
	*cloned = *m

	return cloned
}

type GetProcessesListeningOnPortsResponse struct {
	ListeningEndpoints   []*storage.ProcessListeningOnPort `protobuf:"bytes,1,rep,name=listening_endpoints,json=listeningEndpoints,proto3" json:"listening_endpoints,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *GetProcessesListeningOnPortsResponse) Reset()         { *m = GetProcessesListeningOnPortsResponse{} }
func (m *GetProcessesListeningOnPortsResponse) String() string { return proto.CompactTextString(m) }
func (*GetProcessesListeningOnPortsResponse) ProtoMessage()    {}
func (*GetProcessesListeningOnPortsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8bf40985da37317, []int{1}
}
func (m *GetProcessesListeningOnPortsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetProcessesListeningOnPortsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetProcessesListeningOnPortsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetProcessesListeningOnPortsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProcessesListeningOnPortsResponse.Merge(m, src)
}
func (m *GetProcessesListeningOnPortsResponse) XXX_Size() int {
	return m.Size()
}
func (m *GetProcessesListeningOnPortsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProcessesListeningOnPortsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetProcessesListeningOnPortsResponse proto.InternalMessageInfo

func (m *GetProcessesListeningOnPortsResponse) GetListeningEndpoints() []*storage.ProcessListeningOnPort {
	if m != nil {
		return m.ListeningEndpoints
	}
	return nil
}

func (m *GetProcessesListeningOnPortsResponse) MessageClone() proto.Message {
	return m.Clone()
}
func (m *GetProcessesListeningOnPortsResponse) Clone() *GetProcessesListeningOnPortsResponse {
	if m == nil {
		return nil
	}
	cloned := new(GetProcessesListeningOnPortsResponse)
	*cloned = *m

	if m.ListeningEndpoints != nil {
		cloned.ListeningEndpoints = make([]*storage.ProcessListeningOnPort, len(m.ListeningEndpoints))
		for idx, v := range m.ListeningEndpoints {
			cloned.ListeningEndpoints[idx] = v.Clone()
		}
	}
	return cloned
}

func init() {
	proto.RegisterType((*GetProcessesListeningOnPortsRequest)(nil), "v1.GetProcessesListeningOnPortsRequest")
	proto.RegisterType((*GetProcessesListeningOnPortsResponse)(nil), "v1.GetProcessesListeningOnPortsResponse")
}

func init() {
	proto.RegisterFile("api/v1/process_listening_on_port_service.proto", fileDescriptor_d8bf40985da37317)
}

var fileDescriptor_d8bf40985da37317 = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4a, 0x33, 0x31,
	0x14, 0x85, 0x9b, 0xfe, 0xf0, 0x83, 0x51, 0x37, 0x71, 0x53, 0x4a, 0x19, 0x4b, 0x2b, 0xb4, 0xab,
	0x0c, 0x53, 0x5d, 0xb9, 0x14, 0x44, 0x14, 0xc1, 0x52, 0x37, 0xe2, 0x66, 0x18, 0x3b, 0x97, 0x21,
	0x58, 0x73, 0x63, 0x6e, 0x1c, 0x2a, 0xe2, 0xc6, 0x57, 0x70, 0xe3, 0x4b, 0xf8, 0x1e, 0x2e, 0x45,
	0x5f, 0x40, 0xaa, 0x0f, 0x22, 0x35, 0xad, 0x83, 0x62, 0xab, 0xdb, 0xcb, 0x77, 0x4e, 0xce, 0x3d,
	0x37, 0x5c, 0x26, 0x46, 0x85, 0x79, 0x14, 0x1a, 0x8b, 0x7d, 0x20, 0x8a, 0x07, 0x8a, 0x1c, 0x68,
	0xa5, 0xb3, 0x18, 0x75, 0x6c, 0xd0, 0xba, 0x98, 0xc0, 0xe6, 0xaa, 0x0f, 0xd2, 0x58, 0x74, 0x28,
	0xca, 0x79, 0x54, 0xad, 0x65, 0x88, 0xd9, 0x00, 0xc2, 0xb1, 0x34, 0xd1, 0x1a, 0x5d, 0xe2, 0x14,
	0x6a, 0xf2, 0x44, 0xb5, 0x45, 0x0e, 0x6d, 0x92, 0xc1, 0x6c, 0x4b, 0x0f, 0x36, 0xf6, 0x78, 0x73,
	0x07, 0x5c, 0xd7, 0x53, 0x40, 0xfb, 0x53, 0xec, 0x40, 0x77, 0xd1, 0x3a, 0xea, 0xc1, 0xf9, 0x05,
	0x90, 0x13, 0x4d, 0xbe, 0x9c, 0x82, 0x19, 0xe0, 0xe5, 0x19, 0x68, 0x17, 0xab, 0xb4, 0xc2, 0xea,
	0xac, 0xbd, 0xd0, 0x5b, 0x2a, 0x86, 0xbb, 0x69, 0x63, 0xc8, 0xd7, 0xe6, 0x7b, 0x91, 0x41, 0x4d,
	0x20, 0xba, 0x7c, 0xa5, 0x88, 0x03, 0x3a, 0x35, 0xa8, 0xb4, 0xa3, 0x0a, 0xab, 0xff, 0x6b, 0x2f,
	0x76, 0x56, 0xe5, 0x24, 0xba, 0x9c, 0x18, 0x7d, 0xb3, 0xe9, 0x89, 0x4f, 0xed, 0xf6, 0x54, 0xda,
	0x79, 0x62, 0xbc, 0x3e, 0xf3, 0xdd, 0x43, 0xdf, 0x9d, 0xb8, 0x67, 0xbc, 0x36, 0x2f, 0x9f, 0x68,
	0xc9, 0x3c, 0x92, 0x7f, 0x68, 0xa3, 0xda, 0xfe, 0x1d, 0xf4, 0xab, 0x36, 0x36, 0x6f, 0x9e, 0xdf,
	0x6e, 0xcb, 0x1b, 0xa2, 0x33, 0x3e, 0xef, 0x0f, 0x4b, 0x87, 0x45, 0x81, 0xe1, 0xd5, 0x97, 0x86,
	0xaf, 0xb7, 0xe4, 0xc3, 0x28, 0x60, 0x8f, 0xa3, 0x80, 0xbd, 0x8c, 0x02, 0x76, 0xf7, 0x1a, 0x94,
	0x78, 0x45, 0xa1, 0x24, 0x97, 0xf4, 0x4f, 0x2d, 0x0e, 0xfd, 0xfd, 0xc6, 0x3f, 0x47, 0xe6, 0xd1,
	0x71, 0x39, 0x8f, 0x8e, 0x4a, 0x27, 0xff, 0x3f, 0x66, 0xeb, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x1f, 0x0a, 0x2f, 0x3f, 0x50, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProcessesListeningOnPortsServiceClient is the client API for ProcessesListeningOnPortsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConnInterface.NewStream.
type ProcessesListeningOnPortsServiceClient interface {
	GetProcessesListeningOnPorts(ctx context.Context, in *GetProcessesListeningOnPortsRequest, opts ...grpc.CallOption) (*GetProcessesListeningOnPortsResponse, error)
}

type processesListeningOnPortsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessesListeningOnPortsServiceClient(cc grpc.ClientConnInterface) ProcessesListeningOnPortsServiceClient {
	return &processesListeningOnPortsServiceClient{cc}
}

func (c *processesListeningOnPortsServiceClient) GetProcessesListeningOnPorts(ctx context.Context, in *GetProcessesListeningOnPortsRequest, opts ...grpc.CallOption) (*GetProcessesListeningOnPortsResponse, error) {
	out := new(GetProcessesListeningOnPortsResponse)
	err := c.cc.Invoke(ctx, "/v1.ProcessesListeningOnPortsService/GetProcessesListeningOnPorts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessesListeningOnPortsServiceServer is the server API for ProcessesListeningOnPortsService service.
type ProcessesListeningOnPortsServiceServer interface {
	GetProcessesListeningOnPorts(context.Context, *GetProcessesListeningOnPortsRequest) (*GetProcessesListeningOnPortsResponse, error)
}

// UnimplementedProcessesListeningOnPortsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProcessesListeningOnPortsServiceServer struct {
}

func (*UnimplementedProcessesListeningOnPortsServiceServer) GetProcessesListeningOnPorts(ctx context.Context, req *GetProcessesListeningOnPortsRequest) (*GetProcessesListeningOnPortsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProcessesListeningOnPorts not implemented")
}

func RegisterProcessesListeningOnPortsServiceServer(s *grpc.Server, srv ProcessesListeningOnPortsServiceServer) {
	s.RegisterService(&_ProcessesListeningOnPortsService_serviceDesc, srv)
}

func _ProcessesListeningOnPortsService_GetProcessesListeningOnPorts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProcessesListeningOnPortsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessesListeningOnPortsServiceServer).GetProcessesListeningOnPorts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ProcessesListeningOnPortsService/GetProcessesListeningOnPorts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessesListeningOnPortsServiceServer).GetProcessesListeningOnPorts(ctx, req.(*GetProcessesListeningOnPortsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProcessesListeningOnPortsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ProcessesListeningOnPortsService",
	HandlerType: (*ProcessesListeningOnPortsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProcessesListeningOnPorts",
			Handler:    _ProcessesListeningOnPortsService_GetProcessesListeningOnPorts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/process_listening_on_port_service.proto",
}

func (m *GetProcessesListeningOnPortsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetProcessesListeningOnPortsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetProcessesListeningOnPortsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.DeploymentId) > 0 {
		i -= len(m.DeploymentId)
		copy(dAtA[i:], m.DeploymentId)
		i = encodeVarintProcessListeningOnPortService(dAtA, i, uint64(len(m.DeploymentId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GetProcessesListeningOnPortsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetProcessesListeningOnPortsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetProcessesListeningOnPortsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ListeningEndpoints) > 0 {
		for iNdEx := len(m.ListeningEndpoints) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ListeningEndpoints[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProcessListeningOnPortService(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintProcessListeningOnPortService(dAtA []byte, offset int, v uint64) int {
	offset -= sovProcessListeningOnPortService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GetProcessesListeningOnPortsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DeploymentId)
	if l > 0 {
		n += 1 + l + sovProcessListeningOnPortService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *GetProcessesListeningOnPortsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ListeningEndpoints) > 0 {
		for _, e := range m.ListeningEndpoints {
			l = e.Size()
			n += 1 + l + sovProcessListeningOnPortService(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovProcessListeningOnPortService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProcessListeningOnPortService(x uint64) (n int) {
	return sovProcessListeningOnPortService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GetProcessesListeningOnPortsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessListeningOnPortService
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
			return fmt.Errorf("proto: GetProcessesListeningOnPortsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetProcessesListeningOnPortsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeploymentId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPortService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProcessListeningOnPortService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPortService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeploymentId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessListeningOnPortService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProcessListeningOnPortService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetProcessesListeningOnPortsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessListeningOnPortService
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
			return fmt.Errorf("proto: GetProcessesListeningOnPortsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetProcessesListeningOnPortsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListeningEndpoints", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPortService
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
				return ErrInvalidLengthProcessListeningOnPortService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPortService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ListeningEndpoints = append(m.ListeningEndpoints, &storage.ProcessListeningOnPort{})
			if err := m.ListeningEndpoints[len(m.ListeningEndpoints)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessListeningOnPortService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProcessListeningOnPortService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProcessListeningOnPortService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProcessListeningOnPortService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProcessListeningOnPortService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProcessListeningOnPortService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProcessListeningOnPortService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProcessListeningOnPortService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProcessListeningOnPortService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProcessListeningOnPortService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProcessListeningOnPortService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProcessListeningOnPortService = fmt.Errorf("proto: unexpected end of group")
)
