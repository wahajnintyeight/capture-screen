// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0--rc2
// source: capture-screen-request.proto

package output

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

type ScreenCaptureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceName  string `protobuf:"bytes,1,opt,name=deviceName,proto3" json:"deviceName,omitempty"`
	TimesTamp   string `protobuf:"bytes,2,opt,name=timesTamp,proto3" json:"timesTamp,omitempty"`
	OsName      string `protobuf:"bytes,3,opt,name=osName,proto3" json:"osName,omitempty"`
	MemoryUsage string `protobuf:"bytes,5,opt,name=memoryUsage,proto3" json:"memoryUsage,omitempty"`
	DiskUsage   string `protobuf:"bytes,6,opt,name=diskUsage,proto3" json:"diskUsage,omitempty"`
	LastImage   string `protobuf:"bytes,7,opt,name=lastImage,proto3" json:"lastImage,omitempty"`
}

func (x *ScreenCaptureRequest) Reset() {
	*x = ScreenCaptureRequest{}
	mi := &file_capture_screen_request_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScreenCaptureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScreenCaptureRequest) ProtoMessage() {}

func (x *ScreenCaptureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_capture_screen_request_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScreenCaptureRequest.ProtoReflect.Descriptor instead.
func (*ScreenCaptureRequest) Descriptor() ([]byte, []int) {
	return file_capture_screen_request_proto_rawDescGZIP(), []int{0}
}

func (x *ScreenCaptureRequest) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

func (x *ScreenCaptureRequest) GetTimesTamp() string {
	if x != nil {
		return x.TimesTamp
	}
	return ""
}

func (x *ScreenCaptureRequest) GetOsName() string {
	if x != nil {
		return x.OsName
	}
	return ""
}

func (x *ScreenCaptureRequest) GetMemoryUsage() string {
	if x != nil {
		return x.MemoryUsage
	}
	return ""
}

func (x *ScreenCaptureRequest) GetDiskUsage() string {
	if x != nil {
		return x.DiskUsage
	}
	return ""
}

func (x *ScreenCaptureRequest) GetLastImage() string {
	if x != nil {
		return x.LastImage
	}
	return ""
}

type ScreenCaptureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ScreenCaptureResponse) Reset() {
	*x = ScreenCaptureResponse{}
	mi := &file_capture_screen_request_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScreenCaptureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScreenCaptureResponse) ProtoMessage() {}

func (x *ScreenCaptureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_capture_screen_request_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScreenCaptureResponse.ProtoReflect.Descriptor instead.
func (*ScreenCaptureResponse) Descriptor() ([]byte, []int) {
	return file_capture_screen_request_proto_rawDescGZIP(), []int{1}
}

func (x *ScreenCaptureResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ScreenCaptureResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_capture_screen_request_proto protoreflect.FileDescriptor

var file_capture_screen_request_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x2d, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e,
	0x2d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x22, 0xca, 0x01,
	0x0a, 0x14, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x54,
	0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x54, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x4b, 0x0a, 0x15, 0x53, 0x63,
	0x72, 0x65, 0x65, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x70, 0x0a, 0x14, 0x53, 0x63, 0x72, 0x65, 0x65,
	0x6e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x58, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12, 0x23,
	0x2e, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x53,
	0x63, 0x72, 0x65, 0x65, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x63, 0x61, 0x70, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a, 0x19, 0x63, 0x61, 0x70,
	0x74, 0x75, 0x72, 0x65, 0x2d, 0x73, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x2f, 0x73, 0x72, 0x63, 0x2f,
	0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_capture_screen_request_proto_rawDescOnce sync.Once
	file_capture_screen_request_proto_rawDescData = file_capture_screen_request_proto_rawDesc
)

func file_capture_screen_request_proto_rawDescGZIP() []byte {
	file_capture_screen_request_proto_rawDescOnce.Do(func() {
		file_capture_screen_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_capture_screen_request_proto_rawDescData)
	})
	return file_capture_screen_request_proto_rawDescData
}

var file_capture_screen_request_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_capture_screen_request_proto_goTypes = []any{
	(*ScreenCaptureRequest)(nil),  // 0: screencapture.ScreenCaptureRequest
	(*ScreenCaptureResponse)(nil), // 1: screencapture.ScreenCaptureResponse
}
var file_capture_screen_request_proto_depIdxs = []int32{
	0, // 0: screencapture.ScreenCaptureService.SendCapture:input_type -> screencapture.ScreenCaptureRequest
	1, // 1: screencapture.ScreenCaptureService.SendCapture:output_type -> screencapture.ScreenCaptureResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_capture_screen_request_proto_init() }
func file_capture_screen_request_proto_init() {
	if File_capture_screen_request_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_capture_screen_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_capture_screen_request_proto_goTypes,
		DependencyIndexes: file_capture_screen_request_proto_depIdxs,
		MessageInfos:      file_capture_screen_request_proto_msgTypes,
	}.Build()
	File_capture_screen_request_proto = out.File
	file_capture_screen_request_proto_rawDesc = nil
	file_capture_screen_request_proto_goTypes = nil
	file_capture_screen_request_proto_depIdxs = nil
}
