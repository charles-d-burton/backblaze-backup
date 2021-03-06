// Code generated by protoc-gen-go.
// source: files.proto
// DO NOT EDIT!

/*
Package filesystem is a generated protocol buffer package.

It is generated from these files:
	files.proto

It has these top-level messages:
	MetaData
*/
package filesystem

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MetaData struct {
	Name     string   `protobuf:"bytes,1,opt,name=Name,json=name" json:"Name,omitempty"`
	Size     int64    `protobuf:"varint,2,opt,name=Size,json=size" json:"Size,omitempty"`
	Sha1     string   `protobuf:"bytes,3,opt,name=Sha1,json=sha1" json:"Sha1,omitempty"`
	BackedUp bool     `protobuf:"varint,4,opt,name=BackedUp,json=backedUp" json:"BackedUp,omitempty"`
	Versions []string `protobuf:"bytes,5,rep,name=Versions,json=versions" json:"Versions,omitempty"`
}

func (m *MetaData) Reset()                    { *m = MetaData{} }
func (m *MetaData) String() string            { return proto.CompactTextString(m) }
func (*MetaData) ProtoMessage()               {}
func (*MetaData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MetaData) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetaData) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *MetaData) GetSha1() string {
	if m != nil {
		return m.Sha1
	}
	return ""
}

func (m *MetaData) GetBackedUp() bool {
	if m != nil {
		return m.BackedUp
	}
	return false
}

func (m *MetaData) GetVersions() []string {
	if m != nil {
		return m.Versions
	}
	return nil
}

func init() {
	proto.RegisterType((*MetaData)(nil), "filesystem.MetaData")
}

func init() { proto.RegisterFile("files.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0xcb, 0xcc, 0x49,
	0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x02, 0x73, 0x2a, 0x8b, 0x4b, 0x52, 0x73,
	0x95, 0xea, 0xb8, 0x38, 0x7c, 0x53, 0x4b, 0x12, 0x5d, 0x12, 0x4b, 0x12, 0x85, 0x84, 0xb8, 0x58,
	0xfc, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x58, 0xf2, 0x12, 0x73, 0x53,
	0x41, 0x62, 0xc1, 0x99, 0x55, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x2c, 0xc5, 0x99,
	0x55, 0x10, 0xb1, 0x8c, 0x44, 0x43, 0x09, 0x66, 0x88, 0xba, 0xe2, 0x8c, 0x44, 0x43, 0x21, 0x29,
	0x2e, 0x0e, 0xa7, 0xc4, 0xe4, 0xec, 0xd4, 0x94, 0xd0, 0x02, 0x09, 0x16, 0x05, 0x46, 0x0d, 0x8e,
	0x20, 0x8e, 0x24, 0x28, 0x1f, 0x24, 0x17, 0x96, 0x5a, 0x54, 0x9c, 0x99, 0x9f, 0x57, 0x2c, 0xc1,
	0xaa, 0xc0, 0xac, 0xc1, 0x19, 0xc4, 0x51, 0x06, 0xe5, 0x27, 0xb1, 0x81, 0x9d, 0x64, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0xfa, 0xc0, 0xb0, 0xf3, 0xa1, 0x00, 0x00, 0x00,
}
