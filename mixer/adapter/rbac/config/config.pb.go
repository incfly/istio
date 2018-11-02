// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/rbac/config/config.proto

/*
	Package config is a generated protocol buffer package.

	The `rbac` adapter is deprecated by native RBAC implemented in Envoy proxy.
	See https://istio.io/docs/concepts/security/#enabling-authorization for enabling
	the native RBAC with your existing service role and service role binding.

	It is generated from these files:
		mixer/adapter/rbac/config/config.proto

	It has these top-level messages:
		Params
*/
package config

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/gogoproto"

import time "time"

import types "github.com/gogo/protobuf/types"

import strings "strings"
import reflect "reflect"

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

// Configuration format for the `rbac` adapter.
//
// For example, the following configuration defines a RBAC handler with
// configuration store URL pointing to Kubernetes etcd ("k8s://").
// If you want to run Mixer locally, you can set the configuration store
// URL to a local directory (e.g., "fs:///tmp/testdata/configroot").
//
// ```yaml
// apiVersion: "config.istio.io/v1alpha2"
// kind: rbac
// metadata:
//   name: rbachandler
//   namespace: istio-system
// spec:
//   config_store_url: "fs:///tmp/testdata/config"
// ```
type Params struct {
	// URL for the config store. It is used to initiate a new Store instance.
	// Following are some examples of the config store URL:
	// * "k8s://"
	// * "fs:///tmp/testdata/configroot"
	ConfigStoreUrl string `protobuf:"bytes,1,opt,name=config_store_url,json=configStoreUrl,proto3" json:"config_store_url,omitempty"`
	// The duration for which authorization results may be cached.
	CacheDuration time.Duration `protobuf:"bytes,2,opt,name=cache_duration,json=cacheDuration,stdduration" json:"cache_duration"`
}

func (m *Params) Reset()                    { *m = Params{} }
func (*Params) ProtoMessage()               {}
func (*Params) Descriptor() ([]byte, []int) { return fileDescriptorConfig, []int{0} }

func init() {
	proto.RegisterType((*Params)(nil), "adapter.rbac.config.Params")
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ConfigStoreUrl) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.ConfigStoreUrl)))
		i += copy(dAtA[i:], m.ConfigStoreUrl)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintConfig(dAtA, i, uint64(types.SizeOfStdDuration(m.CacheDuration)))
	n1, err := types.StdDurationMarshalTo(m.CacheDuration, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	return i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Params) Size() (n int) {
	var l int
	_ = l
	l = len(m.ConfigStoreUrl)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = types.SizeOfStdDuration(m.CacheDuration)
	n += 1 + l + sovConfig(uint64(l))
	return n
}

func sovConfig(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Params) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Params{`,
		`ConfigStoreUrl:` + fmt.Sprintf("%v", this.ConfigStoreUrl) + `,`,
		`CacheDuration:` + strings.Replace(strings.Replace(this.CacheDuration.String(), "Duration", "google_protobuf.Duration", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringConfig(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfigStoreUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConfigStoreUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CacheDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := types.StdDurationUnmarshal(&m.CacheDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthConfig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
				return 0, ErrInvalidLengthConfig
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowConfig
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
				next, err := skipConfig(dAtA[start:])
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
	ErrInvalidLengthConfig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("mixer/adapter/rbac/config/config.proto", fileDescriptorConfig) }

var fileDescriptorConfig = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xcb, 0xcd, 0xac, 0x48,
	0x2d, 0xd2, 0x4f, 0x4c, 0x49, 0x2c, 0x28, 0x49, 0x2d, 0xd2, 0x2f, 0x4a, 0x4a, 0x4c, 0xd6, 0x4f,
	0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0x87, 0x52, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xc2, 0x50,
	0x15, 0x7a, 0x20, 0x15, 0x7a, 0x10, 0x29, 0x29, 0xb9, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0x7d,
	0xb0, 0x92, 0xa4, 0xd2, 0x34, 0xfd, 0x94, 0xd2, 0xa2, 0xc4, 0x92, 0xcc, 0xfc, 0x3c, 0x88, 0x26,
	0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0x30, 0x53, 0x1f, 0xc4, 0x82, 0x88, 0x2a, 0xd5, 0x71, 0xb1,
	0x05, 0x24, 0x16, 0x25, 0xe6, 0x16, 0x0b, 0x69, 0x70, 0x09, 0x40, 0x4c, 0x8a, 0x2f, 0x2e, 0xc9,
	0x2f, 0x4a, 0x8d, 0x2f, 0x2d, 0xca, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0xe2, 0x83, 0x88,
	0x07, 0x83, 0x84, 0x43, 0x8b, 0x72, 0x84, 0xbc, 0xb8, 0xf8, 0x92, 0x13, 0x93, 0x33, 0x52, 0xe3,
	0x61, 0x36, 0x48, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x1b, 0x49, 0xea, 0x41, 0x9c, 0xa0, 0x07, 0x73,
	0x82, 0x9e, 0x0b, 0x54, 0x81, 0x13, 0xc7, 0x89, 0x7b, 0xf2, 0x0c, 0x33, 0xee, 0xcb, 0x33, 0x06,
	0xf1, 0x82, 0xb5, 0xc2, 0x25, 0x2c, 0x4e, 0x3c, 0x94, 0x63, 0xb8, 0xf0, 0x50, 0x8e, 0xe1, 0xc6,
	0x43, 0x39, 0x86, 0x0f, 0x0f, 0xe5, 0x18, 0x1a, 0x1e, 0xc9, 0x31, 0xae, 0x78, 0x24, 0xc7, 0x70,
	0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xbe, 0x78, 0x24, 0xc7,
	0xf0, 0xe1, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x51, 0x6c, 0x10, 0xd7, 0x24, 0xb1, 0x81,
	0x6d, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x7c, 0xf7, 0x13, 0xb7, 0x35, 0x01, 0x00, 0x00,
}
