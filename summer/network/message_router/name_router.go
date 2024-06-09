package message_router

import (
	"fmt"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"google.golang.org/protobuf/proto"
)

type Msg struct {
	Sender  network.Connection
	Message proto.Message
}

type MessageHandler struct {
	Op func(msg Msg)
}

func (m MessageHandler) Operator(args ...any) {
	m.Op(args[0].(Msg))
}

// NameRouter 消息分发器,通过string注册
type NameRouter struct {
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

func NewNameRouter() *NameRouter {
	return &NameRouter{
		threadCount:  1,
		workerCount:  atomic.Int32{},
		Running:      false,
		threadEvent:  ns.NewAutoResetEvent(),
		messageQueue: ns.NewTSQueue[Msg](),
		delegateMap:  map[string]ns.Event[MessageHandler]{},
		waitGroup:    sync.WaitGroup{},
		mutex:        sync.Mutex{},
	}
}

var (
	singleMessageRouter = singleton.Singleton{}
)

func GetMessageRouterInstance() *NameRouter {
	instance, _ := singleton.GetOrDo[*NameRouter](&singleMessageRouter, func() (*NameRouter, error) {
		return NewNameRouter(), nil
	})
	return instance
}

// Subscribe 订阅
func (mr *NameRouter) Subscribe(name string, handler MessageHandler) {
	_, ok := mr.delegateMap[name]
	if !ok {
		mr.delegateMap[name] = ns.Event[MessageHandler]{}
	}
	d := mr.delegateMap[name]
	d.AddDelegate(ns.NewDelegate(handler, reflect.ValueOf(handler.Op).String()))
	mr.delegateMap[name] = d
	fmt.Println("Subscribe:", name)
}

// Off 退订
func (mr *NameRouter) Off(name string, handler MessageHandler) {
	_, ok := mr.delegateMap[name]
	if !ok {
		mr.delegateMap[name] = ns.Event[MessageHandler]{}
	}
	d := mr.delegateMap[name]
	d.RemoveDelegates(reflect.ValueOf(handler.Op).String())
	fmt.Println("Off:", name)
}

// fire 触发
func (mr *NameRouter) fire(msg Msg) {
	name := string(msg.Message.ProtoReflect().Descriptor().FullName())
	d, ok := mr.delegateMap[name]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				fmt.Printf("NameRouter fire error: %v\nStack trace:\n%s\n", err, buf[:n])
			}
		}()
		if d.HasDelegate() {
			d.Invoke(msg)
		}
	}
}

// AddMessage 添加新的消息到队列中
func (mr *NameRouter) AddMessage(msg Msg) {
	mr.messageQueue.Push(msg)
	mr.threadEvent.Set()
}

func (mr *NameRouter) Start(threadCount int) {
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

func (mr *NameRouter) Stop() {
	mr.Running = false
	mr.messageQueue.Clear()
	for mr.workerCount.Load() > 0 {
		mr.threadEvent.Set()
	}
	mr.waitGroup.Wait()
	time.Sleep(100 * time.Millisecond)
}

func (mr *NameRouter) messageWork() {
	fmt.Println("worker thread start")
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			fmt.Printf("Panic: %v\nStack trace:\n%s\n", err, buf[:n])
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
func (mr *NameRouter) executeMessage(msg Msg) {
	//触发订阅
	//fmt.Println(msg)
	//if mr.msgOKHandler(msg) {
	mr.fire(msg)
	//}
}
