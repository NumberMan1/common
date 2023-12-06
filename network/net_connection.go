package network

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common/delegate"
	"net"
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
	}
	c.decoder.AddDataReceivedCB(DataReceivedEventHandler{Op: c.OnDataReceived}, "net_con_rec")
	c.decoder.AddDisconnectCB(DisconnectedEventHandler{Op: func(decoder *LengthFieldDecoder) {
		c.disconnectedCallback.Invoke(c)
	}}, "net_con_dis")
	c.decoder.Start(context.Background())
	return c
}

func (c NetConnection) OnDataReceived(*LengthFieldDecoder, []byte) {

}

func (c *NetConnection) Close() {
	err := c.socket.Close()
	if err != nil {
		fmt.Printf("NetConnection Close %v", err.Error())
		return
	}
	c.socket = nil
	if c.disconnectedCallback.HasDelegate() {
		c.disconnectedCallback.Invoke(c)
	}
}
