package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	mq, err := ConnectRabbitMQ("num", "123", "localhost:5672", "customers")
	if err != nil {
		println(err.Error())
		return
	}
	client, err := NewRMQClient(mq)
	if err != nil {
		println(err.Error())
		return
	}
	defer client.Close()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	for i := 1; i != 5; i += 1 {
		err = client.Send(timeout, "customer_events", "customer.*", amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         []byte("你好, 这是Transient"),
		})
		if err != nil {
			println(err.Error())
			return
		}
		err = client.Send(timeout, "customer_events", "customer_queue.*", amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte("你好, 这是Persistent"),
		})
		if err != nil {
			println(err.Error())
			return
		}
	}
	consume, err := client.Consume("customer_queue_test2", "producer", false)
	if err != nil {
		println(err.Error())
		return
	}
	group, _ := errgroup.WithContext(timeout)
	group.SetLimit(10)
	go func() {
		for message := range consume {
			msg := message
			group.Go(func() error {
				if !msg.Redelivered {
					fmt.Printf("new message %v %v第一次\n", msg.MessageId, string(msg.Body))
					if err := msg.Nack(false, true); err != nil {
						println(err.Error())
						return err
					}
				} else {
					fmt.Printf("new message %v %v\n第2次", msg.MessageId, string(msg.Body))
					if err := msg.Ack(false); err != nil {
						println(err.Error())
						return err
					}
				}
				return nil
			})
		}
	}()
	select {}
}

func TestConsume(t *testing.T) {
	mq, err := ConnectRabbitMQ("num", "123", "localhost:5672", "customers")
	if err != nil {
		println(err.Error())
		return
	}
	client, err := NewRMQClient(mq)
	defer client.Close()
	if err != nil {
		println(err.Error())
		return
	}
	err = client.ApplyQos(10, 0, false)
	if err != nil {
		println(err.Error())
		return
	}
	queue, err := client.CreateQueue("customer_queue_test2", true, false)
	if err != nil {
		println(err.Error())
		return
	}
	err = client.CreateBind("customer_queue_test2", "customer.*", "customer_events")
	if err != nil {
		println(err.Error())
		return
	}
	consume, err := client.Consume(queue.Name, "customer", false)
	if err != nil {
		println(err.Error())
		return
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	group, ctx := errgroup.WithContext(timeout)
	group.SetLimit(10)
	go func() {
		for message := range consume {
			msg := message
			group.Go(func() error {
				fmt.Printf("new message %v %v\n", msg.MessageId, string(msg.Body))
				if err := msg.Ack(false); err != nil {
					println(err)
					return err
				}
				if err := client.Send(ctx, "customer_events", "customer.*", amqp.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp.Persistent,
					Body:          []byte("收到消息"),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					println(err)
					return err
				}
				return nil
			})
		}

	}()
	select {}
}
