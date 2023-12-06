package test

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common/network"
	"net"
	"testing"
)

func dCB(*network.LengthFieldDecoder) {
	fmt.Println("关闭")
}
func rCB(sender *network.LengthFieldDecoder, data []byte) {
	fmt.Printf("收到%v\n", data)
	buffer := network.NewByteBufferByBuf(false, data)
	fmt.Printf("收到%v\n", buffer.ReadString())
}

func TestSever(t *testing.T) {
	r := network.DataReceivedEventHandler{
		Op: rCB,
	}
	d := network.DisconnectedEventHandler{
		Op: dCB,
	}
	rs := []byte("你好, 客户端")
	buffer := network.NewByteBufferByCapacity(false, 1024)
	l := int32(len(rs))
	buffer.WriteInt32(l)
	buffer.WriteBytes(rs, 0, int(l))
	listener, _ := net.Listen("tcp4", "127.0.0.1:20000")
	conn, err := listener.Accept()
	if err != nil {
		println(err)
	}
	println(conn.RemoteAddr().String())
	lengthFieldDecoder := network.NewLengthFieldDecoderDefault(conn, 0, 4)
	lengthFieldDecoder.AddDataReceivedCB(r, "r")
	lengthFieldDecoder.AddDisconnectCB(d, "d")
	lengthFieldDecoder.Start(context.Background())

	_, err = conn.Write(buffer.ToArray())
	if err != nil {
		println(err)
	}
	select {}
}
