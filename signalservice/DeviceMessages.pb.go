// Code generated by protoc-gen-go. DO NOT EDIT.
// source: DeviceMessages.proto

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

type ProvisioningUuid struct {
	Uuid                 *string  `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProvisioningUuid) Reset()         { *m = ProvisioningUuid{} }
func (m *ProvisioningUuid) String() string { return proto.CompactTextString(m) }
func (*ProvisioningUuid) ProtoMessage()    {}
func (*ProvisioningUuid) Descriptor() ([]byte, []int) {
	return fileDescriptor_e841f1d488a49120, []int{0}
}

func (m *ProvisioningUuid) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProvisioningUuid.Unmarshal(m, b)
}
func (m *ProvisioningUuid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProvisioningUuid.Marshal(b, m, deterministic)
}
func (m *ProvisioningUuid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProvisioningUuid.Merge(m, src)
}
func (m *ProvisioningUuid) XXX_Size() int {
	return xxx_messageInfo_ProvisioningUuid.Size(m)
}
func (m *ProvisioningUuid) XXX_DiscardUnknown() {
	xxx_messageInfo_ProvisioningUuid.DiscardUnknown(m)
}

var xxx_messageInfo_ProvisioningUuid proto.InternalMessageInfo

func (m *ProvisioningUuid) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

type ProvisionEnvelope struct {
	PublicKey            []byte   `protobuf:"bytes,1,opt,name=publicKey" json:"publicKey,omitempty"`
	Body                 []byte   `protobuf:"bytes,2,opt,name=body" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProvisionEnvelope) Reset()         { *m = ProvisionEnvelope{} }
func (m *ProvisionEnvelope) String() string { return proto.CompactTextString(m) }
func (*ProvisionEnvelope) ProtoMessage()    {}
func (*ProvisionEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_e841f1d488a49120, []int{1}
}

