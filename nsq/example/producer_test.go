package example

import (
	"encoding/json"
	gnsq "github.com/nsqio/go-nsq"
	"github.com/phuhao00/broker/nsq"
	"testing"
)

func TestProducer(t *testing.T) {
	client := nsq.NewProducerClient("127.0.0.1:4150")
	go client.Run()
	msgIDGood := gnsq.MessageID{'1'}
	msgGood := gnsq.NewMessage(msgIDGood, []byte("good"))
	bytes, _ := json.Marshal(msgGood)
	client.Pub("test", bytes)
	select {}
}
