package network

import (
	"fmt"
	pt "github.com/NumberMan1/common/protocol/gen/proto"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
)

// NetService / 网络服务
type NetService struct {
	// 网络监听器
	listener *TcpSocketListener
}

func NewNetService() *NetService {
	return &NetService{}
}

func (ns *NetService) Init(port int) error {
	listener, err := NewTcpSocketListener("0.0.0.0:" + strconv.FormatInt(int64(port), 10))
	if err != nil {
		return err
	}
	listener.SetSocketConnected(SocketConnectedCB{Op: ns.onClientConnected})
	ns.listener = listener
	return nil
}

func (ns *NetService) Start(threadCount int) {
	go ns.listener.Start()
	GetMessageRouterInstance().Start(threadCount)
}

func (ns *NetService) onClientConnected(tsl *TcpSocketListener, socket net.Conn) {
	NewNetConnection(socket, DataReceivedCallback{Op: ns.onDataReceived}, DisconnectedCallback{Op: ns.onDisconnected})
}

func (ns *NetService) onDisconnected(sender *NetConnection) {
	fmt.Println("连接断开")
}

func (ns *NetService) onDataReceived(sender *NetConnection, data []byte) {
	p := &pt.Package{}
	//_ = proto.Unmarshal(data, p)
	err := proto.Unmarshal(data, p)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("%v", int(p.Id))
	//fmt.Printf("%v", data)
	//fmt.Println("want")
	//x := &pt.Package{}
	//x.Id = pt.PackageID_userLoginRequest
	//x.Request = &pt.Request{}
	//x.Request.UserLogin = &pt.UserLoginRequest{}
	//x.Request.UserLogin.Username = "xiazm"
	//x.Request.UserLogin.Password = "123456"
	//marshal, _ := proto.Marshal(x)
	//fmt.Printf("%v", marshal)
	//fmt.Printf("用户名%v", p.Request.UserLogin.Username)
	GetMessageRouterInstance().AddMessage(Msg{
		ID:      int32(p.Id),
		Sender:  sender,
		Message: p,
	})
}