func (m *ProvisionEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProvisionEnvelope.Unmarshal(m, b)
}
func (m *ProvisionEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProvisionEnvelope.Marshal(b, m, deterministic)
}
func (m *ProvisionEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProvisionEnvelope.Merge(m, src)
}
func (m *ProvisionEnvelope) XXX_Size() int {
	return xxx_messageInfo_ProvisionEnvelope.Size(m)
}
func (m *ProvisionEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_ProvisionEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_ProvisionEnvelope proto.InternalMessageInfo

func (m *ProvisionEnvelope) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *ProvisionEnvelope) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type ProvisionMessage struct {
	IdentityKeyPrivate   []byte   `protobuf:"bytes,2,opt,name=identityKeyPrivate" json:"identityKeyPrivate,omitempty"`
	Number               *string  `protobuf:"bytes,3,opt,name=number" json:"number,omitempty"`
	ProvisioningCode     *string  `protobuf:"bytes,4,opt,name=provisioningCode" json:"provisioningCode,omitempty"`
	UserAgent            *string  `protobuf:"bytes,5,opt,name=userAgent" json:"userAgent,omitempty"`
	ProfileKey           []byte   `protobuf:"bytes,6,opt,name=profileKey" json:"profileKey,omitempty"`
	ReadReceipts         *bool    `protobuf:"varint,7,opt,name=readReceipts" json:"readReceipts,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProvisionMessage) Reset()         { *m = ProvisionMessage{} }
func (m *ProvisionMessage) String() string { return proto.CompactTextString(m) }
func (*ProvisionMessage) ProtoMessage()    {}
func (*ProvisionMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_e841f1d488a49120, []int{2}
}

func (m *ProvisionMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProvisionMessage.Unmarshal(m, b)
}
func (m *ProvisionMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProvisionMessage.Marshal(b, m, deterministic)
}
func (m *ProvisionMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProvisionMessage.Merge(m, src)
}
func (m *ProvisionMessage) XXX_Size() int {
	return xxx_messageInfo_ProvisionMessage.Size(m)
}
func (m *ProvisionMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ProvisionMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ProvisionMessage proto.InternalMessageInfo

func (m *ProvisionMessage) GetIdentityKeyPrivate() []byte {
	if m != nil {
		return m.IdentityKeyPrivate
	}
	return nil
}

func (m *ProvisionMessage) GetNumber() string {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return ""
}

func (m *ProvisionMessage) GetProvisioningCode() string {
	if m != nil && m.ProvisioningCode != nil {
		return *m.ProvisioningCode
	}
	return ""
}

func (m *ProvisionMessage) GetUserAgent() string {
	if m != nil && m.UserAgent != nil {
		return *m.UserAgent
	}
	return ""
}

func (m *ProvisionMessage) GetProfileKey() []byte {
	if m != nil {
		return m.ProfileKey
	}
	return nil
}

func (m *ProvisionMessage) GetReadReceipts() bool {
	if m != nil && m.ReadReceipts != nil {
		return *m.ReadReceipts
	}
	return false
}

func init() {
	proto.RegisterType((*ProvisioningUuid)(nil), "signalservice.ProvisioningUuid")
	proto.RegisterType((*ProvisionEnvelope)(nil), "signalservice.ProvisionEnvelope")
	proto.RegisterType((*ProvisionMessage)(nil), "signalservice.ProvisionMessage")
}

func init() { proto.RegisterFile("DeviceMessages.proto", fileDescriptor_e841f1d488a49120) }

var fileDescriptor_e841f1d488a49120 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xdf, 0x4a, 0xc3, 0x30,
	0x14, 0xc6, 0xa9, 0xce, 0xe9, 0x0e, 0x13, 0x66, 0x10, 0xc9, 0x85, 0xc8, 0xe8, 0x85, 0x0c, 0x2f,
	0xf6, 0x0e, 0xa2, 0xbb, 0x1a, 0xc2, 0x28, 0xf8, 0x00, 0x6d, 0xf3, 0x59, 0x0e, 0xd4, 0x24, 0xe4,
	0x4f, 0xa1, 0x6f, 0xec, 0x63, 0x48, 0x43, 0xb1, 0x15, 0xbd, 0x3b, 0xf9, 0x7d, 0x5f, 0x92, 0xdf,
	0xa1, 0xdb, 0x57, 0x74, 0x5c, 0xe3, 0x0d, 0xde, 0x97, 0x0d, 0xfc, 0xde, 0x3a, 0x13, 0x8c, 0xb8,
	0xf6, 0xdc, 0xe8, 0xb2, 0xf5, 0x70, 0x43, 0x98, 0x3f, 0xd2, 0xe6, 0xe4, 0x4c, 0xc7, 0x9e, 0x8d,
	0x66, 0xdd, 0xbc, 0x47, 0x56, 0x42, 0xd0, 0x22, 0x46, 0x56, 0x32, 0xdb, 0x66, 0xbb, 0x55, 0x91,
	0xe6, 0xfc, 0x40, 0x37, 0x3f, 0xbd, 0x83, 0xee, 0xd0, 0x1a, 0x0b, 0x71, 0x4f, 0x2b, 0x1b, 0xab,
	0x96, 0xeb, 0x23, 0xfa, 0xd4, 0x5e, 0x17, 0x13, 0x18, 0x9e, 0xa9, 0x8c, 0xea, 0xe5, 0x59, 0x0a,
	0xd2, 0x9c, 0x7f, 0x65, 0xb3, 0xff, 0x46, 0x33, 0xb1, 0x27, 0xc1, 0x0a, 0x3a, 0x70, 0xe8, 0x8f,
	0xe8, 0x4f, 0x8e, 0xbb, 0x32, 0x60, 0xbc, 0xf6, 0x4f, 0x22, 0xee, 0x68, 0xa9, 0xe3, 0x67, 0x05,
	0x27, 0xcf, 0x93, 0xe1, 0x78, 0x12, 0x4f, 0xb4, 0xb1, 0xb3, 0x5d, 0x5e, 0x8c, 0x82, 0x5c, 0xa4,
	0xc6, 0x1f, 0x3e, 0xa8, 0x47, 0x0f, 0xf7, 0xdc, 0x40, 0x07, 0x79, 0x91, 0x4a, 0x13, 0x10, 0x0f,
	0x44, 0xd6, 0x99, 0x0f, 0x6e, 0x31, 0x6c, 0xb6, 0x4c, 0x26, 0x33, 0x22, 0x72, 0x5a, 0x3b, 0x94,
	0xaa, 0x40, 0x0d, 0xb6, 0xc1, 0xcb, 0xcb, 0x6d, 0xb6, 0xbb, 0x2a, 0x7e, 0xb1, 0xef, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x5c, 0xff, 0xe4, 0x9e, 0x7f, 0x01, 0x00, 0x00,
}
