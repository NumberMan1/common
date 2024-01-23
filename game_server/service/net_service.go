package service

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/core"
)

type NetService struct {
	tcpServer *core.TcpServer
}

func NewNetService() *NetService {
	server, _ := core.NewTcpServer("127.0.0.1:32510")
	server.SetConnectedCallback(core.TcpServerConnectedCallback{Op: OnClientConnected})
	server.SetDisconnectedCallback(core.TcpServerDisconnectedCallback{Op: OnDisconnected})
	return &NetService{
		tcpServer: server,
	}
}

func (n *NetService) Start() {
	n.tcpServer.Start()
	network.GetMessageRouterInstance().Start(10)
}

func OnClientConnected(conn network.Connection) {
	logger.SLCInfo("客户端接入")
}

func OnDisconnected(conn network.Connection) {
	logger.SLCInfo("连接断开:%v", conn.Socket().RemoteAddr().String())
}
