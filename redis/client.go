package redis

import (
	"github.com/NumberMan1/common"
	"github.com/redis/go-redis/v9"
)

type ClusterClient struct {
	*common.BaseComponent
	*redis.ClusterClient
}

func NewClusterClient(conf *Config) *ClusterClient {
	c := &ClusterClient{
		BaseComponent: common.NewBaseComponent(),
		ClusterClient: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: conf.Addrs,
			//Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
			// To route commands by latency or randomly, enable one of the following.
			//RouteByLatency: true,
			//RouteRandomly: true,
		}),
	}
	return c
}
