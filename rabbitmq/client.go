package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQClient struct {
	// 客户端使用的连接
	conn *amqp.Connection
	// 用于处理/发送消息
	ch *amqp.Channel
}

// ConnectRabbitMQ connection用于复用, 但是只有同一类型的才能复用, 发布用一个, 消费一个, 不能发布和消费同时用
func ConnectRabbitMQ(user, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", user, password, host, vhost))
}

func NewRMQClient(connection *amqp.Connection) (RMQClient, error) {
	ch, err := connection.Channel()
	if err != nil {
		return RMQClient{}, err
	}
	if err = ch.Confirm(false); err != nil {
		return RMQClient{}, err
	}
	return RMQClient{
		conn: connection,
		ch:   ch,
	}, nil
}

// Close 不关闭conn防止其它客户端需要
func (rc RMQClient) Close() error {
	return rc.ch.Close()
}

// CreateQueue 创建新的队列
func (rc RMQClient) CreateQueue(name string, durable, autoDelete bool) (amqp.Queue, error) {
	queue, err := rc.ch.QueueDeclare(name, durable, autoDelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, nil
	}
	return queue, nil
}

// CreateBind 通过指定的交换机和路由key来绑定队列
func (rc RMQClient) CreateBind(name, key, exchange string) error {
	return rc.ch.QueueBind(name, key, exchange, false, nil)
}

func (rc RMQClient) Send(ctx context.Context, exchange, key string, option amqp.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		key,
		// 表示如果发生错误, 返回error
		true,
		// 被弃用, 应该为false
		false,
		option,
	)
	if err != nil {
		return err
	}
	confirmation.Wait()
	return nil
}

func (rc RMQClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(
		queue, consumer, autoAck,
		// 该队列是否为该消费者一人独享
		false,
		// 用于分布式只有一个服务器消费, 实际不支持
		false,
		false,
		nil,
	)
}

// ApplyQos
// count允许未确认的消息的数量
// size允许多少byte
// global是否为全局设置
func (rc RMQClient) ApplyQos(count, size int, global bool) error {
	return rc.ch.Qos(count, size, global)
}
