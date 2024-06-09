package network

import (
	"net"
	"sync/atomic"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"google.golang.org/protobuf/proto"
)

type TcpServerEventHandler struct {
	Op func(tsl *TcpServer, socket *net.TCPConn)
}

func (s TcpServerEventHandler) Operator(args ...any) {
	s.Op(args[0].(*TcpServer), args[1].(*net.TCPConn))
}

type TcpServerConnectedCallback struct {
	Op func(connection Connection)
}

func (c TcpServerConnectedCallback) Operator(args ...any) {
	c.Op(args[0].(Connection))
}

type TcpServerDataReceivedCallback = ConnectionDataReceivedCallback

type TcpServerDisconnectedCallback = ConnectionDisconnectedCallback

// TcpServer 负责监听TCP网络端口，异步接收Socket连接
// 负责监听TCP网络端口，异步接收Socket连接
// -- Connected        有新的连接
// -- DataReceived     有新的消息
// -- Disconnected     有连接断开
// Start()         启动服务器
// Stop()          关闭服务器
// IsRunning     是否正在运行
type TcpServer struct {
	endAddr        *net.TCPAddr
	serverListener *net.TCPListener //服务端监听对象
	//客户端接入事件
	socketConnected ns.Event[TcpServerEventHandler]
	//事件委托：新的连接
	connected ns.Event[TcpServerConnectedCallback]
	//事件委托：收到消息
	dataReceived ns.Event[TcpServerDataReceivedCallback]
	//事件委托：连接断开
	disconnected ns.Event[TcpServerDisconnectedCallback]
	//可以排队接受的传入连接数
	maxBacklog int
	//当前的传入连接数
	curBacklog  atomic.Int32
	acceptEvent *ns.AutoResetEvent
	chStop      chan struct{} // 用于发送停止信号
}

func (ts *TcpServer) SetSocketConnected(socketConnectedCB TcpServerEventHandler) {
	ts.socketConnected.AddDelegate(ns.NewDelegate(socketConnectedCB, "socketConnected"))
}

func (ts *TcpServer) SetConnectedCallback(connectedCallback TcpServerConnectedCallback) {
	ts.connected.AddDelegate(ns.NewDelegate(connectedCallback, "connectedCallback"))
}

func (ts *TcpServer) SetDataReceivedCallback(dataReceivedCallback TcpServerDataReceivedCallback) {
	ts.dataReceived.AddDelegate(ns.NewDelegate(dataReceivedCallback, "dataReceivedCallback"))
}

func (ts *TcpServer) SetDisconnectedCallback(disconnectedCallback TcpServerDisconnectedCallback) {
	ts.disconnected.AddDelegate(ns.NewDelegate(disconnectedCallback, "disconnectedCallback"))
}

func NewTcpServer(address string) (*TcpServer, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	return &TcpServer{
		endAddr:         addr,
		serverListener:  nil,
		socketConnected: ns.Event[TcpServerEventHandler]{},
		connected:       ns.Event[TcpServerConnectedCallback]{},
		dataReceived:    ns.Event[TcpServerDataReceivedCallback]{},
		disconnected:    ns.Event[TcpServerDisconnectedCallback]{},
		maxBacklog:      100,
		curBacklog:      atomic.Int32{},
		acceptEvent:     ns.NewAutoResetEvent(),
		chStop:          make(chan struct{}, 1),
	}, nil
}

func (ts *TcpServer) IsRunning() bool {
	return ts.serverListener != nil
}

func (ts *TcpServer) Start() {
	if !ts.IsRunning() {
		logger.SLCDebug("start listen %v", ts.endAddr.String())
		ts.serverListener, _ = net.ListenTCP("tcp4", ts.endAddr)
		go func() {
			accept, err := ts.serverListener.AcceptTCP()
			ts.onAccept(accept, err)
		}()
	}
}

func (ts *TcpServer) Stop() error {
	ts.chStop <- struct{}{}
	if ts.serverListener == nil {
		return nil
	}
	err := ts.serverListener.Close()
	ts.serverListener = nil
	return err
}

func (ts *TcpServer) onAccept(socket *net.TCPConn, err error) {
	for {
		select {
		case <-ts.chStop:
			return
		default:
			if err == nil {
				if socket != nil {
					ts.curBacklog.Add(1)
					ts.OnSocketConnected(socket)
					ts.socketConnected.Invoke(ts, socket)
				}
			}
			if int(ts.curBacklog.Load()) > ts.maxBacklog {
				ts.acceptEvent.Wait()
			}
			socket, err = ts.serverListener.AcceptTCP()
		}
	}
}

func (ts *TcpServer) OnSocketConnected(socket *net.TCPConn) {
	if ts.socketConnected.HasDelegate() {
		ts.socketConnected.Invoke(ts, socket)
	}
	connection := NewConnection(socket)
	receivedCallback := TcpServerDataReceivedCallback{Op: func(sender Connection, data proto.Message) {
		if ts.dataReceived.HasDelegate() {
			ts.dataReceived.Invoke(sender, data)
		}
	}}
	disconnectedCallback := TcpServerDisconnectedCallback{Op: func(sender Connection) {
		ts.curBacklog.Add(-1)
		if int(ts.curBacklog.Load()) > ts.maxBacklog { // 如果数量超过则代表正在等待唤醒
			ts.acceptEvent.Set()
		}
		if ts.disconnected.HasDelegate() {
			ts.disconnected.Invoke(sender)
		}
	}}
	connection.SetDataReceivedCallback(receivedCallback)
	connection.SetDisconnectedCallback(disconnectedCallback)
	if ts.connected.HasDelegate() {
		ts.connected.Invoke(connection)
	}
}
