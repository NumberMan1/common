package network

import (
	"container/list"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Msg struct {
	Sender  *NetConnection
	Message proto.Message
}

type MessageRouter struct {
	ThreadCount int  //工作协程数
	WorkerCount int  //正在工作的协程数
	Running     bool //是否正在运行状态
	//threadEvent  delegate.Event[]
	messageQueue list.List
}

var (
	instance *MessageRouter
	once     sync.Once
)

func GetMessageRouterInstance() *MessageRouter {
	once.Do(func() {
		instance = &MessageRouter{}
	})
	return instance
}
