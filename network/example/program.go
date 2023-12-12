package main

import (
	"fmt"
	"github.com/NumberMan1/common/network"
	pt "github.com/NumberMan1/common/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
)

func f(msg network.Msg) bool {
	return true
}

func OnUserLoginRequest(sender *network.NetConnection, p proto.Message) {
	//println(p)
	msg := p.(*pt.Package)
	fmt.Printf("发现用户登录请求:%v %v\n", msg.Request.UserLogin.Username, msg.Request.UserLogin.Password)
}

func main() {
	network.GetMessageRouterInstance().SetMsgOKHandler(f)
	service := network.NewNetService()
	err := service.Init(32510)
	if err != nil {
		fmt.Println(err)
		return
	}
	service.Start(2)
	network.GetMessageRouterInstance().On(network.MessageHandler[proto.Message]{ID: int32(pt.PackageID_userLoginRequest), Op: OnUserLoginRequest})
	select {}
}
