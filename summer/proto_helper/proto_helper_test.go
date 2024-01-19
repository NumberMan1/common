package proto_helper

import (
	"fmt"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	request := &pt.UserLoginRequest{Username: "123"}
	name := request.ProtoReflect().Descriptor().FullName()
	bytes, err := proto.Marshal(request)
	fmt.Println(name.Name())
	m, err := Parse([]byte(name), bytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(m)
}

func TestParseFrom(t *testing.T) {
	request := &pt.UserLoginRequest{Username: "123"}
	bytes, _ := proto.Marshal(request)
	code := SeqCode(reflect.TypeOf(request))
	fmt.Println(code)
	msg, _ := ParseFrom(code, bytes, 0, len(bytes))
	fmt.Println(msg.(*pt.UserLoginRequest).Username)
}
