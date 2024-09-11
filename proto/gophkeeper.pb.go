// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.9
// source: proto/gophkeeper.proto

package keeperproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Data_DType int32

const (
	Data_CREDENTIALS Data_DType = 0
	Data_TEXT        Data_DType = 1
	Data_BINARY      Data_DType = 2
	Data_CARD        Data_DType = 3
)

// Enum value maps for Data_DType.
var (
	Data_DType_name = map[int32]string{
		0: "CREDENTIALS",
		1: "TEXT",
		2: "BINARY",
		3: "CARD",
	}
	Data_DType_value = map[string]int32{
		"CREDENTIALS": 0,
		"TEXT":        1,
		"BINARY":      2,
		"CARD":        3,
	}
)

func (x Data_DType) Enum() *Data_DType {
	p := new(Data_DType)
	*p = x
	return p
}

func (x Data_DType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Data_DType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_gophkeeper_proto_enumTypes[0].Descriptor()
}

func (Data_DType) Type() protoreflect.EnumType {
	return &file_proto_gophkeeper_proto_enumTypes[0]
}

func (x Data_DType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Data_DType.Descriptor instead.
func (Data_DType) EnumDescriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{0, 0}
}

type Data_OpType int32

const (
	Data_ADD    Data_OpType = 0
	Data_DELETE Data_OpType = 1
)

// Enum value maps for Data_OpType.
var (
	Data_OpType_name = map[int32]string{
		0: "ADD",
		1: "DELETE",
	}
	Data_OpType_value = map[string]int32{
		"ADD":    0,
		"DELETE": 1,
	}
)

func (x Data_OpType) Enum() *Data_OpType {
	p := new(Data_OpType)
	*p = x
	return p
}

func (x Data_OpType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Data_OpType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_gophkeeper_proto_enumTypes[1].Descriptor()
}

func (Data_OpType) Type() protoreflect.EnumType {
	return &file_proto_gophkeeper_proto_enumTypes[1]
}

