package proto_helper

// ProtobufTool Protobuf序列化与反序列化

import (
	"github.com/NumberMan1/common/logger"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"sort"
)

var (
	registry = make(map[string]reflect.Type)
	dict1    = make(map[int]reflect.Type)
	dict2    = make(map[reflect.Type]int)
)

func init() {
	ts := protoregistry.GlobalTypes
	ts.RangeMessages(func(messageType protoreflect.MessageType) bool {
		name := messageType.Descriptor().FullName()
		if messageType.Descriptor().FullName() != "google.protobuf.Any" { // 排除该类型, 并不直接作为通讯类型
			registry[string(name)] = reflect.TypeOf(messageType.New().Interface())
		}
		return true
	})
	list := make([]string, 0)
	for k := range registry {
		list = append(list, k)
	}
	sort.Slice(list, func(i, j int) bool {
		if len(list[i]) != len(list[j]) {
			return len(list[i]) < len(list[j])
		}
		return list[i] < list[j]
	})
	for i, fullName := range list {
		t := registry[fullName]
		logger.SLCDebug("Proto类型注册：%d - %s", i, fullName)
		dict1[i] = t
		dict2[t] = i
	}
}

// Parse 依据给定的FullName解析数据
func Parse(fullName []byte, data []byte) (msg proto.Message, err error) {
	messageName := protoreflect.FullName(fullName)
	pbType, err := protoregistry.GlobalTypes.FindMessageByName(messageName)
	if err != nil {
		return
	}
	msg = pbType.New().Interface()
	err = proto.Unmarshal(data, msg)
	return
}

func SeqCode(t reflect.Type) int {
	return dict2[t]
}

func SeqType(code int) reflect.Type {
	return dict1[code]
}

func ParseFrom(typeCode int, data []byte, offset, len int) (msg proto.Message, err error) {
	seqType := SeqType(typeCode)
	messageName := protoreflect.FullName(seqType.String()[1:]) // 去除*,防止2重指针
	pbType, err := protoregistry.GlobalTypes.FindMessageByName(messageName)
	if err != nil {
		return
	}
	msg = pbType.New().Interface()
	err = proto.Unmarshal(data[offset:offset+len], msg)
	logger.SLCInfo("解析消息：code=%s - %v", seqType.String(), msg)
	return
}
