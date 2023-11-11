package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisClusterClient(t *testing.T) {
	conf := &Config{Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}}
	c := NewClusterClient(conf)
	statusCmd := c.ClusterClient.Set(context.Background(), "hello", "world", time.Second*1000)
	fmt.Println(statusCmd)
}
