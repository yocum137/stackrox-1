// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/process_listening_on_port.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
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

type ProcessListeningOnPort struct {
	Port                 uint32                     `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Protocol             L4Protocol                 `protobuf:"varint,2,opt,name=protocol,proto3,enum=storage.L4Protocol" json:"protocol,omitempty"`
	Process              *ProcessIndicatorUniqueKey `protobuf:"bytes,3,opt,name=process,proto3" json:"process,omitempty"`
	CloseTimestamp       *types.Timestamp           `protobuf:"bytes,4,opt,name=close_timestamp,json=closeTimestamp,proto3" json:"close_timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ProcessListeningOnPort) Reset()         { *m = ProcessListeningOnPort{} }
func (m *ProcessListeningOnPort) String() string { return proto.CompactTextString(m) }
func (*ProcessListeningOnPort) ProtoMessage()    {}
func (*ProcessListeningOnPort) Descriptor() ([]byte, []int) {
	return fileDescriptor_44bd1925a567394f, []int{0}
}
func (m *ProcessListeningOnPort) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessListeningOnPort) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessListeningOnPort.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessListeningOnPort) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessListeningOnPort.Merge(m, src)
}
func (m *ProcessListeningOnPort) XXX_Size() int {
	return m.Size()
}
func (m *ProcessListeningOnPort) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessListeningOnPort.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessListeningOnPort proto.InternalMessageInfo

func (m *ProcessListeningOnPort) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *ProcessListeningOnPort) GetProtocol() L4Protocol {
	if m != nil {
		return m.Protocol
	}
	return L4Protocol_L4_PROTOCOL_UNKNOWN
}

func (m *ProcessListeningOnPort) GetProcess() *ProcessIndicatorUniqueKey {
	if m != nil {
		return m.Process
	}
	return nil
}

func (m *ProcessListeningOnPort) GetCloseTimestamp() *types.Timestamp {
	if m != nil {
		return m.CloseTimestamp
	}
	return nil
}

func (m *ProcessListeningOnPort) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ProcessListeningOnPort) Clone() *ProcessListeningOnPort {
	if m == nil {
		return nil
	}
	cloned := new(ProcessListeningOnPort)
	*cloned = *m

	cloned.Process = m.Process.Clone()
	cloned.CloseTimestamp = m.CloseTimestamp.Clone()
	return cloned
}

type ProcessListeningOnPortStorage struct {
	// Ideally it has to be GENERATED ALWAYS AS IDENTITY, which will make it a
	// bigint with a sequence. Unfortunately at the moment some bits of store
	// generator assume an id has to be a string.
	Id                   string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk,type(uuid)"`
	Port                 uint32           `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Protocol             L4Protocol       `protobuf:"varint,3,opt,name=protocol,proto3,enum=storage.L4Protocol" json:"protocol,omitempty"`
	CloseTimestamp       *types.Timestamp `protobuf:"bytes,4,opt,name=close_timestamp,json=closeTimestamp,proto3" json:"close_timestamp,omitempty"`
	ProcessIndicatorId   string           `protobuf:"bytes,5,opt,name=process_indicator_id,json=processIndicatorId,proto3" json:"process_indicator_id,omitempty" search:"Process ID,store" sql:"fk(ProcessIndicator:id),no-fk-constraint,index=btree,type(uuid)"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ProcessListeningOnPortStorage) Reset()         { *m = ProcessListeningOnPortStorage{} }
func (m *ProcessListeningOnPortStorage) String() string { return proto.CompactTextString(m) }
func (*ProcessListeningOnPortStorage) ProtoMessage()    {}
func (*ProcessListeningOnPortStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_44bd1925a567394f, []int{1}
}
func (m *ProcessListeningOnPortStorage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessListeningOnPortStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessListeningOnPortStorage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessListeningOnPortStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessListeningOnPortStorage.Merge(m, src)
}
func (m *ProcessListeningOnPortStorage) XXX_Size() int {
	return m.Size()
}
func (m *ProcessListeningOnPortStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessListeningOnPortStorage.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessListeningOnPortStorage proto.InternalMessageInfo

func (m *ProcessListeningOnPortStorage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProcessListeningOnPortStorage) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *ProcessListeningOnPortStorage) GetProtocol() L4Protocol {
	if m != nil {
		return m.Protocol
	}
	return L4Protocol_L4_PROTOCOL_UNKNOWN
}

