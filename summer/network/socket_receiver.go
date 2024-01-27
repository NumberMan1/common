package network

import (
	"context"
	"github.com/NumberMan1/common/logger"
	"net"
	"time"
)

type DataReceivedEventHandler func(data []byte) error

type DisconnectedEventHandler func()

type SocketReceiver struct {
	buf          []byte
	startIndex   int
	mSocket      net.Conn
	DataReceived DataReceivedEventHandler
	Disconnected DisconnectedEventHandler
}

func NewSocketReceiver(mSocket net.Conn) *SocketReceiver {
	return &SocketReceiver{mSocket: mSocket, buf: make([]byte, 64*1024), startIndex: 0}
}

func (s *SocketReceiver) Start() {
	s.beginReceive(context.WithCancel(context.TODO()))
}

func (s *SocketReceiver) beginReceive(ctx context.Context, cancelFunc context.CancelFunc) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.SLCError("SocketReceiver's beginReceive recover error %v", err)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				s.disconnected()
				return
			default:
				// 最多等待10秒
				err := s.mSocket.SetReadDeadline(time.Now().Add(time.Second * 10))
				if err != nil {
					cancelFunc()
					continue
				}
				n, err := s.mSocket.Read(s.buf[s.startIndex:])
				s.receiveCB(cancelFunc, n, err)
			}
		}
	}()
}

func (s *SocketReceiver) receiveCB(cancelFunc context.CancelFunc, l int, err error) {
	// 连接出错
	if l == 0 || err != nil {
		cancelFunc()
		return
	}
	// 继续接受数据
	s.doReceive(l)
}

func (s *SocketReceiver) doReceive(l int) {
	remain, offset := s.startIndex+l, 0
	for remain > 4 {
		msgLen := s.getInt32BE(s.buf, offset)
		if remain < int(msgLen+4) {
			break
		}
		data := make([]byte, msgLen)
		copy(data, s.buf[offset+4:int(msgLen)+offset+4])
		if s.DataReceived != nil {
			err := s.DataReceived(data)
			if err != nil {
				logger.SLCError("SocketReceiver消息解析异常: %v\n", err.Error())
			}
		}
		offset += int(msgLen) + 4
		remain -= int(msgLen) + 4
	}
	if remain > 0 {
		copy(s.buf, s.buf[offset:remain+offset])
	}
	s.startIndex = remain
}

func (s *SocketReceiver) disconnected() {
	if s.Disconnected != nil {
		s.Disconnected()
		_ = s.mSocket.Close()
	}
	s.mSocket = nil
}

// 获取大端模式int值
func (s *SocketReceiver) getInt32BE(data []byte, index int) int32 {
	return int32(data[index])<<0x18 | int32(data[index+1])<<0x10 | int32(data[index+2])<<8 | int32(data[index+3])
}
