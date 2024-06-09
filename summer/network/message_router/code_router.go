package message_router

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type msg struct {
	Code int64
	data interface{}
	Conn interface{}
}

// CodeRouter 消息分发器,通过int64注册
type CodeRouter struct {
	routineCount int           //工作协程数
	workerCount  atomic.Int32  //正在工作的协程数
	Running      bool          //是否正在运行状态
	threadEvent  chan struct{} //通知1个协程停止
	// 消息队列，所有客户端发来的消息都暂存在这里
	messageQueue chan msg
	// 频道字典（订阅记录）
	delegateMap map[int64]func(data, conn interface{})
	waitGroup   *sync.WaitGroup
	mutex       *sync.Mutex
}

func NewCodeRouter(queueSize int) *CodeRouter {
	return &CodeRouter{
		routineCount: 1,
		workerCount:  atomic.Int32{},
		Running:      false,
		messageQueue: make(chan msg, queueSize),
		threadEvent:  make(chan struct{}),
		delegateMap:  map[int64]func(data, conn interface{}){},
		waitGroup:    &sync.WaitGroup{},
		mutex:        &sync.Mutex{},
	}
}

// Subscribe 订阅
func (mr *CodeRouter) Subscribe(code int64, handler func(data, conn interface{})) {
	mr.delegateMap[code] = handler
	fmt.Printf("Subscribe:%v\n", code)
}

// Off 退订
func (mr *CodeRouter) Off(code int64) {
	delete(mr.delegateMap, code)
	fmt.Printf("Off:%v\n", code)
}

// AddMessage 添加新的消息到队列中
func (mr *CodeRouter) AddMessage(code int64, data interface{}, conn interface{}) {
	mr.messageQueue <- msg{code, data, conn}
}

// ClearMsgQueue 清空消息队列
func (mr *CodeRouter) ClearMsgQueue() {
	for {
		select {
		case <-mr.messageQueue:
			// 接收值，继续循环尝试接收
		default:
			// channel为空或已关闭，退出循环
			return
		}
	}
}

// Start 启动路由
func (mr *CodeRouter) Start(routineCount int) {
	if mr.Running {
		return
	}
	mr.Running = true
	mr.routineCount = int(math.Min(math.Max(float64(routineCount), 1), 200))
	for i := 0; i < mr.routineCount; i += 1 {
		mr.waitGroup.Add(i)
		go func(i int) {
			mr.work()
			mr.waitGroup.Done()
		}(i)
	}
	for int(mr.workerCount.Load()) < mr.routineCount {
		time.Sleep(100 * time.Millisecond)
	}
}

// Stop 停止路由
func (mr *CodeRouter) Stop() {
	mr.Running = false
	mr.ClearMsgQueue()
	for mr.workerCount.Load() > 0 {
		mr.threadEvent <- struct{}{}
	}
	mr.waitGroup.Wait()
	time.Sleep(100 * time.Millisecond)
}

func (mr *CodeRouter) work() {
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
		select {
		case msg := <-mr.messageQueue:
			mr.fire(msg)
		case <-mr.threadEvent:
			// TODO:退出工作
		}
	}
	fmt.Println("worker thread end")
}

// fire 触发
func (mr *CodeRouter) fire(m msg) {
	handler, ok := mr.delegateMap[m.Code]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				fmt.Printf("CodeRouter fire error: %v\nStack trace:\n%s\n", err, buf[:n])
			}
		}()
		if handler != nil {
			handler(m.data, m.Conn)
		}
	}
}
