package network

import (
	"encoding/binary"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/core"
	"github.com/NumberMan1/common/summer/proto_helper"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
	"net"
	"reflect"
	"slices"
	"sync"
	"time"
)

// ConnectionDataReceivedCallback 格式如下
type ConnectionDataReceivedCallback struct {
	Op func(sender Connection, data proto.Message)
}

func (d ConnectionDataReceivedCallback) Operator(args ...any) {
	d.Op(args[0].(Connection), args[1].(proto.Message))
}

// ConnectionDisconnectedCallback 格式如下
type ConnectionDisconnectedCallback struct {
	Op func(sender Connection)
}

func (d ConnectionDisconnectedCallback) Operator(args ...any) {
	d.Op(args[0].(Connection))
}

// Connection 通用网络连接，可以继承此类实现功能拓展
// 职责：发送消息，关闭连接，断开回调，接收消息回调，
type Connection interface {
	Socket() net.Conn
	Close()
	SocketSend(data []byte, offset, count int)
	Send(p proto.Message)
	SetDataReceivedCallback(dataReceivedCallback core.Func)
	SetDisconnectedCallback(disconnectedCallback core.Func)
	//sendCallBack()
}

type connection struct {
	socket net.Conn
	// 接收到数据
	dataReceivedCallback core.Event[ConnectionDataReceivedCallback]
	// 连接断开
	disconnectedCallback core.Event[ConnectionDisconnectedCallback]
	_package             *pt.Package
	mutex                sync.Mutex
}

func (c *connection) Socket() net.Conn {
	return c.socket
}

// SetDataReceivedCallback Func应为DataReceivedCallback
func (c *connection) SetDataReceivedCallback(dataReceivedCallback core.Func) {
	c.dataReceivedCallback.AddDelegate(core.NewDelegate(dataReceivedCallback.(ConnectionDataReceivedCallback), "dataReceivedCallback"))
}

// SetDisconnectedCallback Func应为DisconnectedCallback
func (c *connection) SetDisconnectedCallback(disconnectedCallback core.Func) {
	c.disconnectedCallback.AddDelegate(core.NewDelegate(disconnectedCallback.(ConnectionDisconnectedCallback), "disconnectedCallback"))
}

func NewConnection(socket net.Conn) Connection {
	c := &connection{
		socket:               socket,
		dataReceivedCallback: core.Event[ConnectionDataReceivedCallback]{},
		disconnectedCallback: core.Event[ConnectionDisconnectedCallback]{},
		mutex:                sync.Mutex{},
	}
	lfd := NewSocketReceiver(socket)
	lfd.DataReceived = c._received
	lfd.Disconnected = func() {
		if c.disconnectedCallback.HasDelegate() {
			c.disconnectedCallback.Invoke(c)
		}
	}
	lfd.Start()
	return c
}

func (c *connection) _received(data []byte) error {
	code := binary.BigEndian.Uint16(data)
	msg, err := proto_helper.ParseFrom(int(code), data, 2, len(data)-2)
	if GetMessageRouterInstance().Running {
		GetMessageRouterInstance().AddMessage(Msg{
			Sender:  c,
			Message: msg,
		})
	}
	if c.dataReceivedCallback.HasDelegate() {
		c.dataReceivedCallback.Invoke(c, msg)
	}
	return err
}

func (c *connection) Close() {
	err := c.socket.Close()
	if err != nil {
		logger.SLCDebug("NetConnection Close %s", err.Error())
		return
	}
	c.socket = nil
	if c.disconnectedCallback.HasDelegate() {
		c.disconnectedCallback.Invoke(c)
	}
}

// SocketSend 前提是data必须是大端字节序
func (c *connection) SocketSend(data []byte, offset, count int) {
	go func() {
		c.mutex.Lock()
		var err error = nil
		if c.socket != nil {
			num := 0
			for {
				err = c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					break
				}
				n, err := c.socket.Write(data[offset : offset+count])
				if err != nil {
					break
				}
				num += n
				if num == count {
					break
				}
			}
			//c.sendCallBack()
		}
		c.mutex.Unlock()
	}()
}

func (c *connection) Send(p proto.Message) {
	stream := core.AllocateDataStream()
	bs, _ := proto.Marshal(p)
	if core.IsLittleEndian() {
		slices.Reverse(bs)
	}
	code := proto_helper.SeqCode(reflect.TypeOf(p))
	stream.WriteInt32(int32(len(bs) + 2))
	stream.WriteUInt16(uint16(code))
	stream.Write(bs)
	c.SocketSend(bs, 0, len(bs))
	//buf := core.NewByteBufferByCapacity(false, 4+len(bs))
	//buf.WriteInt32(int32(len(bs)))
	//buf.WriteBytes(bs, 0, len(bs))
	//data := buf.ToArray()
	//c.SendBytes(data, 0, len(data))
}

//func (c *connection) sendCallBack() {
//
//}