func (m *ProcessListeningOnPortStorage) GetCloseTimestamp() *types.Timestamp {
	if m != nil {
		return m.CloseTimestamp
	}
	return nil
}

func (m *ProcessListeningOnPortStorage) GetProcessIndicatorId() string {
	if m != nil {
		return m.ProcessIndicatorId
	}
	return ""
}

func (m *ProcessListeningOnPortStorage) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ProcessListeningOnPortStorage) Clone() *ProcessListeningOnPortStorage {
	if m == nil {
		return nil
	}
	cloned := new(ProcessListeningOnPortStorage)
	*cloned = *m

	cloned.CloseTimestamp = m.CloseTimestamp.Clone()
	return cloned
}

func init() {
	proto.RegisterType((*ProcessListeningOnPort)(nil), "storage.ProcessListeningOnPort")
	proto.RegisterType((*ProcessListeningOnPortStorage)(nil), "storage.ProcessListeningOnPortStorage")
}

func init() {
	proto.RegisterFile("storage/process_listening_on_port.proto", fileDescriptor_44bd1925a567394f)
}

var fileDescriptor_44bd1925a567394f = []byte{
	// 434 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0x4f, 0x6b, 0xd4, 0x40,
	0x14, 0x77, 0xd2, 0x6a, 0x75, 0xc4, 0x0a, 0xd3, 0xa2, 0x61, 0xd1, 0xdd, 0x90, 0x4b, 0x53, 0x48,
	0x13, 0xa8, 0x9e, 0x16, 0xbd, 0x54, 0x2f, 0x8b, 0x05, 0x97, 0xa8, 0x17, 0x2f, 0x21, 0x9b, 0x99,
	0xc4, 0x21, 0xe9, 0xbc, 0x74, 0x66, 0x42, 0xed, 0x07, 0x11, 0xfc, 0x48, 0x1e, 0xfd, 0x04, 0xa5,
	0xac, 0x37, 0x8f, 0x7b, 0xf2, 0x28, 0x3b, 0x99, 0x04, 0x5c, 0x15, 0x3c, 0x78, 0x7b, 0xc9, 0xfb,
	0xfd, 0xf2, 0xfb, 0xf3, 0x82, 0x0f, 0x94, 0x06, 0x99, 0x95, 0x2c, 0x6e, 0x24, 0xe4, 0x4c, 0xa9,
	0xb4, 0xe6, 0x4a, 0x33, 0xc1, 0x45, 0x99, 0x82, 0x48, 0x1b, 0x90, 0x3a, 0x6a, 0x24, 0x68, 0x20,
	0x3b, 0x16, 0x38, 0x9a, 0x94, 0x00, 0x65, 0x6d, 0x08, 0x1a, 0x16, 0x6d, 0x11, 0x6b, 0x7e, 0xc6,
	0x94, 0xce, 0xce, 0x9a, 0x0e, 0x39, 0xda, 0x2f, 0xa1, 0x04, 0x33, 0xc6, 0xeb, 0xc9, 0xbe, 0x1d,
	0xf5, 0x42, 0x82, 0xe9, 0x0b, 0x90, 0x55, 0x5a, 0xd4, 0x70, 0x61, 0x77, 0x93, 0x4d, 0x13, 0x5c,
	0x50, 0x9e, 0x67, 0x1a, 0x64, 0x07, 0xf0, 0xaf, 0x11, 0x7e, 0x30, 0xef, 0x76, 0xa7, 0xbd, 0xbf,
	0xd7, 0x62, 0x0e, 0x52, 0x13, 0x82, 0xb7, 0xd7, 0x2e, 0x5d, 0xe4, 0xa1, 0xe0, 0x5e, 0x62, 0x66,
	0x12, 0xe3, 0xdb, 0x86, 0x97, 0x43, 0xed, 0x3a, 0x1e, 0x0a, 0x76, 0x8f, 0xf7, 0x22, 0x2b, 0x11,
	0x9d, 0x3e, 0x9d, 0xdb, 0x55, 0x32, 0x80, 0xc8, 0x33, 0xbc, 0x63, 0xa5, 0xdd, 0x2d, 0x0f, 0x05,
	0x77, 0x8f, 0xfd, 0x01, 0x6f, 0x65, 0x67, 0xbd, 0xa3, 0x77, 0x82, 0x9f, 0xb7, 0xec, 0x15, 0xbb,
	0x4c, 0x7a, 0x0a, 0x79, 0x81, 0xef, 0xe7, 0x35, 0x28, 0x96, 0x0e, 0x4d, 0xb8, 0xdb, 0xe6, 0x2b,
	0xa3, 0xa8, 0xeb, 0x2a, 0xea, 0xbb, 0x8a, 0xde, 0xf6, 0x88, 0x64, 0xd7, 0x50, 0x86, 0x67, 0xff,
	0xbb, 0x83, 0x1f, 0xff, 0x39, 0xe2, 0x9b, 0xce, 0x09, 0x39, 0xc0, 0x0e, 0xa7, 0x26, 0xe7, 0x9d,
	0x93, 0x87, 0xab, 0xab, 0xc9, 0x9e, 0x3a, 0xaf, 0xa7, 0x7e, 0x53, 0x85, 0xfa, 0xb2, 0x61, 0x41,
	0xdb, 0x72, 0x7a, 0xe8, 0x27, 0x0e, 0xa7, 0x43, 0x25, 0xce, 0x5f, 0x2a, 0xd9, 0xfa, 0x97, 0x4a,
	0xfe, 0x47, 0x28, 0xf2, 0x09, 0xe1, 0xfd, 0xdf, 0x6e, 0x9a, 0x72, 0xea, 0xde, 0x34, 0x29, 0xf2,
	0xd5, 0xd5, 0x24, 0x55, 0x2c, 0x93, 0xf9, 0x87, 0xa9, 0x6f, 0xc3, 0x7b, 0xb3, 0x97, 0xe1, 0xda,
	0x18, 0xf3, 0x3d, 0x93, 0xaf, 0xa8, 0x82, 0xcd, 0x13, 0x4c, 0x39, 0x3d, 0x0c, 0x05, 0x1c, 0x15,
	0xd5, 0x51, 0x0e, 0x42, 0x69, 0x99, 0x71, 0xa1, 0x43, 0x2e, 0x28, 0xfb, 0xf8, 0x7c, 0xa1, 0x25,
	0x63, 0xbf, 0x34, 0x42, 0x9a, 0x0d, 0xfa, 0x8c, 0x9e, 0x3c, 0xfa, 0xb2, 0x1c, 0xa3, 0xaf, 0xcb,
	0x31, 0xba, 0x5e, 0x8e, 0xd1, 0xe7, 0x6f, 0xe3, 0x1b, 0xef, 0xfb, 0xdf, 0xfb, 0x07, 0x42, 0x8b,
	0x5b, 0x26, 0xd9, 0x93, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x68, 0xa2, 0x43, 0x9c, 0x1c, 0x03,
	0x00, 0x00,
}

