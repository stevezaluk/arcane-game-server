// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.6.1
// source: state.proto

package models

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

type GameState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	GameMode    string             `protobuf:"bytes,2,opt,name=gameMode,proto3" json:"gameMode,omitempty"`
	Players     map[string]*Player `protobuf:"bytes,3,rep,name=players,proto3" json:"players,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Battlefield *Zone              `protobuf:"bytes,4,opt,name=battlefield,proto3" json:"battlefield,omitempty"`
	Exile       *Zone              `protobuf:"bytes,5,opt,name=exile,proto3" json:"exile,omitempty"`
	Command     *Zone              `protobuf:"bytes,6,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *GameState) Reset() {
	*x = GameState{}
	mi := &file_state_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameState) ProtoMessage() {}

func (x *GameState) ProtoReflect() protoreflect.Message {
	mi := &file_state_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameState.ProtoReflect.Descriptor instead.
func (*GameState) Descriptor() ([]byte, []int) {
	return file_state_proto_rawDescGZIP(), []int{0}
}

func (x *GameState) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GameState) GetGameMode() string {
	if x != nil {
		return x.GameMode
	}
	return ""
}

func (x *GameState) GetPlayers() map[string]*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *GameState) GetBattlefield() *Zone {
	if x != nil {
		return x.Battlefield
	}
	return nil
}

func (x *GameState) GetExile() *Zone {
	if x != nil {
		return x.Exile
	}
	return nil
}

func (x *GameState) GetCommand() *Zone {
	if x != nil {
		return x.Command
	}
	return nil
}

var File_state_proto protoreflect.FileDescriptor

var file_state_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67,
	0x61, 0x6d, 0x65, 0x1a, 0x0c, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0a, 0x7a, 0x6f, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb3, 0x02,
	0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x61, 0x6d, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x73, 0x12, 0x2c, 0x0a, 0x0b, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e,
	0x5a, 0x6f, 0x6e, 0x65, 0x52, 0x0b, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x20, 0x0a, 0x05, 0x65, 0x78, 0x69, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x52, 0x05, 0x65, 0x78,
	0x69, 0x6c, 0x65, 0x12, 0x24, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x5a, 0x6f, 0x6e, 0x65,
	0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x48, 0x0a, 0x0c, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x22, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x74, 0x65, 0x76, 0x65, 0x7a, 0x61, 0x6c, 0x75, 0x6b, 0x2f, 0x61, 0x72, 0x63,
	0x61, 0x6e, 0x65, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_state_proto_rawDescOnce sync.Once
	file_state_proto_rawDescData = file_state_proto_rawDesc
)

func file_state_proto_rawDescGZIP() []byte {
	file_state_proto_rawDescOnce.Do(func() {
		file_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_state_proto_rawDescData)
	})
	return file_state_proto_rawDescData
}

var file_state_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_state_proto_goTypes = []any{
	(*GameState)(nil), // 0: game.GameState
	nil,               // 1: game.GameState.PlayersEntry
	(*Zone)(nil),      // 2: game.Zone
	(*Player)(nil),    // 3: game.Player
}
var file_state_proto_depIdxs = []int32{
	1, // 0: game.GameState.players:type_name -> game.GameState.PlayersEntry
	2, // 1: game.GameState.battlefield:type_name -> game.Zone
	2, // 2: game.GameState.exile:type_name -> game.Zone
	2, // 3: game.GameState.command:type_name -> game.Zone
	3, // 4: game.GameState.PlayersEntry.value:type_name -> game.Player
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_state_proto_init() }
func file_state_proto_init() {
	if File_state_proto != nil {
		return
	}
	file_player_proto_init()
	file_zone_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_state_proto_goTypes,
		DependencyIndexes: file_state_proto_depIdxs,
		MessageInfos:      file_state_proto_msgTypes,
	}.Build()
	File_state_proto = out.File
	file_state_proto_rawDesc = nil
	file_state_proto_goTypes = nil
	file_state_proto_depIdxs = nil
}
