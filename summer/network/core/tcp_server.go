package core

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/core"
	"github.com/NumberMan1/common/summer/network"
	"google.golang.org/protobuf/proto"
	"net"
)

type TcpServerEventHandler struct {
	Op func(tsl *TcpServer, socket *net.TCPConn)
}

func (s TcpServerEventHandler) Operator(args ...any) {
	s.Op(args[0].(*TcpServer), args[1].(*net.TCPConn))
}

type TcpServerConnectedCallback struct {
	Op func(connection network.Connection)
}

func (c TcpServerConnectedCallback) Operator(args ...any) {
	c.Op(args[0].(network.Connection))
}

type TcpServerDataReceivedCallback = network.ConnectionDataReceivedCallback

type TcpServerDisconnectedCallback = network.ConnectionDisconnectedCallback

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
	socketConnected core.Event[TcpServerEventHandler]
	//事件委托：新的连接
	connected core.Event[TcpServerConnectedCallback]
	//事件委托：收到消息
	dataReceived core.Event[TcpServerDataReceivedCallback]
	//事件委托：连接断开
	disconnected core.Event[TcpServerDisconnectedCallback]
	//可以排队接受的传入连接数
	maxBacklog int
	//当前的传入连接数
	curBacklog  int
	acceptEvent *core.AutoResetEvent
	chStop      chan struct{} // 用于发送停止信号
}

func (ts *TcpServer) SetSocketConnected(socketConnectedCB TcpServerEventHandler) {
	ts.socketConnected.AddDelegate(core.NewDelegate(socketConnectedCB, "socketConnected"))
}

func (ts *TcpServer) SetConnectedCallback(connectedCallback TcpServerConnectedCallback) {
	ts.connected.AddDelegate(core.NewDelegate(connectedCallback, "connectedCallback"))
}

func (ts *TcpServer) SetDataReceivedCallback(dataReceivedCallback TcpServerDataReceivedCallback) {
	ts.dataReceived.AddDelegate(core.NewDelegate(dataReceivedCallback, "dataReceivedCallback"))
}

func (ts *TcpServer) SetDisconnectedCallback(disconnectedCallback TcpServerDisconnectedCallback) {
	ts.disconnected.AddDelegate(core.NewDelegate(disconnectedCallback, "disconnectedCallback"))
}

func NewTcpServer(address string) (*TcpServer, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	return &TcpServer{
		endAddr:         addr,
		serverListener:  nil,
		socketConnected: core.Event[TcpServerEventHandler]{},
		connected:       core.Event[TcpServerConnectedCallback]{},
		dataReceived:    core.Event[TcpServerDataReceivedCallback]{},
		disconnected:    core.Event[TcpServerDisconnectedCallback]{},
		maxBacklog:      100,
		curBacklog:      0,
		acceptEvent:     core.NewAutoResetEvent(),
		chStop:          make(chan struct{}, 1),
	}, nil
}

func (ts *TcpServer) IsRunning() bool {
	return ts.serverListener != nil
}

func (ts *TcpServer) Start() {
	if !ts.IsRunning() {
		logger.SLCDebug("start listen %v\n", ts.endAddr.String())
		ts.serverListener, _ = net.ListenTCP("tcp4", ts.endAddr)
		go func() {
			accept, err := ts.serverListener.AcceptTCP()
			ts.onAccept(accept, err)
		}()
	}
}

func (ts *TcpServer) Stop() error {
	if ts.serverListener == nil {
		return nil
	}
	err := ts.serverListener.Close()
	ts.serverListener = nil
	return err
}

func (ts *TcpServer) onAccept(socket *net.TCPConn, err error) {
	select {
	case <-ts.chStop:
		return
	default:
		if err == nil {
			if socket != nil {
				ts.curBacklog += 1
				ts.OnSocketConnected(socket)
				ts.socketConnected.Invoke(ts, socket)
			}
		}
		if ts.curBacklog < ts.maxBacklog {
			ts.acceptEvent.Wait()
		}
		accept, err := ts.serverListener.AcceptTCP()
		ts.onAccept(accept, err)
	}
}

func (ts *TcpServer) OnSocketConnected(socket *net.TCPConn) {
	if ts.socketConnected.HasDelegate() {
		ts.socketConnected.Invoke(ts, socket)
	}
	connection := network.NewConnection(socket)
	receivedCallback := TcpServerDataReceivedCallback{Op: func(sender network.Connection, data proto.Message) {
		if ts.dataReceived.HasDelegate() {
			ts.dataReceived.Invoke(sender, data)
		}
	}}
	disconnectedCallback := TcpServerDisconnectedCallback{Op: func(sender network.Connection) {
		ts.curBacklog -= 1
		if ts.curBacklog > ts.maxBacklog { // 如果数量超过则代表正在等待唤醒
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
