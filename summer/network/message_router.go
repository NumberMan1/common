package network

import (
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network/message_router"
)

var (
	singleMessageRouter = singleton.Singleton{}
)

// GetMessageRouterInstance 获取消息路由单例,默认使用string为订阅类型
func GetMessageRouterInstance() *message_router.NameRouter {
	instance, _ := singleton.GetOrDo[*message_router.NameRouter](&singleMessageRouter, func() (*message_router.NameRouter, error) {
		return message_router.NewNameRouter(), nil
	})
	return instance
}
