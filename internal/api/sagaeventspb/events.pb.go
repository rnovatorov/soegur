// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: internal/api/sagaeventspb/events.proto

package sagaeventspb

import (
	sagaspecpb "github.com/rnovatorov/soegur/internal/api/sagaspecpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SagaBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Spec   *sagaspecpb.Saga `protobuf:"bytes,1,opt,name=spec,proto3" json:"spec,omitempty"`
	Config *structpb.Struct `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *SagaBegun) Reset() {
	*x = SagaBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SagaBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaBegun) ProtoMessage() {}

func (x *SagaBegun) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaBegun.ProtoReflect.Descriptor instead.
func (*SagaBegun) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{0}
}

func (x *SagaBegun) GetSpec() *sagaspecpb.Saga {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *SagaBegun) GetConfig() *structpb.Struct {
	if x != nil {
		return x.Config
	}
	return nil
}

type SagaEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SagaEnded) Reset() {
	*x = SagaEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SagaEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaEnded) ProtoMessage() {}

func (x *SagaEnded) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaEnded.ProtoReflect.Descriptor instead.
func (*SagaEnded) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{1}
}

type StepBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StepBegun) Reset() {
	*x = StepBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepBegun) ProtoMessage() {}

func (x *StepBegun) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepBegun.ProtoReflect.Descriptor instead.
func (*StepBegun) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{2}
}

func (x *StepBegun) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StepEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Output *structpb.Value `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *StepEnded) Reset() {
	*x = StepEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepEnded) ProtoMessage() {}

func (x *StepEnded) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepEnded.ProtoReflect.Descriptor instead.
func (*StepEnded) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{3}
}

func (x *StepEnded) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StepEnded) GetOutput() *structpb.Value {
	if x != nil {
		return x.Output
	}
	return nil
}

type StepAborted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reason string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *StepAborted) Reset() {
	*x = StepAborted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepAborted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepAborted) ProtoMessage() {}

func (x *StepAborted) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepAborted.ProtoReflect.Descriptor instead.
func (*StepAborted) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{4}
}

func (x *StepAborted) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StepAborted) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type StepCompensationBegun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StepCompensationBegun) Reset() {
	*x = StepCompensationBegun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepCompensationBegun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepCompensationBegun) ProtoMessage() {}

func (x *StepCompensationBegun) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepCompensationBegun.ProtoReflect.Descriptor instead.
func (*StepCompensationBegun) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{5}
}

func (x *StepCompensationBegun) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StepCompensationEnded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StepCompensationEnded) Reset() {
	*x = StepCompensationEnded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StepCompensationEnded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StepCompensationEnded) ProtoMessage() {}

func (x *StepCompensationEnded) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_sagaeventspb_events_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StepCompensationEnded.ProtoReflect.Descriptor instead.
func (*StepCompensationEnded) Descriptor() ([]byte, []int) {
	return file_internal_api_sagaeventspb_events_proto_rawDescGZIP(), []int{6}
}

func (x *StepCompensationEnded) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_internal_api_sagaeventspb_events_proto protoreflect.FileDescriptor

