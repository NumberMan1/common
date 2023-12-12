package network

import (
	"fmt"
	"github.com/NumberMan1/common/delegate"
	"net"
)

type SocketConnectedCB struct {
	Op func(tsl *TcpSocketListener, socket net.Conn)
}

func (s SocketConnectedCB) Operator(args ...any) {
	s.Op(args[0].(*TcpSocketListener), args[1].(net.Conn))
}

// TcpSocketListener 负责监听TCP网络端口，异步接收Socket连接
type TcpSocketListener struct {
	endAddr         *net.TCPAddr
	serverListener  *net.TCPListener //服务端监听对象
	socketConnected delegate.Event[SocketConnectedCB]
	chStop          chan struct{} // 用于发送停止信号
}

func (tsl *TcpSocketListener) SetSocketConnected(socketConnectedCB SocketConnectedCB) {
	tsl.socketConnected.AddDelegate(delegate.NewDelegate(socketConnectedCB, "socketConnected"))
}

func NewTcpSocketListener(address string) (*TcpSocketListener, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	return &TcpSocketListener{
		endAddr:         addr,
		serverListener:  nil,
		socketConnected: delegate.Event[SocketConnectedCB]{},
		chStop:          make(chan struct{}, 1),
	}, nil
}

func (tsl *TcpSocketListener) IsRunning() bool {
	return tsl.serverListener != nil
}

func (tsl *TcpSocketListener) Start() {
	if !tsl.IsRunning() {
		fmt.Printf("start listen %v\n", tsl.endAddr.String())
		tsl.serverListener, _ = net.ListenTCP("tcp4", tsl.endAddr)
		accept, err := tsl.serverListener.Accept()
		go tsl.onAccept(accept, err)
	}
}

func (tsl *TcpSocketListener) Stop() error {
	if tsl.serverListener == nil {
		return nil
	}
	err := tsl.serverListener.Close()
	tsl.serverListener = nil
	return err
}

func (tsl *TcpSocketListener) onAccept(socket net.Conn, err error) {
	select {
	case <-tsl.chStop:
		return
	default:
		if err == nil {
			if socket != nil {
				tsl.socketConnected.Invoke(tsl, socket)
			}
		}
		accept, err := tsl.serverListener.Accept()
		tsl.onAccept(accept, err)
	}
}
