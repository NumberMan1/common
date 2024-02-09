package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisClusterClient(t *testing.T) {
	//conf := &Config{Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}}
	//c := NewClusterClient(conf)
	conf := &Config{Addrs: []string{":12000"}}
	c := NewNormalClient(conf)
	//statusCmd := c.ClusterClient.Set(context.Background(), "hello", "world", timeunit.Second*1000)
	statusCmd := c.Client.Set(context.Background(), "hello", "world", time.Second*1000)
	fmt.Println(statusCmd)
}
