package main

import (
	"github.com/NumberMan1/common/game_server/service"
	"github.com/NumberMan1/common/logger"
	n2 "github.com/NumberMan1/common/summer/network"
)

func f(msg n2.Msg) bool {
	return true
}

func OnUserLoginRequest(msg n2.Msg) {
	//println(p)
	//fmt.Printf("发现用户登录请求:%v %v\n", msg.UserLogin.Username, msg.Request.UserLogin.Password)
}

func initServices() {
	netService := service.NewNetService()
	netService.Start()
	logger.SLCDebug("网络服务启动完成")
	service.GetSpaceServiceInstance().Start()
	logger.SLCDebug("地图服务启动完成")
	service.GetUserServiceInstance().Start()
	logger.SLCDebug("玩家服务启动完成")
}

func main() {
	initServices()
	//n2.GetMessageRouterInstance().Subscribe("proto.UserLoginRequest", n2.MessageHandler{Op: OnUserLoginRequest})
	select {}
}
