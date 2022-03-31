// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: middleware/cache/testdata/stream.proto

package testdata

import (
	gitalypb "gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Destination *gitalypb.Repository `protobuf:"bytes,1,opt,name=destination,proto3" json:"destination,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_middleware_cache_testdata_stream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_cache_testdata_stream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_middleware_cache_testdata_stream_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetDestination() *gitalypb.Repository {
	if x != nil {
		return x.Destination
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_middleware_cache_testdata_stream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_middleware_cache_testdata_stream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_middleware_cache_testdata_stream_proto_rawDescGZIP(), []int{1}
}

var File_middleware_cache_testdata_stream_proto protoreflect.FileDescriptor

var file_middleware_cache_testdata_stream_proto_rawDesc = []byte{
	0x0a, 0x26, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61,
	0x74, 0x61, 0x1a, 0x0a, 0x6c, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67,
	0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x42, 0x04, 0x98, 0xc6, 0x2c, 0x01, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0x52, 0x0a, 0x12, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x0d, 0x49, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x64,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x04, 0xf0,
	0x97, 0x28, 0x01, 0x32, 0xb9, 0x02, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x4a, 0x0a, 0x17, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x70, 0x6f, 0x4d, 0x75, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x11,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x06, 0xfa, 0x97, 0x28, 0x02, 0x08, 0x01, 0x30, 0x01, 0x12,
	0x4b, 0x0a, 0x18, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x70, 0x6f, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x11, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x06, 0xfa, 0x97, 0x28, 0x02, 0x08, 0x02, 0x30, 0x01, 0x12, 0x47, 0x0a, 0x16,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x4d,
	0x75, 0x74, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x06, 0xfa,
	0x97, 0x28, 0x02, 0x08, 0x01, 0x12, 0x48, 0x0a, 0x17, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55,
	0x6e, 0x61, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x12, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x06, 0xfa, 0x97, 0x28, 0x02, 0x08, 0x02, 0x42,
	0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2f,
	0x76, 0x31, 0x34, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x2f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_middleware_cache_testdata_stream_proto_rawDescOnce sync.Once
	file_middleware_cache_testdata_stream_proto_rawDescData = file_middleware_cache_testdata_stream_proto_rawDesc
)

func file_middleware_cache_testdata_stream_proto_rawDescGZIP() []byte {
	file_middleware_cache_testdata_stream_proto_rawDescOnce.Do(func() {
		file_middleware_cache_testdata_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file_middleware_cache_testdata_stream_proto_rawDescData)
	})
	return file_middleware_cache_testdata_stream_proto_rawDescData
}

var file_middleware_cache_testdata_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_middleware_cache_testdata_stream_proto_goTypes = []interface{}{
	(*Request)(nil),             // 0: testdata.Request
	(*Response)(nil),            // 1: testdata.Response
	(*gitalypb.Repository)(nil), // 2: gitaly.Repository
}
var file_middleware_cache_testdata_stream_proto_depIdxs = []int32{
	2, // 0: testdata.Request.destination:type_name -> gitaly.Repository
	0, // 1: testdata.InterceptedService.IgnoredMethod:input_type -> testdata.Request
	0, // 2: testdata.TestService.ClientStreamRepoMutator:input_type -> testdata.Request
	0, // 3: testdata.TestService.ClientStreamRepoAccessor:input_type -> testdata.Request
	0, // 4: testdata.TestService.ClientUnaryRepoMutator:input_type -> testdata.Request
	0, // 5: testdata.TestService.ClientUnaryRepoAccessor:input_type -> testdata.Request
	1, // 6: testdata.InterceptedService.IgnoredMethod:output_type -> testdata.Response
	1, // 7: testdata.TestService.ClientStreamRepoMutator:output_type -> testdata.Response
	1, // 8: testdata.TestService.ClientStreamRepoAccessor:output_type -> testdata.Response
	1, // 9: testdata.TestService.ClientUnaryRepoMutator:output_type -> testdata.Response
	1, // 10: testdata.TestService.ClientUnaryRepoAccessor:output_type -> testdata.Response
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_middleware_cache_testdata_stream_proto_init() }
func file_middleware_cache_testdata_stream_proto_init() {
	if File_middleware_cache_testdata_stream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_middleware_cache_testdata_stream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_middleware_cache_testdata_stream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_middleware_cache_testdata_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_middleware_cache_testdata_stream_proto_goTypes,
		DependencyIndexes: file_middleware_cache_testdata_stream_proto_depIdxs,
		MessageInfos:      file_middleware_cache_testdata_stream_proto_msgTypes,
	}.Build()
	File_middleware_cache_testdata_stream_proto = out.File
	file_middleware_cache_testdata_stream_proto_rawDesc = nil
	file_middleware_cache_testdata_stream_proto_goTypes = nil
	file_middleware_cache_testdata_stream_proto_depIdxs = nil
}