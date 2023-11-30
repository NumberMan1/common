package network

import (
	"context"
	"net"
	"testing"
)

//func dCB(args ...*LengthFieldDecoder) {
//	fmt.Printf("%v关闭", args[0])
//}
//func rCB(args ...interface{}) {
//	fmt.Printf("收到%v", args[1].([]byte))
//}

func TestClient(t *testing.T) {
	rs := []byte("你好, 服务器")
	buffer := NewByteBufferByCapacity(false, 1024)
	l := int32(len(rs))
	buffer.WriteInt32(l)
	buffer.WriteBytes(rs, 0, int(l))
	conn, err := net.Dial("tcp4", "127.0.0.1:20000")
	if err != nil {
		println(err)
	}
	println(conn.RemoteAddr().String())
	lengthFieldDecoder := NewLengthFieldDecoderDefault(conn, 0, 4)
	lengthFieldDecoder.AddDataReceivedCB(rCB, "r")
	lengthFieldDecoder.AddDisconnectCB(dCB, "d")
	lengthFieldDecoder.Start(context.Background())
	conn.Write(buffer.ToArray())
	select {}
}