func (m *ProcessListeningOnPort) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessListeningOnPort) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessListeningOnPort) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.CloseTimestamp != nil {
		{
			size, err := m.CloseTimestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Process != nil {
		{
			size, err := m.Process.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Protocol != 0 {
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(m.Protocol))
		i--
		dAtA[i] = 0x10
	}
	if m.Port != 0 {
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(m.Port))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ProcessListeningOnPortStorage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessListeningOnPortStorage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessListeningOnPortStorage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ProcessIndicatorId) > 0 {
		i -= len(m.ProcessIndicatorId)
		copy(dAtA[i:], m.ProcessIndicatorId)
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(len(m.ProcessIndicatorId)))
		i--
		dAtA[i] = 0x2a
	}
	if m.CloseTimestamp != nil {
		{
			size, err := m.CloseTimestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Protocol != 0 {
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(m.Protocol))
		i--
		dAtA[i] = 0x18
	}
	if m.Port != 0 {
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(m.Port))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintProcessListeningOnPort(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProcessListeningOnPort(dAtA []byte, offset int, v uint64) int {
	offset -= sovProcessListeningOnPort(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProcessListeningOnPort) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Port != 0 {
		n += 1 + sovProcessListeningOnPort(uint64(m.Port))
	}
	if m.Protocol != 0 {
		n += 1 + sovProcessListeningOnPort(uint64(m.Protocol))
	}
	if m.Process != nil {
		l = m.Process.Size()
		n += 1 + l + sovProcessListeningOnPort(uint64(l))
	}
	if m.CloseTimestamp != nil {
		l = m.CloseTimestamp.Size()
		n += 1 + l + sovProcessListeningOnPort(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ProcessListeningOnPortStorage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovProcessListeningOnPort(uint64(l))
	}
	if m.Port != 0 {
		n += 1 + sovProcessListeningOnPort(uint64(m.Port))
	}
	if m.Protocol != 0 {
		n += 1 + sovProcessListeningOnPort(uint64(m.Protocol))
	}
	if m.CloseTimestamp != nil {
		l = m.CloseTimestamp.Size()
		n += 1 + l + sovProcessListeningOnPort(uint64(l))
	}
	l = len(m.ProcessIndicatorId)
	if l > 0 {
		n += 1 + l + sovProcessListeningOnPort(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovProcessListeningOnPort(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProcessListeningOnPort(x uint64) (n int) {
	return sovProcessListeningOnPort(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProcessListeningOnPort) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessListeningOnPort
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
			return fmt.Errorf("proto: ProcessListeningOnPort: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessListeningOnPort: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Port", wireType)
			}
			m.Port = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Port |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Protocol", wireType)
			}
			m.Protocol = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Protocol |= L4Protocol(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Process", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
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
				return ErrInvalidLengthProcessListeningOnPort
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPort
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Process == nil {
				m.Process = &ProcessIndicatorUniqueKey{}
			}
			if err := m.Process.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CloseTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
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
				return ErrInvalidLengthProcessListeningOnPort
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPort
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CloseTimestamp == nil {
				m.CloseTimestamp = &types.Timestamp{}
			}
			if err := m.CloseTimestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessListeningOnPort(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProcessListeningOnPort
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
func (m *ProcessListeningOnPortStorage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessListeningOnPort
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
			return fmt.Errorf("proto: ProcessListeningOnPortStorage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessListeningOnPortStorage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
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
				return ErrInvalidLengthProcessListeningOnPort
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPort
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Port", wireType)
			}
			m.Port = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Port |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Protocol", wireType)
			}
			m.Protocol = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Protocol |= L4Protocol(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CloseTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
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
				return ErrInvalidLengthProcessListeningOnPort
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPort
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CloseTimestamp == nil {
				m.CloseTimestamp = &types.Timestamp{}
			}
			if err := m.CloseTimestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProcessIndicatorId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPort
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
				return ErrInvalidLengthProcessListeningOnPort
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPort
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProcessIndicatorId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessListeningOnPort(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProcessListeningOnPort
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
func skipProcessListeningOnPort(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProcessListeningOnPort
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
					return 0, ErrIntOverflowProcessListeningOnPort
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
					return 0, ErrIntOverflowProcessListeningOnPort
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
				return 0, ErrInvalidLengthProcessListeningOnPort
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProcessListeningOnPort
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProcessListeningOnPort
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProcessListeningOnPort        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProcessListeningOnPort          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProcessListeningOnPort = fmt.Errorf("proto: unexpected end of group")
)
