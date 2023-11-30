package network

import (
	"context"
	"fmt"
	"net"
	"testing"
)

func dCB(args ...*LengthFieldDecoder) {
	fmt.Println("关闭")
}
func rCB(args ...interface{}) {
	fmt.Printf("收到%v\n", string(args[1].([]byte)))
}

func TestSever(t *testing.T) {
	rs := []byte("你好, 客户端")
	buffer := NewByteBufferByCapacity(false, 1024)
	l := int32(len(rs))
	buffer.WriteInt32(l)
	buffer.WriteBytes(rs, 0, int(l))
	listener, _ := net.Listen("tcp4", "127.0.0.1:20000")
	conn, err := listener.Accept()
	if err != nil {
		println(err)
	}
	println(conn.RemoteAddr().String())
	lengthFieldDecoder := NewLengthFieldDecoderDefault(conn, 0, 4)
	lengthFieldDecoder.AddDataReceivedCB(rCB, "r")
	lengthFieldDecoder.AddDisconnectCB(dCB, "d")
	lengthFieldDecoder.Start(context.Background())

	_, err = conn.Write(buffer.ToArray())
	if err != nil {
		println(err)
	}
	select {}
}
