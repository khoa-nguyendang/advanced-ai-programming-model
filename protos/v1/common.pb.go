// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: common.proto

package appsvc

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

type ResponseCode int32

const (
	ResponseCode_NO_INFO      ResponseCode = 0
	ResponseCode_SUCCESSFULLY ResponseCode = 1
	ResponseCode_FAILED       ResponseCode = 2
	ResponseCode_ERROR        ResponseCode = 3
)

// Enum value maps for ResponseCode.
var (
	ResponseCode_name = map[int32]string{
		0: "NO_INFO",
		1: "SUCCESSFULLY",
		2: "FAILED",
		3: "ERROR",
	}
	ResponseCode_value = map[string]int32{
		"NO_INFO":      0,
		"SUCCESSFULLY": 1,
		"FAILED":       2,
		"ERROR":        3,
	}
)

func (x ResponseCode) Enum() *ResponseCode {
	p := new(ResponseCode)
	*p = x
	return p
}

func (x ResponseCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseCode) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[0].Descriptor()
}

func (ResponseCode) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[0]
}

func (x ResponseCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseCode.Descriptor instead.
func (ResponseCode) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

type AdministratorState int32

const (
	AdministratorState_ACTIVATED   AdministratorState = 0
	AdministratorState_DEACTIVATED AdministratorState = 1
	AdministratorState_LIMITED     AdministratorState = 2
	AdministratorState_DELETED     AdministratorState = 3
	AdministratorState_RESTRICTED  AdministratorState = 4
)

// Enum value maps for AdministratorState.
var (
	AdministratorState_name = map[int32]string{
		0: "ACTIVATED",
		1: "DEACTIVATED",
		2: "LIMITED",
		3: "DELETED",
		4: "RESTRICTED",
	}
	AdministratorState_value = map[string]int32{
		"ACTIVATED":   0,
		"DEACTIVATED": 1,
		"LIMITED":     2,
		"DELETED":     3,
		"RESTRICTED":  4,
	}
)

func (x AdministratorState) Enum() *AdministratorState {
	p := new(AdministratorState)
	*p = x
	return p
}

func (x AdministratorState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AdministratorState) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[1].Descriptor()
}

func (AdministratorState) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[1]
}

func (x AdministratorState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AdministratorState.Descriptor instead.
func (AdministratorState) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

type UserState int32

const (
	UserState_USER_ENABLED  UserState = 0
	UserState_USER_DISABLED UserState = 1
	UserState_USER_DELETED  UserState = 2
)

// Enum value maps for UserState.
var (
	UserState_name = map[int32]string{
		0: "USER_ENABLED",
		1: "USER_DISABLED",
		2: "USER_DELETED",
	}
	UserState_value = map[string]int32{
		"USER_ENABLED":  0,
		"USER_DISABLED": 1,
		"USER_DELETED":  2,
	}
)

func (x UserState) Enum() *UserState {
	p := new(UserState)
	*p = x
	return p
}

func (x UserState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserState) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[2].Descriptor()
}

func (UserState) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[2]
}

func (x UserState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserState.Descriptor instead.
func (UserState) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

type RoleType int32

const (
	RoleType_SYSTEM_USER_ROLE RoleType = 0
	RoleType_FACEID_USER_ROLE RoleType = 1
)

// Enum value maps for RoleType.
var (
	RoleType_name = map[int32]string{
		0: "SYSTEM_USER_ROLE",
		1: "FACEID_USER_ROLE",
	}
	RoleType_value = map[string]int32{
		"SYSTEM_USER_ROLE": 0,
		"FACEID_USER_ROLE": 1,
	}
)

func (x RoleType) Enum() *RoleType {
	p := new(RoleType)
	*p = x
	return p
}

func (x RoleType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoleType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[3].Descriptor()
}

func (RoleType) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[3]
}

func (x RoleType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoleType.Descriptor instead.
func (RoleType) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

type Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data    []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	ImageId string `protobuf:"bytes,2,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *Image) Reset() {
	*x = Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Image) ProtoMessage() {}

func (x *Image) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Image.ProtoReflect.Descriptor instead.
func (*Image) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *Image) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Image) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x70, 0x70, 0x73, 0x76, 0x63, 0x22, 0x36, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x2a, 0x44,
	0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b,
	0x0a, 0x07, 0x4e, 0x4f, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x53,
	0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x46, 0x55, 0x4c, 0x4c, 0x59, 0x10, 0x01, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x03, 0x2a, 0x5e, 0x0a, 0x12, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x41, 0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x45, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4c, 0x49,
	0x4d, 0x49, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x44, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x45, 0x53, 0x54, 0x52, 0x49, 0x43, 0x54,
	0x45, 0x44, 0x10, 0x04, 0x2a, 0x42, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x4e, 0x41, 0x42, 0x4c, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x44, 0x49, 0x53, 0x41,
	0x42, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x2a, 0x36, 0x0a, 0x08, 0x52, 0x6f, 0x6c, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x46, 0x41,
	0x43, 0x45, 0x49, 0x44, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x10, 0x01,
	0x42, 0x12, 0x5a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70,
	0x70, 0x73, 0x76, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_proto_goTypes = []interface{}{
	(ResponseCode)(0),       // 0: appsvc.ResponseCode
	(AdministratorState)(0), // 1: appsvc.AdministratorState
	(UserState)(0),          // 2: appsvc.UserState
	(RoleType)(0),           // 3: appsvc.RoleType
	(*Image)(nil),           // 4: appsvc.Image
}
var file_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Image); i {
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
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		EnumInfos:         file_common_proto_enumTypes,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
