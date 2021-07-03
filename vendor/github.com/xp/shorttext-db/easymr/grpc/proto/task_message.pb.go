// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task_message.proto

package proto

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

type Task struct {
	Type                 uint32       `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Priority             uint32       `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	Consumable           string       `protobuf:"bytes,3,opt,name=consumable,proto3" json:"consumable,omitempty"`
	Source               *Collection  `protobuf:"bytes,4,opt,name=source,proto3" json:"source,omitempty"`
	Result               *Collection  `protobuf:"bytes,5,opt,name=result,proto3" json:"result,omitempty"`
	Context              *TaskContext `protobuf:"bytes,6,opt,name=context,proto3" json:"context,omitempty"`
	Stage                uint32       `protobuf:"varint,7,opt,name=stage,proto3" json:"stage,omitempty"`
	RunType              uint32       `protobuf:"varint,8,opt,name=RunType,proto3" json:"RunType,omitempty"`
	Object               *ObjectItem  `protobuf:"bytes,9,opt,name=Object,proto3" json:"Object,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_c29fba0774a43188, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Task) GetPriority() uint32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *Task) GetConsumable() string {
	if m != nil {
		return m.Consumable
	}
	return ""
}

func (m *Task) GetSource() *Collection {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Task) GetResult() *Collection {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *Task) GetContext() *TaskContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *Task) GetStage() uint32 {
	if m != nil {
		return m.Stage
	}
	return 0
}

func (m *Task) GetRunType() uint32 {
	if m != nil {
		return m.RunType
	}
	return 0
}

func (m *Task) GetObject() *ObjectItem {
	if m != nil {
		return m.Object
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "proto.Task")
}

func init() { proto.RegisterFile("task_message.proto", fileDescriptor_c29fba0774a43188) }

var fileDescriptor_c29fba0774a43188 = []byte{
	// 286 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x4f, 0x4b, 0xc4, 0x30,
	0x10, 0xc5, 0xe9, 0xda, 0x3f, 0xbb, 0x11, 0x0f, 0x86, 0x3d, 0xc4, 0x1e, 0xa4, 0x08, 0x62, 0x05,
	0x69, 0x41, 0xbf, 0x81, 0x7b, 0xf2, 0x24, 0x84, 0x3d, 0x79, 0x91, 0x34, 0x8c, 0xa5, 0x6e, 0xd3,
	0x09, 0x49, 0x0a, 0xf6, 0x83, 0xf9, 0xfd, 0xa4, 0x69, 0x77, 0x71, 0x17, 0xf6, 0x34, 0x7d, 0x6f,
	0x7e, 0x43, 0xdf, 0x0b, 0xa1, 0x4e, 0xd8, 0xdd, 0xa7, 0x02, 0x6b, 0x45, 0x0d, 0x85, 0x36, 0xe8,
	0x90, 0x46, 0x7e, 0xa4, 0x4c, 0x62, 0xdb, 0x82, 0x74, 0x0d, 0x76, 0xc7, 0x40, 0x7a, 0x33, 0x1e,
	0x49, 0xec, 0x1c, 0xfc, 0xb8, 0x93, 0xd5, 0x1a, 0xab, 0x6f, 0x90, 0x27, 0xee, 0xdd, 0xef, 0x82,
	0x84, 0x5b, 0x61, 0x77, 0x94, 0x92, 0xd0, 0x0d, 0x1a, 0x58, 0x90, 0x05, 0xf9, 0x15, 0xf7, 0xdf,
	0x34, 0x25, 0x4b, 0x6d, 0x1a, 0x34, 0x8d, 0x1b, 0xd8, 0xc2, 0xfb, 0x07, 0x4d, 0x6f, 0x09, 0x91,
	0xd8, 0xd9, 0x5e, 0x89, 0xaa, 0x05, 0x76, 0x91, 0x05, 0xf9, 0x8a, 0xff, 0x73, 0xe8, 0x23, 0x89,
	0x2d, 0xf6, 0x46, 0x02, 0x0b, 0xb3, 0x20, 0xbf, 0x7c, 0xbe, 0x9e, 0x7e, 0x58, 0x6c, 0x0e, 0xd1,
	0xf9, 0x0c, 0x8c, 0xa8, 0x01, 0xdb, 0xb7, 0x8e, 0x45, 0x67, 0xd1, 0x09, 0xa0, 0x4f, 0x24, 0x99,
	0xdb, 0xb1, 0xd8, 0xb3, 0x74, 0x66, 0xc7, 0x0e, 0x9b, 0x69, 0xc3, 0xf7, 0x08, 0x5d, 0x93, 0xc8,
	0x3a, 0x51, 0x03, 0x4b, 0x7c, 0xf8, 0x49, 0x50, 0x46, 0x12, 0xde, 0x77, 0xdb, 0xb1, 0xec, 0xd2,
	0xfb, 0x7b, 0x39, 0x06, 0x79, 0xf7, 0x8f, 0xc4, 0x56, 0x47, 0x41, 0x26, 0xf3, 0xcd, 0x81, 0xe2,
	0x33, 0xf0, 0xfa, 0xf0, 0x71, 0x2f, 0x51, 0x15, 0x1d, 0x80, 0x2e, 0x6b, 0xd4, 0xad, 0x70, 0x5f,
	0x68, 0x54, 0x09, 0xc2, 0x0e, 0xca, 0x94, 0xb5, 0xd1, 0xb2, 0xf4, 0xb7, 0x55, 0xec, 0xc7, 0xcb,
	0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x24, 0x5b, 0x68, 0xf5, 0xcf, 0x01, 0x00, 0x00,
}