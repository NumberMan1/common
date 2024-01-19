package test

import (
	"context"
	"github.com/NumberMan1/common/summer/core"
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
	r := core.DataReceivedEventHandler{
		Op: rCB,
	}
	d := core.DisconnectedEventHandler{
		Op: dCB,
	}
	rs := []byte("你好, 服务器")
	buffer := core.NewByteBufferByCapacity(false, 1024)
	l := int32(len(rs))
	buffer.WriteInt32(l)
	buffer.WriteBytes(rs, 0, int(l))
	conn, err := net.Dial("tcp4", "127.0.0.1:20000")
	if err != nil {
		println(err)
	}
	println(conn.RemoteAddr().String())
	lengthFieldDecoder := core.NewLengthFieldDecoderDefault(conn, 0, 4)
	lengthFieldDecoder.AddDataReceivedCB(r, "r")
	lengthFieldDecoder.AddDisconnectCB(d, "d")
	lengthFieldDecoder.Start(context.Background())
	conn.Write(buffer.ToArray())
	//rs = []byte("你好, 服务器 第2次")
	//buffer = core.NewByteBufferByCapacity(false, 1024)
	//l = int32(len(rs))
	//buffer.WriteInt32(l)
	//buffer.WriteBytes(rs, 0, int(l))
	//_, err = conn.Write(buffer.ToArray())
	//if err != nil {
	//	println(err)
	//	return
	//}
	select {}
}
