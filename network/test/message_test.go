package test

import (
	pt "github.com/NumberMan1/common/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func Pr(message proto.Message) {
	s := reflect.TypeOf(message)
	println(s.String())
}

func TestMsg(t *testing.T) {
	Pr(&pt.Package{
		Id:       0,
		Request:  nil,
		Response: nil,
	})
	Pr(&pt.Response{
		UserRegister: nil,
		UserLogin:    nil,
	})
	Pr(&pt.UserLoginRequest{
		Username: "",
		Password: "",
	})
}
