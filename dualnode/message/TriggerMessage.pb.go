// Code generated by protoc-gen-go. DO NOT EDIT.
// source: TriggerMessage.proto

package message

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
//const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// TriggerMessage is message sent from Kardia to dual node to have it execute specific method based on given address, method and params
// After finish executing, txid will be appended into params within every callBack in callBacks and
// they are sent back to Kardia
type TriggerMessage struct {
	ContractAddress      string            `protobuf:"bytes,1,opt,name=contractAddress,proto3" json:"contractAddress,omitempty"`
	MethodName           string            `protobuf:"bytes,2,opt,name=methodName,proto3" json:"methodName,omitempty"`
	Params               []string          `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
	CallBacks            []*TriggerMessage `protobuf:"bytes,4,rep,name=callBacks,proto3" json:"callBacks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TriggerMessage) Reset()         { *m = TriggerMessage{} }
func (m *TriggerMessage) String() string { return proto.CompactTextString(m) }
func (*TriggerMessage) ProtoMessage()    {}
func (*TriggerMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_034f86483671695c, []int{0}
}

func (m *TriggerMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TriggerMessage.Unmarshal(m, b)
}
func (m *TriggerMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TriggerMessage.Marshal(b, m, deterministic)
}
func (m *TriggerMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TriggerMessage.Merge(m, src)
}
func (m *TriggerMessage) XXX_Size() int {
	return xxx_messageInfo_TriggerMessage.Size(m)
}
func (m *TriggerMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TriggerMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TriggerMessage proto.InternalMessageInfo

func (m *TriggerMessage) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func (m *TriggerMessage) GetMethodName() string {
	if m != nil {
		return m.MethodName
	}
	return ""
}

func (m *TriggerMessage) GetParams() []string {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *TriggerMessage) GetCallBacks() []*TriggerMessage {
	if m != nil {
		return m.CallBacks
	}
	return nil
}

func init() {
	proto.RegisterType((*TriggerMessage)(nil), "protocol.TriggerMessage")
}

func init() { proto.RegisterFile("TriggerMessage.proto", fileDescriptor_034f86483671695c) }

var fileDescriptor_034f86483671695c = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x09, 0x29, 0xca, 0x4c,
	0x4f, 0x4f, 0x2d, 0xf2, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x00, 0x53, 0xc9, 0xf9, 0x39, 0x4a, 0xab, 0x18, 0xb9, 0xf8, 0x50, 0x95, 0x08, 0x69,
	0x70, 0xf1, 0x27, 0xe7, 0xe7, 0x95, 0x14, 0x25, 0x26, 0x97, 0x38, 0xa6, 0xa4, 0x14, 0xa5, 0x16,
	0x17, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xa1, 0x0b, 0x0b, 0xc9, 0x71, 0x71, 0xe5, 0xa6,
	0x96, 0x64, 0xe4, 0xa7, 0xf8, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x81, 0x15, 0x21, 0x89, 0x08, 0x89,
	0x71, 0xb1, 0x15, 0x24, 0x16, 0x25, 0xe6, 0x16, 0x4b, 0x30, 0x2b, 0x30, 0x6b, 0x70, 0x06, 0x41,
	0x79, 0x42, 0x66, 0x5c, 0x9c, 0xc9, 0x89, 0x39, 0x39, 0x4e, 0x89, 0xc9, 0xd9, 0xc5, 0x12, 0x2c,
	0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x12, 0x7a, 0x30, 0x27, 0xe9, 0xa1, 0x3a, 0x27, 0x08, 0xa1, 0xd4,
	0x49, 0x86, 0x8b, 0x23, 0xbf, 0x28, 0x5d, 0xaf, 0xa4, 0x28, 0x3f, 0xcf, 0x89, 0x1d, 0xaa, 0x2c,
	0x8a, 0x3d, 0x17, 0xa2, 0x30, 0x89, 0x0d, 0x6c, 0x82, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x9f,
	0x6d, 0xa4, 0xb2, 0xf3, 0x00, 0x00, 0x00,
}
