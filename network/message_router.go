package network

import "sync"

type MessageRouter struct {
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
