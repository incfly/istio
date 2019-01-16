// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/client_ssl_auth/v2/client_ssl_auth.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/lyft/protoc-gen-validate/validate"

import time "time"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ClientSSLAuth struct {
	// The :ref:`cluster manager <arch_overview_cluster_manager>` cluster that runs
	// the authentication service. The filter will connect to the service every 60s to fetch the list
	// of principals. The service must support the expected :ref:`REST API
	// <config_network_filters_client_ssl_auth_rest_api>`.
	AuthApiCluster string `protobuf:"bytes,1,opt,name=auth_api_cluster,json=authApiCluster,proto3" json:"auth_api_cluster,omitempty"`
	// The prefix to use when emitting :ref:`statistics
	// <config_network_filters_client_ssl_auth_stats>`.
	StatPrefix string `protobuf:"bytes,2,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	// Time in milliseconds between principal refreshes from the
	// authentication service. Default is 60000 (60s). The actual fetch time
	// will be this value plus a random jittered value between
	// 0-refresh_delay_ms milliseconds.
	RefreshDelay *time.Duration `protobuf:"bytes,3,opt,name=refresh_delay,json=refreshDelay,proto3,stdduration" json:"refresh_delay,omitempty"`
	// An optional list of IP address and subnet masks that should be white
	// listed for access by the filter. If no list is provided, there is no
	// IP white list.
	IpWhiteList          []*core.CidrRange `protobuf:"bytes,4,rep,name=ip_white_list,json=ipWhiteList,proto3" json:"ip_white_list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ClientSSLAuth) Reset()         { *m = ClientSSLAuth{} }
func (m *ClientSSLAuth) String() string { return proto.CompactTextString(m) }
func (*ClientSSLAuth) ProtoMessage()    {}
func (*ClientSSLAuth) Descriptor() ([]byte, []int) {
	return fileDescriptor_client_ssl_auth_4986a886d9153439, []int{0}
}
func (m *ClientSSLAuth) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClientSSLAuth) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClientSSLAuth.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ClientSSLAuth) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientSSLAuth.Merge(dst, src)
}
func (m *ClientSSLAuth) XXX_Size() int {
	return m.Size()
}
func (m *ClientSSLAuth) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientSSLAuth.DiscardUnknown(m)
}

var xxx_messageInfo_ClientSSLAuth proto.InternalMessageInfo

func (m *ClientSSLAuth) GetAuthApiCluster() string {
	if m != nil {
		return m.AuthApiCluster
	}
	return ""
}

func (m *ClientSSLAuth) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *ClientSSLAuth) GetRefreshDelay() *time.Duration {
	if m != nil {
		return m.RefreshDelay
	}
	return nil
}

func (m *ClientSSLAuth) GetIpWhiteList() []*core.CidrRange {
	if m != nil {
		return m.IpWhiteList
	}
	return nil
}