func (x Data_OpType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Data_OpType.Descriptor instead.
func (Data_OpType) EnumDescriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{0, 1}
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dtype   Data_DType  `protobuf:"varint,2,opt,name=dtype,proto3,enum=gophkeeper.Data_DType" json:"dtype,omitempty"`
	Optype  Data_OpType `protobuf:"varint,3,opt,name=optype,proto3,enum=gophkeeper.Data_OpType" json:"optype,omitempty"`
	Payload []byte      `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gophkeeper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gophkeeper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Data) GetDtype() Data_DType {
	if x != nil {
		return x.Dtype
	}
	return Data_CREDENTIALS
}

func (x *Data) GetOptype() Data_OpType {
	if x != nil {
		return x.Optype
	}
	return Data_ADD
}

func (x *Data) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type SyncPushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *SyncPushRequest) Reset() {
	*x = SyncPushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gophkeeper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncPushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncPushRequest) ProtoMessage() {}

func (x *SyncPushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gophkeeper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncPushRequest.ProtoReflect.Descriptor instead.
func (*SyncPushRequest) Descriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{1}
}

func (x *SyncPushRequest) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type SyncPushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *SyncPushResponse) Reset() {
	*x = SyncPushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gophkeeper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncPushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncPushResponse) ProtoMessage() {}

func (x *SyncPushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gophkeeper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncPushResponse.ProtoReflect.Descriptor instead.
func (*SyncPushResponse) Descriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{2}
}

func (x *SyncPushResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type SyncPullRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name []string `protobuf:"bytes,1,rep,name=name,proto3" json:"name,omitempty"`
}

func (x *SyncPullRequest) Reset() {
	*x = SyncPullRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gophkeeper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncPullRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncPullRequest) ProtoMessage() {}

func (x *SyncPullRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gophkeeper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncPullRequest.ProtoReflect.Descriptor instead.
func (*SyncPullRequest) Descriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{3}
}

func (x *SyncPullRequest) GetName() []string {
	if x != nil {
		return x.Name
	}
	return nil
}

type SyncPullResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *SyncPullResponse) Reset() {
	*x = SyncPullResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_gophkeeper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncPullResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncPullResponse) ProtoMessage() {}

func (x *SyncPullResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_gophkeeper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncPullResponse.ProtoReflect.Descriptor instead.
func (*SyncPullResponse) Descriptor() ([]byte, []int) {
	return file_proto_gophkeeper_proto_rawDescGZIP(), []int{4}
}

func (x *SyncPullResponse) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_gophkeeper_proto protoreflect.FileDescriptor

var file_proto_gophkeeper_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x22, 0xec, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x2c, 0x0a, 0x05, 0x64, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x44, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x64, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x6f, 0x70, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x17, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x2e, 0x4f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x6f, 0x70, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x38, 0x0a, 0x05, 0x44, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x52, 0x45, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x41,
	0x4c, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x58, 0x54, 0x10, 0x01, 0x12, 0x0a,
	0x0a, 0x06, 0x42, 0x49, 0x4e, 0x41, 0x52, 0x59, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x41,
	0x52, 0x44, 0x10, 0x03, 0x22, 0x1d, 0x0a, 0x06, 0x4f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07,
	0x0a, 0x03, 0x41, 0x44, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x10, 0x01, 0x22, 0x37, 0x0a, 0x0f, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x28, 0x0a, 0x10,
	0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x25, 0x0a, 0x0f, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75,
	0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x38, 0x0a,
	0x10, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x8c, 0x01, 0x0a, 0x04, 0x53, 0x79, 0x6e, 0x63,
	0x12, 0x41, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12, 0x1b, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x04, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x1b, 0x2e, 0x67, 0x6f,
	0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x6c,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x6c, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x63, 0x64, 0x2f, 0x67, 0x6f,
	0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_gophkeeper_proto_rawDescOnce sync.Once
	file_proto_gophkeeper_proto_rawDescData = file_proto_gophkeeper_proto_rawDesc
)

func file_proto_gophkeeper_proto_rawDescGZIP() []byte {
	file_proto_gophkeeper_proto_rawDescOnce.Do(func() {
		file_proto_gophkeeper_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_gophkeeper_proto_rawDescData)
	})
	return file_proto_gophkeeper_proto_rawDescData
}

var file_proto_gophkeeper_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_gophkeeper_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_gophkeeper_proto_goTypes = []any{
	(Data_DType)(0),          // 0: gophkeeper.Data.DType
	(Data_OpType)(0),         // 1: gophkeeper.Data.OpType
	(*Data)(nil),             // 2: gophkeeper.Data
	(*SyncPushRequest)(nil),  // 3: gophkeeper.SyncPushRequest
	(*SyncPushResponse)(nil), // 4: gophkeeper.SyncPushResponse
	(*SyncPullRequest)(nil),  // 5: gophkeeper.SyncPullRequest
	(*SyncPullResponse)(nil), // 6: gophkeeper.SyncPullResponse
}
var file_proto_gophkeeper_proto_depIdxs = []int32{
	0, // 0: gophkeeper.Data.dtype:type_name -> gophkeeper.Data.DType
	1, // 1: gophkeeper.Data.optype:type_name -> gophkeeper.Data.OpType
	2, // 2: gophkeeper.SyncPushRequest.data:type_name -> gophkeeper.Data
	2, // 3: gophkeeper.SyncPullResponse.data:type_name -> gophkeeper.Data
	3, // 4: gophkeeper.Sync.Push:input_type -> gophkeeper.SyncPushRequest
	5, // 5: gophkeeper.Sync.Pull:input_type -> gophkeeper.SyncPullRequest
	4, // 6: gophkeeper.Sync.Push:output_type -> gophkeeper.SyncPushResponse
	6, // 7: gophkeeper.Sync.Pull:output_type -> gophkeeper.SyncPullResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_gophkeeper_proto_init() }
func file_proto_gophkeeper_proto_init() {
	if File_proto_gophkeeper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_gophkeeper_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Data); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_gophkeeper_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SyncPushRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_gophkeeper_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SyncPushResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_gophkeeper_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SyncPullRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_gophkeeper_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SyncPullResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_gophkeeper_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_gophkeeper_proto_goTypes,
		DependencyIndexes: file_proto_gophkeeper_proto_depIdxs,
		EnumInfos:         file_proto_gophkeeper_proto_enumTypes,
		MessageInfos:      file_proto_gophkeeper_proto_msgTypes,
	}.Build()
	File_proto_gophkeeper_proto = out.File
	file_proto_gophkeeper_proto_rawDesc = nil
	file_proto_gophkeeper_proto_goTypes = nil
	file_proto_gophkeeper_proto_depIdxs = nil
}
