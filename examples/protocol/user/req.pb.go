// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.15.4
// source: user/req.proto

package user

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GetUserReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   int64  `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	UserName string `protobuf:"bytes,2,opt,name=userName,proto3" json:"userName,omitempty"`
}

func (x *GetUserReq) Reset() {
	*x = GetUserReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_req_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserReq) ProtoMessage() {}

func (x *GetUserReq) ProtoReflect() protoreflect.Message {
	mi := &file_user_req_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserReq.ProtoReflect.Descriptor instead.
func (*GetUserReq) Descriptor() ([]byte, []int) {
	return file_user_req_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetUserReq) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

var File_user_req_proto protoreflect.FileDescriptor

var file_user_req_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x40, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x20, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x50, 0x01, 0x5a, 0x0d, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_user_req_proto_rawDescOnce sync.Once
	file_user_req_proto_rawDescData = file_user_req_proto_rawDesc
)

func file_user_req_proto_rawDescGZIP() []byte {
	file_user_req_proto_rawDescOnce.Do(func() {
		file_user_req_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_req_proto_rawDescData)
	})
	return file_user_req_proto_rawDescData
}

var file_user_req_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_user_req_proto_goTypes = []interface{}{
	(*GetUserReq)(nil), // 0: user.GetUserReq
}
var file_user_req_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_req_proto_init() }
func file_user_req_proto_init() {
	if File_user_req_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_req_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserReq); i {
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
			RawDescriptor: file_user_req_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_req_proto_goTypes,
		DependencyIndexes: file_user_req_proto_depIdxs,
		MessageInfos:      file_user_req_proto_msgTypes,
	}.Build()
	File_user_req_proto = out.File
	file_user_req_proto_rawDesc = nil
	file_user_req_proto_goTypes = nil
	file_user_req_proto_depIdxs = nil
}
