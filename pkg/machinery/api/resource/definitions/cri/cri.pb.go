// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: resource/definitions/cri/cri.proto

package cri

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SeccompProfileSpec represents the SeccompProfile.
type SeccompProfileSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value *structpb.Struct `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *SeccompProfileSpec) Reset() {
	*x = SeccompProfileSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_resource_definitions_cri_cri_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeccompProfileSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeccompProfileSpec) ProtoMessage() {}

func (x *SeccompProfileSpec) ProtoReflect() protoreflect.Message {
	mi := &file_resource_definitions_cri_cri_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeccompProfileSpec.ProtoReflect.Descriptor instead.
func (*SeccompProfileSpec) Descriptor() ([]byte, []int) {
	return file_resource_definitions_cri_cri_proto_rawDescGZIP(), []int{0}
}

func (x *SeccompProfileSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SeccompProfileSpec) GetValue() *structpb.Struct {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_resource_definitions_cri_cri_proto protoreflect.FileDescriptor

var file_resource_definitions_cri_cri_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x72, 0x69, 0x2f, 0x63, 0x72, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x63, 0x72, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x57, 0x0a, 0x12, 0x53, 0x65, 0x63, 0x63, 0x6f, 0x6d, 0x70, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x48, 0x5a, 0x46, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x64, 0x65, 0x72, 0x6f,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d,
	0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x63, 0x72, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resource_definitions_cri_cri_proto_rawDescOnce sync.Once
	file_resource_definitions_cri_cri_proto_rawDescData = file_resource_definitions_cri_cri_proto_rawDesc
)

func file_resource_definitions_cri_cri_proto_rawDescGZIP() []byte {
	file_resource_definitions_cri_cri_proto_rawDescOnce.Do(func() {
		file_resource_definitions_cri_cri_proto_rawDescData = protoimpl.X.CompressGZIP(file_resource_definitions_cri_cri_proto_rawDescData)
	})
	return file_resource_definitions_cri_cri_proto_rawDescData
}

var file_resource_definitions_cri_cri_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resource_definitions_cri_cri_proto_goTypes = []interface{}{
	(*SeccompProfileSpec)(nil), // 0: talos.resource.definitions.cri.SeccompProfileSpec
	(*structpb.Struct)(nil),    // 1: google.protobuf.Struct
}
var file_resource_definitions_cri_cri_proto_depIdxs = []int32{
	1, // 0: talos.resource.definitions.cri.SeccompProfileSpec.value:type_name -> google.protobuf.Struct
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_resource_definitions_cri_cri_proto_init() }
func file_resource_definitions_cri_cri_proto_init() {
	if File_resource_definitions_cri_cri_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_resource_definitions_cri_cri_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeccompProfileSpec); i {
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
			RawDescriptor: file_resource_definitions_cri_cri_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resource_definitions_cri_cri_proto_goTypes,
		DependencyIndexes: file_resource_definitions_cri_cri_proto_depIdxs,
		MessageInfos:      file_resource_definitions_cri_cri_proto_msgTypes,
	}.Build()
	File_resource_definitions_cri_cri_proto = out.File
	file_resource_definitions_cri_cri_proto_rawDesc = nil
	file_resource_definitions_cri_cri_proto_goTypes = nil
	file_resource_definitions_cri_cri_proto_depIdxs = nil
}
