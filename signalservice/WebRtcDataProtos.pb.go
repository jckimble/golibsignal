// Code generated by protoc-gen-go. DO NOT EDIT.
// source: WebRtcDataProtos.proto

package signalservice

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type Connected struct {
	Id                   *uint64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Connected) Reset()         { *m = Connected{} }
func (m *Connected) String() string { return proto.CompactTextString(m) }
func (*Connected) ProtoMessage()    {}
func (*Connected) Descriptor() ([]byte, []int) {
	return fileDescriptor_1173f613e5ad5dbe, []int{0}
}

func (m *Connected) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Connected.Unmarshal(m, b)
}
func (m *Connected) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Connected.Marshal(b, m, deterministic)
}
func (m *Connected) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Connected.Merge(m, src)
}
func (m *Connected) XXX_Size() int {
	return xxx_messageInfo_Connected.Size(m)
}
func (m *Connected) XXX_DiscardUnknown() {
	xxx_messageInfo_Connected.DiscardUnknown(m)
}

var xxx_messageInfo_Connected proto.InternalMessageInfo

func (m *Connected) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

type Hangup struct {
	Id                   *uint64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hangup) Reset()         { *m = Hangup{} }
func (m *Hangup) String() string { return proto.CompactTextString(m) }
func (*Hangup) ProtoMessage()    {}
func (*Hangup) Descriptor() ([]byte, []int) {
	return fileDescriptor_1173f613e5ad5dbe, []int{1}
}

func (m *Hangup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hangup.Unmarshal(m, b)
}
func (m *Hangup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hangup.Marshal(b, m, deterministic)
}
func (m *Hangup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hangup.Merge(m, src)
}
func (m *Hangup) XXX_Size() int {
	return xxx_messageInfo_Hangup.Size(m)
}
func (m *Hangup) XXX_DiscardUnknown() {
	xxx_messageInfo_Hangup.DiscardUnknown(m)
}

var xxx_messageInfo_Hangup proto.InternalMessageInfo

func (m *Hangup) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

type VideoStreamingStatus struct {
	Id                   *uint64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Enabled              *bool    `protobuf:"varint,2,opt,name=enabled" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VideoStreamingStatus) Reset()         { *m = VideoStreamingStatus{} }
func (m *VideoStreamingStatus) String() string { return proto.CompactTextString(m) }
func (*VideoStreamingStatus) ProtoMessage()    {}
func (*VideoStreamingStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_1173f613e5ad5dbe, []int{2}
}

func (m *VideoStreamingStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VideoStreamingStatus.Unmarshal(m, b)
}
func (m *VideoStreamingStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VideoStreamingStatus.Marshal(b, m, deterministic)
}
func (m *VideoStreamingStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VideoStreamingStatus.Merge(m, src)
}
func (m *VideoStreamingStatus) XXX_Size() int {
	return xxx_messageInfo_VideoStreamingStatus.Size(m)
}
func (m *VideoStreamingStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_VideoStreamingStatus.DiscardUnknown(m)
}

var xxx_messageInfo_VideoStreamingStatus proto.InternalMessageInfo

func (m *VideoStreamingStatus) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *VideoStreamingStatus) GetEnabled() bool {
	if m != nil && m.Enabled != nil {
		return *m.Enabled
	}
	return false
}

type Data struct {
	Connected            *Connected            `protobuf:"bytes,1,opt,name=connected" json:"connected,omitempty"`
	Hangup               *Hangup               `protobuf:"bytes,2,opt,name=hangup" json:"hangup,omitempty"`
	VideoStreamingStatus *VideoStreamingStatus `protobuf:"bytes,3,opt,name=videoStreamingStatus" json:"videoStreamingStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_1173f613e5ad5dbe, []int{3}
}

