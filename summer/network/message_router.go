package network

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/ns/singleton"
	"google.golang.org/protobuf/proto"
)

type Msg struct {
	Sender  Connection
	Message proto.Message
}

type MessageHandler struct {
	Op func(msg Msg)
}

func (m MessageHandler) Operator(args ...any) {
	m.Op(args[0].(Msg))
}

// MessageRouter 消息分发器
// 都应通过GetMessageRouterInstance来获取对象
type MessageRouter struct {
	threadCount int                //工作协程数
	workerCount atomic.Int32       //正在工作的协程数
	Running     bool               //是否正在运行状态
	threadEvent *ns.AutoResetEvent //通过Set每次可以唤醒1个线程
	con         sync.Cond
	// 消息队列，所有客户端发来的消息都暂存在这里
	messageQueue *ns.TSQueue[Msg]
	// 频道字典（订阅记录）
	delegateMap map[string]ns.Event[MessageHandler]
	waitGroup   sync.WaitGroup
	mutex       sync.Mutex
}

var (
	singleMessageRouter = singleton.Singleton{}
)

func GetMessageRouterInstance() *MessageRouter {
	instance, _ := singleton.GetOrDo[*MessageRouter](&singleMessageRouter, func() (*MessageRouter, error) {
		return &MessageRouter{
			threadCount:  1,
			workerCount:  atomic.Int32{},
			Running:      false,
			threadEvent:  ns.NewAutoResetEvent(),
			messageQueue: ns.NewTSQueue[Msg](),
			delegateMap:  map[string]ns.Event[MessageHandler]{},
			waitGroup:    sync.WaitGroup{},
			mutex:        sync.Mutex{},
		}, nil
	})
	return instance
}

// Subscribe 订阅
func (mr *MessageRouter) Subscribe(name string, handler MessageHandler) {
	_, ok := mr.delegateMap[name]
	if !ok {
		mr.delegateMap[name] = ns.Event[MessageHandler]{}
	}
	d := mr.delegateMap[name]
	d.AddDelegate(ns.NewDelegate(handler, reflect.ValueOf(handler.Op).String()))
	mr.delegateMap[name] = d
	logger.SLCDebug("Subscribe:%v", name)
}

// Off 退订
func (mr *MessageRouter) Off(name string, handler MessageHandler) {
	_, ok := mr.delegateMap[name]
	if !ok {
		mr.delegateMap[name] = ns.Event[MessageHandler]{}
	}
	d := mr.delegateMap[name]
	d.RemoveDelegates(reflect.ValueOf(handler.Op).String())
	logger.SLCDebug("Off:%v", name)
}

// fire 触发
func (mr *MessageRouter) fire(msg Msg) {
	name := string(msg.Message.ProtoReflect().Descriptor().FullName())
	d, ok := mr.delegateMap[name]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				logger.SLCError("MessageRouter fire error %v", err)
			}
		}()
		if d.HasDelegate() {
			d.Invoke(msg)
		}
	}
}

// AddMessage 添加新的消息到队列中
func (mr *MessageRouter) AddMessage(msg Msg) {
	mr.messageQueue.Push(msg)
	mr.threadEvent.Set()
}

func (mr *MessageRouter) Start(threadCount int) {
	if mr.Running {
		return
	}
	mr.Running = true
	mr.threadCount = min(max(threadCount, 1), 200)
	for i := 0; i < mr.threadCount; i += 1 {
		mr.waitGroup.Add(i)
		go func(i int) {
			mr.messageWork()
			mr.waitGroup.Done()
		}(i)
	}
	for int(mr.workerCount.Load()) < mr.threadCount {
		time.Sleep(100 * time.Millisecond)
	}
}

func (mr *MessageRouter) Stop() {
	mr.Running = false
	mr.messageQueue.Clear()
	for mr.workerCount.Load() > 0 {
		mr.threadEvent.Set()
	}
	mr.waitGroup.Wait()
	time.Sleep(100 * time.Millisecond)
}

func (mr *MessageRouter) messageWork() {
	logger.SLCInfo("worker thread start")
	defer func() {
		if err := recover(); err != nil {
			logger.SLCError("MessageRouter messageWork error %v", err)
		}
		mr.workerCount.Add(-1)
	}()
	mr.workerCount.Add(1)
	for mr.Running {
		if mr.messageQueue.Empty() {
			mr.threadEvent.Wait() //可以通过Set()唤醒
			continue
		}
		//从消息队列取出一个元素
		mr.mutex.Lock()
		if mr.messageQueue.Empty() {
			continue
		}
		msg := mr.messageQueue.Pop()
		mr.mutex.Unlock()
		if p := msg.Message; p != nil {
			mr.executeMessage(msg)
		}
	}
	logger.SLCInfo("worker thread end")
}

// executeMessage 处理消息
func (mr *MessageRouter) executeMessage(msg Msg) {
	//触发订阅
	//fmt.Println(msg)
	//if mr.msgOKHandler(msg) {
	mr.fire(msg)
	//}
}
