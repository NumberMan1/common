package main

import (
	n2 "github.com/NumberMan1/common/summer/network"
	service2 "github.com/NumberMan1/common/summer/network/game_server/service"
)

func f(msg n2.Msg) bool {
	return true
}

func OnUserLoginRequest(msg n2.Msg) {
	//println(p)
	//fmt.Printf("发现用户登录请求:%v %v\n", msg.UserLogin.Username, msg.Request.UserLogin.Password)
}

func main() {
	n2.GetMessageRouterInstance()
	service := service2.NewNetService()
	service.Start()
	n2.GetMessageRouterInstance().Subscribe("proto.UserLoginRequest", n2.MessageHandler{Op: OnUserLoginRequest})
	select {}
}