func (m *Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Data.Unmarshal(m, b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Data.Marshal(b, m, deterministic)
}
func (m *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(m, src)
}
func (m *Data) XXX_Size() int {
	return xxx_messageInfo_Data.Size(m)
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetConnected() *Connected {
	if m != nil {
		return m.Connected
	}
	return nil
}

func (m *Data) GetHangup() *Hangup {
	if m != nil {
		return m.Hangup
	}
	return nil
}

func (m *Data) GetVideoStreamingStatus() *VideoStreamingStatus {
	if m != nil {
		return m.VideoStreamingStatus
	}
	return nil
}

func init() {
	proto.RegisterType((*Connected)(nil), "signalservice.Connected")
	proto.RegisterType((*Hangup)(nil), "signalservice.Hangup")
	proto.RegisterType((*VideoStreamingStatus)(nil), "signalservice.VideoStreamingStatus")
	proto.RegisterType((*Data)(nil), "signalservice.Data")
}

func init() { proto.RegisterFile("WebRtcDataProtos.proto", fileDescriptor_1173f613e5ad5dbe) }

var fileDescriptor_1173f613e5ad5dbe = []byte{
	// 249 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xd9, 0xb5, 0x54, 0x3b, 0x45, 0x91, 0x50, 0x25, 0xe0, 0xa5, 0xae, 0x97, 0x5e, 0xcc,
	0xa1, 0xa0, 0x67, 0xa9, 0x1e, 0x3c, 0x4a, 0x0a, 0xf6, 0x9c, 0x4d, 0x86, 0x6c, 0xa0, 0x9b, 0x94,
	0x64, 0xb6, 0xbe, 0xa1, 0xcf, 0x25, 0x5d, 0x5d, 0xc5, 0xed, 0x9e, 0x42, 0x98, 0xff, 0xfb, 0xf9,
	0x86, 0x81, 0xeb, 0x0d, 0x96, 0x92, 0xf4, 0x8b, 0x22, 0xf5, 0x16, 0x03, 0x85, 0x24, 0x76, 0x87,
	0x87, 0x9d, 0x27, 0x67, 0xbd, 0xda, 0x26, 0x8c, 0x7b, 0xa7, 0xb1, 0xb8, 0x81, 0xc9, 0x73, 0xf0,
	0x1e, 0x35, 0xa1, 0x61, 0x17, 0x90, 0x3b, 0xc3, 0xb3, 0x79, 0xb6, 0x18, 0xc9, 0xdc, 0x99, 0x82,
	0xc3, 0xf8, 0x55, 0x79, 0xdb, 0xec, 0x8e, 0x26, 0x4f, 0x30, 0x7b, 0x77, 0x06, 0xc3, 0x9a, 0x22,
	0xaa, 0xda, 0x79, 0xbb, 0x26, 0x45, 0x4d, 0xea, 0xe7, 0x18, 0x87, 0x53, 0xf4, 0xaa, 0xdc, 0xa2,
	0xe1, 0xf9, 0x3c, 0x5b, 0x9c, 0xc9, 0xee, 0x5b, 0x7c, 0x66, 0x30, 0x3a, 0xc8, 0xb1, 0x47, 0x98,
	0xe8, 0xce, 0xa0, 0x25, 0xa7, 0x4b, 0x2e, 0xfe, 0x49, 0x8a, 0x5f, 0x43, 0xf9, 0x17, 0x65, 0xf7,
	0x30, 0xae, 0x5a, 0xb9, 0xb6, 0x79, 0xba, 0xbc, 0xea, 0x41, 0xdf, 0xe6, 0xf2, 0x27, 0xc4, 0x36,
	0x30, 0xdb, 0x0f, 0x18, 0xf3, 0x93, 0x16, 0xbe, 0xeb, 0xc1, 0x43, 0xcb, 0xc9, 0xc1, 0x82, 0xd5,
	0x03, 0xdc, 0x86, 0x68, 0x05, 0x55, 0xa1, 0xb1, 0x15, 0xe9, 0xe8, 0x6a, 0x14, 0x09, 0x75, 0x13,
	0x31, 0xd5, 0x49, 0x7c, 0x60, 0x19, 0x49, 0xaf, 0x2e, 0xfb, 0xd7, 0xf8, 0x0a, 0x00, 0x00, 0xff,
	0xff, 0xd2, 0x87, 0x14, 0x4b, 0xa0, 0x01, 0x00, 0x00,
}