var file_internal_api_sagaeventspb_events_proto_rawDesc = []byte{
	0x0a, 0x26, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73,
	0x61, 0x67, 0x61, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x70, 0x62, 0x2f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x61, 0x67, 0x61, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x22, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x61,
	0x67, 0x61, 0x73, 0x70, 0x65, 0x63, 0x70, 0x62, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x09, 0x53, 0x61, 0x67, 0x61, 0x42, 0x65, 0x67, 0x75, 0x6e,
	0x12, 0x2f, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x61,
	0x67, 0x61, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x53, 0x61, 0x67, 0x61, 0x52, 0x04, 0x73, 0x70, 0x65,
	0x63, 0x12, 0x2f, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x22, 0x0b, 0x0a, 0x09, 0x53, 0x61, 0x67, 0x61, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x22,
	0x1b, 0x0a, 0x09, 0x53, 0x74, 0x65, 0x70, 0x42, 0x65, 0x67, 0x75, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4b, 0x0a, 0x09,
	0x53, 0x74, 0x65, 0x70, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x35, 0x0a, 0x0b, 0x53, 0x74, 0x65,
	0x70, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x22, 0x27, 0x0a, 0x15, 0x53, 0x74, 0x65, 0x70, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0x65, 0x67, 0x75, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x27, 0x0a, 0x15, 0x53, 0x74, 0x65,
	0x70, 0x43, 0x6f, 0x6d, 0x70, 0x65, 0x6e, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64,
	0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x6e, 0x6f, 0x76, 0x61, 0x74, 0x6f, 0x72, 0x6f, 0x76, 0x2f, 0x73, 0x6f, 0x65, 0x67,
	0x75, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x73, 0x61, 0x67, 0x61, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_sagaeventspb_events_proto_rawDescOnce sync.Once
	file_internal_api_sagaeventspb_events_proto_rawDescData = file_internal_api_sagaeventspb_events_proto_rawDesc
)

func file_internal_api_sagaeventspb_events_proto_rawDescGZIP() []byte {
	file_internal_api_sagaeventspb_events_proto_rawDescOnce.Do(func() {
		file_internal_api_sagaeventspb_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_sagaeventspb_events_proto_rawDescData)
	})
	return file_internal_api_sagaeventspb_events_proto_rawDescData
}

var file_internal_api_sagaeventspb_events_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_api_sagaeventspb_events_proto_goTypes = []any{
	(*SagaBegun)(nil),             // 0: internal.api.sagaevents.SagaBegun
	(*SagaEnded)(nil),             // 1: internal.api.sagaevents.SagaEnded
	(*StepBegun)(nil),             // 2: internal.api.sagaevents.StepBegun
	(*StepEnded)(nil),             // 3: internal.api.sagaevents.StepEnded
	(*StepAborted)(nil),           // 4: internal.api.sagaevents.StepAborted
	(*StepCompensationBegun)(nil), // 5: internal.api.sagaevents.StepCompensationBegun
	(*StepCompensationEnded)(nil), // 6: internal.api.sagaevents.StepCompensationEnded
	(*sagaspecpb.Saga)(nil),       // 7: internal.api.sagaspec.Saga
	(*structpb.Struct)(nil),       // 8: google.protobuf.Struct
	(*structpb.Value)(nil),        // 9: google.protobuf.Value
}
var file_internal_api_sagaeventspb_events_proto_depIdxs = []int32{
	7, // 0: internal.api.sagaevents.SagaBegun.spec:type_name -> internal.api.sagaspec.Saga
	8, // 1: internal.api.sagaevents.SagaBegun.config:type_name -> google.protobuf.Struct
	9, // 2: internal.api.sagaevents.StepEnded.output:type_name -> google.protobuf.Value
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_api_sagaeventspb_events_proto_init() }
func file_internal_api_sagaeventspb_events_proto_init() {
	if File_internal_api_sagaeventspb_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_api_sagaeventspb_events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SagaBegun); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SagaEnded); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*StepBegun); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*StepEnded); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*StepAborted); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*StepCompensationBegun); i {
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
		file_internal_api_sagaeventspb_events_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*StepCompensationEnded); i {
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
			RawDescriptor: file_internal_api_sagaeventspb_events_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_api_sagaeventspb_events_proto_goTypes,
		DependencyIndexes: file_internal_api_sagaeventspb_events_proto_depIdxs,
		MessageInfos:      file_internal_api_sagaeventspb_events_proto_msgTypes,
	}.Build()
	File_internal_api_sagaeventspb_events_proto = out.File
	file_internal_api_sagaeventspb_events_proto_rawDesc = nil
	file_internal_api_sagaeventspb_events_proto_goTypes = nil
	file_internal_api_sagaeventspb_events_proto_depIdxs = nil
}
