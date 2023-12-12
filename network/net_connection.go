package network

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common/delegate"
	pt "github.com/NumberMan1/common/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
	"net"
	"slices"
	"sync"
	"time"
)

// DataReceivedCallback 格式如下
type DataReceivedCallback struct {
	Op func(sender *NetConnection, data []byte)
}

func (d DataReceivedCallback) Operator(args ...any) {
	d.Op(args[0].(*NetConnection), args[1].([]byte))
}

// DisconnectedCallback 格式如下
type DisconnectedCallback struct {
	Op func(sender *NetConnection)
}

func (d DisconnectedCallback) Operator(args ...any) {
	d.Op(args[0].(*NetConnection))
}

type NetConnection struct {
	socket               net.Conn
	dataReceivedCallback delegate.Event[DataReceivedCallback]
	disconnectedCallback delegate.Event[DisconnectedCallback]
	decoder              *LengthFieldDecoder
	_package             *pt.Package
	mutex                sync.Mutex
}

func NewNetConnection(socket net.Conn, dataReceivedCallback DataReceivedCallback, disconnectedCallback DisconnectedCallback) *NetConnection {
	recDelegate := delegate.NewDelegate(dataReceivedCallback, "rec_cb")
	recEvent := delegate.Event[DataReceivedCallback]{}
	recEvent.AddDelegate(recDelegate)
	disDelegate := delegate.NewDelegate(disconnectedCallback, "dis_cb")
	disEvent := delegate.Event[DisconnectedCallback]{}
	disEvent.AddDelegate(disDelegate)
	c := &NetConnection{
		socket:               socket,
		dataReceivedCallback: recEvent,
		disconnectedCallback: disEvent,
		decoder:              NewLengthFieldDecoder(socket, 0, 4, 0, 4, 64*1024),
		mutex:                sync.Mutex{},
	}
	c.decoder.AddDataReceivedCB(DataReceivedEventHandler{Op: c.onDataReceived}, "net_con_rec")
	c.decoder.AddDisconnectCB(DisconnectedEventHandler{Op: func(decoder *LengthFieldDecoder) {
		c.disconnectedCallback.Invoke(c)
	}}, "net_con_dis")
	c.decoder.Start(context.Background())
	return c
}

func (c *NetConnection) onDataReceived(sender *LengthFieldDecoder, buffer []byte) {
	//fmt.Printf("buf:%v", len(buffer))
	c.dataReceivedCallback.Invoke(c, buffer)
}

func (c *NetConnection) Close() {
	err := c.socket.Close()
	if err != nil {
		fmt.Printf("NetConnection Close %v\n", err.Error())
		return
	}
	c.socket = nil
	if c.disconnectedCallback.HasDelegate() {
		c.disconnectedCallback.Invoke(c)
	}
}

func (c *NetConnection) GetRequest() *pt.Request {
	if c._package == nil {
		c._package = &pt.Package{}
	}
	if c._package.Request == nil {
		c._package.Request = &pt.Request{}
	}
	return c._package.Request
}

func (c *NetConnection) GetResponse() *pt.Response {
	if c._package == nil {
		c._package = &pt.Package{}
	}
	if c._package.Response == nil {
		c._package.Response = &pt.Response{}
	}
	return c._package.Response
}

func (c *NetConnection) SendBytes(data []byte, offset, count int) error {
	mutex.Lock()
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
		c.sendCallBack()
	}
	mutex.Unlock()
	return err
}

func (c *NetConnection) SendPackage(p *pt.Package) error {
	bs, _ := proto.Marshal(p)
	if IsLittleEndian() {
		slices.Reverse(bs)
	}
	buf := NewByteBufferByCapacity(false, 4+len(bs))
	buf.WriteInt32(int32(len(bs)))
	buf.WriteBytes(bs, 0, len(bs))
	data := buf.ToArray()
	err := c.SendBytes(data, 0, len(data))
	return err
}

func (c *NetConnection) Send() error {
	var err error = nil
	if c._package != nil {
		err = c.SendPackage(c._package)
	}
	c._package = nil
	return err
}

func (c *NetConnection) sendCallBack() {

}