func init() {
	proto.RegisterType((*ClientSSLAuth)(nil), "envoy.config.filter.network.client_ssl_auth.v2.ClientSSLAuth")
}
func (m *ClientSSLAuth) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientSSLAuth) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.AuthApiCluster) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintClientSslAuth(dAtA, i, uint64(len(m.AuthApiCluster)))
		i += copy(dAtA[i:], m.AuthApiCluster)
	}
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintClientSslAuth(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if m.RefreshDelay != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintClientSslAuth(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdDuration(*m.RefreshDelay)))
		n1, err := github_com_gogo_protobuf_types.StdDurationMarshalTo(*m.RefreshDelay, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.IpWhiteList) > 0 {
		for _, msg := range m.IpWhiteList {
			dAtA[i] = 0x22
			i++
			i = encodeVarintClientSslAuth(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintClientSslAuth(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ClientSSLAuth) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AuthApiCluster)
	if l > 0 {
		n += 1 + l + sovClientSslAuth(uint64(l))
	}
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovClientSslAuth(uint64(l))
	}
	if m.RefreshDelay != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdDuration(*m.RefreshDelay)
		n += 1 + l + sovClientSslAuth(uint64(l))
	}
	if len(m.IpWhiteList) > 0 {
		for _, e := range m.IpWhiteList {
			l = e.Size()
			n += 1 + l + sovClientSslAuth(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovClientSslAuth(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozClientSslAuth(x uint64) (n int) {
	return sovClientSslAuth(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClientSSLAuth) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClientSslAuth
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ClientSSLAuth: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientSSLAuth: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthApiCluster", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientSslAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthClientSslAuth
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AuthApiCluster = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientSslAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthClientSslAuth
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RefreshDelay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientSslAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthClientSslAuth
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RefreshDelay == nil {
				m.RefreshDelay = new(time.Duration)
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(m.RefreshDelay, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IpWhiteList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientSslAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthClientSslAuth
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IpWhiteList = append(m.IpWhiteList, &core.CidrRange{})
			if err := m.IpWhiteList[len(m.IpWhiteList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClientSslAuth(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClientSslAuth
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
func skipClientSslAuth(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowClientSslAuth
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
					return 0, ErrIntOverflowClientSslAuth
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowClientSslAuth
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthClientSslAuth
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowClientSslAuth
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipClientSslAuth(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthClientSslAuth = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowClientSslAuth   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/client_ssl_auth/v2/client_ssl_auth.proto", fileDescriptor_client_ssl_auth_4986a886d9153439)
}

var fileDescriptor_client_ssl_auth_4986a886d9153439 = []byte{
	// 366 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcf, 0x6a, 0xe3, 0x30,
	0x10, 0xc6, 0x51, 0x12, 0x16, 0x62, 0x6f, 0x96, 0x60, 0x16, 0xd6, 0x1b, 0x16, 0xc7, 0xec, 0x29,
	0xec, 0x41, 0x02, 0xe7, 0x05, 0x36, 0x89, 0x8f, 0x39, 0x14, 0xe7, 0x50, 0xe8, 0x45, 0x28, 0xb1,
	0x6c, 0x0f, 0x15, 0x96, 0x90, 0x64, 0xa7, 0x79, 0x93, 0x3e, 0x4b, 0x4f, 0x3d, 0xf6, 0xd8, 0x37,
	0x68, 0xc9, 0xad, 0x6f, 0xd0, 0x63, 0xf1, 0x9f, 0x5c, 0x72, 0x1b, 0xe9, 0x9b, 0xef, 0x37, 0xf3,
	0x8d, 0x13, 0xf3, 0xb2, 0x96, 0x27, 0x72, 0x90, 0x65, 0x06, 0x39, 0xc9, 0x40, 0x58, 0xae, 0x49,
	0xc9, 0xed, 0x51, 0xea, 0x7b, 0x72, 0x10, 0xc0, 0x4b, 0x4b, 0x8d, 0x11, 0x94, 0x55, 0xb6, 0x20,
	0x75, 0x74, 0xfd, 0x85, 0x95, 0x96, 0x56, 0x7a, 0xb8, 0xa5, 0xe0, 0x8e, 0x82, 0x3b, 0x0a, 0xee,
	0x29, 0xf8, 0xda, 0x52, 0x47, 0xb3, 0x79, 0x37, 0x95, 0x29, 0x68, 0x99, 0x52, 0x73, 0xc2, 0xd2,
	0x54, 0x73, 0x63, 0x3a, 0xe0, 0x2c, 0xc8, 0xa5, 0xcc, 0x05, 0x27, 0xed, 0x6b, 0x5f, 0x65, 0x24,
	0xad, 0x34, 0xb3, 0x20, 0xcb, 0x5e, 0xff, 0x55, 0x33, 0x01, 0x29, 0xb3, 0x9c, 0x5c, 0x8a, 0x5e,
	0xf8, 0x99, 0xcb, 0x5c, 0xb6, 0x25, 0x69, 0xaa, 0xee, 0xf7, 0xef, 0x27, 0x72, 0x26, 0x9b, 0x76,
	0x8d, 0xdd, 0x6e, 0xbb, 0xaa, 0x6c, 0xe1, 0x2d, 0x9d, 0x69, 0xb3, 0x0c, 0x65, 0x0a, 0xe8, 0x41,
	0x54, 0xc6, 0x72, 0xed, 0xa3, 0x10, 0x2d, 0xc6, 0xeb, 0xf1, 0xd3, 0xc7, 0xf3, 0x70, 0xa4, 0x07,
	0x21, 0x4a, 0x7e, 0x34, 0x2d, 0x2b, 0x05, 0x9b, 0xae, 0xc1, 0xfb, 0xe7, 0xb8, 0xc6, 0x32, 0x4b,
	0x95, 0xe6, 0x19, 0x3c, 0xf8, 0x83, 0xeb, 0x7e, 0xa7, 0x51, 0x6f, 0x5a, 0xd1, 0x8b, 0x9d, 0x89,
	0xe6, 0x99, 0xe6, 0xa6, 0xa0, 0x29, 0x17, 0xec, 0xe4, 0x0f, 0x43, 0xb4, 0x70, 0xa3, 0xdf, 0xb8,
	0x4b, 0x86, 0x2f, 0xc9, 0x70, 0xdc, 0x27, 0x5b, 0x8f, 0x1e, 0xdf, 0xe6, 0x28, 0xf9, 0xde, 0xbb,
	0xe2, 0xc6, 0xe4, 0xfd, 0x77, 0x26, 0xa0, 0xe8, 0xb1, 0x00, 0xcb, 0xa9, 0x00, 0x63, 0xfd, 0x51,
	0x38, 0x5c, 0xb8, 0xd1, 0x9f, 0xfe, 0xe0, 0x4c, 0x01, 0xae, 0x23, 0xdc, 0x1c, 0x10, 0x6f, 0x20,
	0xd5, 0x09, 0x2b, 0x73, 0x9e, 0xb8, 0xa0, 0x6e, 0x1b, 0xc7, 0x16, 0x8c, 0x5d, 0x4f, 0x5f, 0xce,
	0x01, 0x7a, 0x3d, 0x07, 0xe8, 0xfd, 0x1c, 0xa0, 0xbb, 0x41, 0x1d, 0xed, 0xbf, 0xb5, 0xa3, 0x97,
	0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x49, 0xe3, 0x8a, 0x3f, 0xfb, 0x01, 0x00, 0x00,
}
