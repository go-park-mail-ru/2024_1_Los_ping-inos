// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: image.proto

package proto

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

type GetImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestID int64  `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Id        int64  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`    // User id
	Cell      string `protobuf:"bytes,3,opt,name=cell,proto3" json:"cell,omitempty"` // Number of the cell that contains the image
}

func (x *GetImageRequest) Reset() {
	*x = GetImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageRequest) ProtoMessage() {}

func (x *GetImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageRequest.ProtoReflect.Descriptor instead.
func (*GetImageRequest) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{0}
}

func (x *GetImageRequest) GetRequestID() int64 {
	if x != nil {
		return x.RequestID
	}
	return 0
}

func (x *GetImageRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetImageRequest) GetCell() string {
	if x != nil {
		return x.Cell
	}
	return ""
}

type GetImageResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetImageResponce) Reset() {
	*x = GetImageResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageResponce) ProtoMessage() {}

func (x *GetImageResponce) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageResponce.ProtoReflect.Descriptor instead.
func (*GetImageResponce) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{1}
}

func (x *GetImageResponce) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_image_proto protoreflect.FileDescriptor

var file_image_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x22, 0x53, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x6c, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x6c, 0x6c, 0x22, 0x24, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x32,
	0x44, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x12, 0x16, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x63, 0x65, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_image_proto_rawDescOnce sync.Once
	file_image_proto_rawDescData = file_image_proto_rawDesc
)

func file_image_proto_rawDescGZIP() []byte {
	file_image_proto_rawDescOnce.Do(func() {
		file_image_proto_rawDescData = protoimpl.X.CompressGZIP(file_image_proto_rawDescData)
	})
	return file_image_proto_rawDescData
}

var file_image_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_image_proto_goTypes = []interface{}{
	(*GetImageRequest)(nil),  // 0: image.GetImageRequest
	(*GetImageResponce)(nil), // 1: image.GetImageResponce
}
var file_image_proto_depIdxs = []int32{
	0, // 0: image.Image.GetImage:input_type -> image.GetImageRequest
	1, // 1: image.Image.GetImage:output_type -> image.GetImageResponce
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_image_proto_init() }
func file_image_proto_init() {
	if File_image_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_image_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageRequest); i {
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
		file_image_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageResponce); i {
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
			RawDescriptor: file_image_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_image_proto_goTypes,
		DependencyIndexes: file_image_proto_depIdxs,
		MessageInfos:      file_image_proto_msgTypes,
	}.Build()
	File_image_proto = out.File
	file_image_proto_rawDesc = nil
	file_image_proto_goTypes = nil
	file_image_proto_depIdxs = nil
}