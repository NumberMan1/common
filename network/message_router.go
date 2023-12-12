package network

import (
	"fmt"
	"github.com/NumberMan1/common/delegate"
	"google.golang.org/protobuf/proto"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Msg struct {
	ID      int32
	Sender  *NetConnection
	Message proto.Message
}

type MessageHandler[T proto.Message] struct {
	ID int32
	Op func(connection *NetConnection, msg T)
}

func (m MessageHandler[T]) Operator(args ...any) {
	m.Op(args[0].(*NetConnection), args[1].(T))
}

type MessageRouter struct {
	ThreadCount int             //工作协程数
	WorkerCount int             //正在工作的协程数
	Running     bool            //是否正在运行状态
	threadEvent *AutoResetEvent //通过Set每次可以唤醒1个线程
	// 消息队列，所有客户端发来的消息都暂存在这里
	messageQueue *TSQueue[Msg]
	// 频道字典（订阅记录）
	delegateMap map[int32]delegate.Event[MessageHandler[proto.Message]]
	// 处理哪种Msg为有效Msg
	msgOKHandler func(msg Msg) bool
}

// SetMsgOKHandler 处理哪种Msg为有效Msg
func (mr *MessageRouter) SetMsgOKHandler(msgOKHandler func(msg Msg) bool) {
	mr.msgOKHandler = msgOKHandler
}

var (
	instance *MessageRouter
	once     sync.Once
)

func GetMessageRouterInstance() *MessageRouter {
	once.Do(func() {
		instance = &MessageRouter{
			ThreadCount:  1,
			WorkerCount:  0,
			Running:      false,
			threadEvent:  NewAutoResetEvent(),
			messageQueue: NewTSQueue[Msg](),
			delegateMap:  map[int32]delegate.Event[MessageHandler[proto.Message]]{},
			msgOKHandler: nil,
		}
	})
	return instance
}

// On 订阅
func (mr *MessageRouter) On(handler MessageHandler[proto.Message]) {
	_, ok := mr.delegateMap[handler.ID]
	if !ok {
		mr.delegateMap[handler.ID] = delegate.Event[MessageHandler[proto.Message]]{}
	}
	d := mr.delegateMap[handler.ID]
	d.AddDelegate(delegate.NewDelegate(handler, strconv.FormatInt(int64(handler.ID), 10)))
	mr.delegateMap[handler.ID] = d
	//fmt.Printf("%v:On\n", handler.ID)
}

// Off 退订
func (mr *MessageRouter) Off(handler MessageHandler[proto.Message]) {
	_, ok := mr.delegateMap[handler.ID]
	if !ok {
		mr.delegateMap[handler.ID] = delegate.Event[MessageHandler[proto.Message]]{}
	}
	d := mr.delegateMap[handler.ID]
	d.RemoveDelegates(strconv.FormatInt(int64(handler.ID), 10))
	//fmt.Printf("%v:Off\n", handler.ID)
}

// fire 触发
func (mr *MessageRouter) fire(msg Msg) {
	//fmt.Println(mr.delegateMap)
	d, ok := mr.delegateMap[msg.ID]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("MessageRouter fire error %v\n", err)
			}
		}()
		//println(d.HasDelegate())
		d.Invoke(msg.Sender, msg.Message)
	}
}

// AddMessage 添加新的消息到队列中
func (mr *MessageRouter) AddMessage(msg Msg) {
	mr.messageQueue.Push(msg)
	mr.threadEvent.Set()
}

func (mr *MessageRouter) Start(threadCount int) {
	mr.Running = true
	mr.ThreadCount = min(max(threadCount, 1), 100)
	for i := 0; i < mr.ThreadCount; i += 1 {
		go mr.messageWork()
	}
	for mr.WorkerCount < mr.ThreadCount {
		time.Sleep(100 * time.Millisecond)
	}
}

func (mr *MessageRouter) Stop() {
	mr.Running = false
	mr.messageQueue.Clear()
	for mr.WorkerCount > 0 {
		mr.threadEvent.Set()
	}
	time.Sleep(100 * time.Millisecond)
}

func (mr *MessageRouter) messageWork() {
	fmt.Println("worker thread start")
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("MessageRouter fire error %v\n", err)
		}
		a := atomic.Int32{}
		a.Store(int32(mr.WorkerCount))
		mr.WorkerCount = int(a.Add(-1))
	}()
	a := atomic.Int32{}
	a.Store(int32(mr.WorkerCount))
	mr.WorkerCount = int(a.Add(1))
	for mr.Running {
		//fmt.Println(mr.messageQueue.Size())
		if mr.messageQueue.Empty() {
			mr.threadEvent.Wait() //可以通过Set()唤醒
			continue
		}
		//从消息队列取出一个元素
		msg := mr.messageQueue.Pop()
		if p := msg.Message; p != nil {
			mr.executeMessage(msg)
		}
	}
}

// executeMessage 递归处理消息
func (mr *MessageRouter) executeMessage(msg Msg) {
	//触发订阅
	//fmt.Println(msg)
	if mr.msgOKHandler(msg) {
		mr.fire(msg)
	}
}
