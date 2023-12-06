package network

import (
	"net"
)

// DataReceivedCallback 格式如下
//type DataReceivedCallback func(sender *NetConnection, data []byte)

// DisconnectedCallback 格式如下
// type DisconnectedCallback func(sender *NetConnection)
type DisconnectedCallback func(...*NetConnection)

type NetConnection struct {
	socket net.Conn
	//dataReceivedCallback delegate.Event[DataReceivedCallback]
	//disconnectedCallback delegate.Event[DisconnectedCallback]
}

func NewNetConnection(socket net.Conn) *NetConnection {
	//dataReceivedCallback DataReceivedCallback,
	//disconnectedCallback DisconnectedCallback) *NetConnection {
	//recDelegate := delegate.NewDelegate(dataReceivedCallback, "rec_cb")
	//recEvent := delegate.Event[DataReceivedCallback]{}
	//recEvent.AddDelegate(recDelegate)
	return &NetConnection{
		socket: socket,
		//dataReceivedCallback: delegate.
	}
}
