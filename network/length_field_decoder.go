package network

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/NumberMan1/common/delegate"
	"github.com/pkg/errors"
	"net"
	"time"
)

// DataReceivedEventHandler 成功收到消息的委托事件
// 实际的DataReceivedEventHandler参数为*LengthFieldDecoder, []byte
// type DataReceivedEventHandler func(*LengthFieldDecoder, []byte)
type DataReceivedEventHandler func(...interface{})

// DisconnectedEventHandler 连接断开的委托事件
// 实际只处理第一个*LengthFieldDecoder
type DisconnectedEventHandler func(...*LengthFieldDecoder)

type LengthFieldDecoder struct {
	isStart           bool
	mSocket           net.Conn
	lengthFieldOffset int //第几个是body长度字段
	lengthFieldLength int //长度字段占几个字节
	//长度字段和内容之间距离几个字节，
	//负数代表向前偏移，body实际长度要减去这个绝对值
	lengthAdjustment    int
	initialBytesToStrip int    //结果数据中前几个字节不需要的字节数
	mBuffer             []byte //接收数据的缓存空间
	mOffset             int    //读取位置
	mSize               int    ///	一次性接收数据的最大字节，默认64k
	// 连接断开的委托事件
	mDisconnectEvent delegate.Event[*LengthFieldDecoder]
	mReceivedEvent   delegate.Event[interface{}]
}

func (d *LengthFieldDecoder) AddDisconnectCB(handler DisconnectedEventHandler, tag string) {
	d.mDisconnectEvent.AddDelegate(delegate.NewDelegate[*LengthFieldDecoder](handler, tag))
}

func (d *LengthFieldDecoder) AddDataReceivedCB(handler DataReceivedEventHandler, tag string) {
	d.mReceivedEvent.AddDelegate(delegate.NewDelegate[interface{}](handler, tag))
}

func NewLengthFieldDecoder(conn net.Conn, lengthFieldOffset int,
	lengthFieldLength int, lengthAdjustment int,
	initialBytesToStrip int, maxBufferLength int) *LengthFieldDecoder {
	return &LengthFieldDecoder{
		isStart:             false,
		mSocket:             conn,
		lengthFieldOffset:   lengthFieldOffset,
		lengthFieldLength:   lengthFieldLength,
		lengthAdjustment:    lengthAdjustment,
		initialBytesToStrip: initialBytesToStrip,
		mSize:               maxBufferLength,
		mBuffer:             make([]byte, maxBufferLength),
		mDisconnectEvent:    delegate.Event[*LengthFieldDecoder]{},
		mReceivedEvent:      delegate.Event[interface{}]{},
	}
}

func NewLengthFieldDecoderDefault(conn net.Conn, lengthFieldOffset int,
	lengthFieldLength int) *LengthFieldDecoder {
	return NewLengthFieldDecoder(conn, lengthFieldOffset,
		lengthFieldLength, 0, lengthFieldLength, 64*1024)
}

func (d *LengthFieldDecoder) Start(ctx context.Context) {
	if d.mSocket != nil && !d.isStart {
		d.beginAsyncReceive(context.WithCancel(ctx))
		d.isStart = true
	}
}
func (d *LengthFieldDecoder) beginAsyncReceive(ctx context.Context, cancelFunc context.CancelFunc) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("LengthFieldDecoder's beginAsyncReceive recover error %v\n", err)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				d.doDisconnected()
				return
			default:
				// 最多等待10秒
				err := d.mSocket.SetReadDeadline(time.Now().Add(time.Second * 3))
				if err != nil {
					cancelFunc()
					continue
				}
				n, err := d.mSocket.Read(d.mBuffer[d.mOffset:])
				d.receive(cancelFunc, n, err)
			}
		}
	}()
}
func (d *LengthFieldDecoder) receive(cancelFunc context.CancelFunc, len int, err error) {
	// 连接出错
	if len == 0 || err != nil {
		cancelFunc()
		return
	}
	// 继续接受数据
	err = d.onReceiveData(len)
	if err != nil {
		fmt.Printf("LengthFieldDecoder'receive func onReceiveData error: %v\n", err)
		cancelFunc()
	}
}
func (d *LengthFieldDecoder) doDisconnected() {
	if d.mDisconnectEvent.HasDelegate() {
		d.mDisconnectEvent.Invoke()
	}
	if d.mSocket != nil {
		_ = d.mSocket.Close()
		d.mSocket = nil
	}
}
func (d *LengthFieldDecoder) onReceiveData(len int) error {
	//headLen+bodyLen=totalLen
	headLen := d.lengthFieldOffset + d.lengthFieldLength
	adj := d.lengthAdjustment //body偏移量
	//size是待处理的数据长度，mOffset每次都从0开始，
	//循环开始之前mOffset代表上次剩余长度
	size := len
	if d.mOffset > 0 {
		size += d.mOffset
		d.mOffset = 0
	}
	//循环解析
	for {
		remain := size - d.mOffset //剩余未处理的长度

		//如果未处理的数据超出限制
		if remain > d.mSize {
			return errors.New("未处理的数据超出限制")
		}
		if remain < headLen {
			//接收的数据不够一个完整的包，继续接收
			copy(d.mBuffer[0:remain], d.mBuffer[d.mOffset:])
			d.mOffset = remain
			return nil
		}

		//获取包长度
		temp := make([]byte, 4)
		lenOffset := d.mOffset + d.lengthFieldOffset
		copy(temp, d.mBuffer[lenOffset:lenOffset+4])
		bytesBuffer := bytes.NewBuffer(temp)
		var bodyLen int32
		_ = binary.Read(bytesBuffer, binary.BigEndian, &bodyLen)
		if remain < headLen+adj+int(bodyLen) {
			copy(d.mBuffer[0:remain], d.mBuffer[d.mOffset:])
			//接收的数据不够一个完整的包，继续接收
			d.mOffset = remain
			return nil
		}

		////body的读取位置
		//bodyStart := d.mOffset + max(headLen, headLen+adj)
		////body的真实长度
		//bodyCount := min(int(bodyLen), int(bodyLen)+adj)
		//fmt.Printf("bodyStart=%v, bodyCount=%v, remain=%v\n", bodyStart, bodyCount, remain)
		//获取包体
		total := headLen + adj + int(bodyLen) //数据包总长度
		count := total - d.initialBytesToStrip
		data := make([]byte, count)
		copy(data[0:count], d.mBuffer[d.mOffset+d.initialBytesToStrip:])
		d.mOffset += total

		//完成一个数据包
		if d.mReceivedEvent.HasDelegate() {
			d.mReceivedEvent.Invoke(d, data)
			//fmt.Println("完成一个数据包")
			break
		}
	}
	return nil
}
