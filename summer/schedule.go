package summer

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/timeunit"
	"reflect"
	"sync"
	"time"
)

var (
	singleSchedule = singleton.Singleton{}
)

type task struct {
	TaskMethod   func()
	StartTime    int64
	Interval     int64
	RepeatCount  int
	currentCount int
	lastTick     int64 //上一次执行开始的时间
	Completed    bool  //是否已经执行完毕
}

func (t *task) Run() {
	t.lastTick = time.Now().UnixMilli()
	defer func() {
		if err := recover(); err != nil {
			logger.SLCError("task Run error %v", err)
		}
	}()
	t.TaskMethod()
	t.currentCount += 1
	if (t.currentCount == t.RepeatCount) && (t.RepeatCount != 0) {
		logger.SLCDebug("Task complete")
		t.Completed = true
	}
}

func (t *task) ShouldRun() bool {
	if (t.currentCount == t.RepeatCount) && (t.RepeatCount != 0) {
		logger.SLCInfo("RepeatCount=%v", t.RepeatCount)
		return false
	}
	now := time.Now().UnixMilli()
	if (now >= t.StartTime) && ((now - t.lastTick) >= t.Interval) {
		return true
	}
	return false
}

// repeatCount 为0表示无限重复
func newTask(taskMethod func(), startTime int64, interval int64, repeatCount int) *task {
	return &task{
		TaskMethod:   taskMethod,
		StartTime:    startTime,
		Interval:     interval,
		RepeatCount:  repeatCount,
		currentCount: 0,
		lastTick:     0,
		Completed:    false,
	}
}

type Schedule struct {
	tasks       []*task
	addQueue    *ns.TSQueue[*task]
	removeQueue *ns.TSQueue[func()]
	isStart     bool
	stop        chan struct{}
	fps         int // 每秒帧数
	ticker      *time.Ticker
	next        int64 //下一帧执行的时间
	mutex       sync.Mutex
}

func GetScheduleInstance() *Schedule {
	result, _ := singleton.GetOrDo[*Schedule](&singleSchedule, func() (*Schedule, error) {
		return &Schedule{
			tasks:       make([]*task, 0),
			addQueue:    ns.NewTSQueue[*task](),
			removeQueue: ns.NewTSQueue[func()](),
			isStart:     false,
			stop:        make(chan struct{}, 1),
			fps:         50,
			ticker:      nil,
			mutex:       sync.Mutex{},
		}, nil
	})
	return result
}

func (s *Schedule) Start() *Schedule {
	if !s.isStart {
		s.isStart = true
		s.ticker = time.NewTicker(1 * time.Millisecond)
		go s.execute()
	}
	return s
}

func (s *Schedule) Stop() *Schedule {
	if s.isStart {
		s.ticker.Stop()
		s.ticker = nil
		s.isStart = false
		s.stop <- struct{}{}
	}
	return s
}

func (s *Schedule) AddTask(action func(), timeUnit timeunit.TimeUnit, timeValue, repeatCount int) {
	interval, err := timeunit.GetInterval(timeValue, timeUnit)
	if err != nil {
		logger.SLCError("Schedule Failed to AddTask: %s", err.Error())
		return
	}
	startTime := time.Now().UnixMilli()
	t := newTask(action, startTime, interval, repeatCount)
	s.addQueue.Push(t)
}

func (s *Schedule) RemoveTask(action func()) {
	s.removeQueue.Push(action)
}

// Update 每帧都会执行
func (s *Schedule) Update(action func()) {
	t := newTask(action, 0, 0, 0)
	s.addQueue.Push(t)
}

func (s *Schedule) execute() {
	for {
		select {
		case <-s.ticker.C:
			// tick间隔
			interval := 1000 / s.fps
			startTime := time.Now().UnixMilli()
			if startTime < s.next {
				return
			}
			s.next = startTime + int64(interval)
			timeunit.Tick()
			s.mutex.Lock()
			//移除队列
			for item := s.removeQueue.Pop(); item != nil; item = s.removeQueue.Pop() {
				for i, task := range s.tasks {
					if reflect.ValueOf(task.TaskMethod) == reflect.ValueOf(item) {
						s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
						break
					}
				}
			}
			// 移除完毕的任务
			for i, task := range s.tasks {
				if task.Completed {
					s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
					break
				}
			}
			// 添加队列任务
			for item := s.addQueue.Pop(); item != nil; item = s.addQueue.Pop() {
				s.tasks = append(s.tasks, item)
			}
			// 执行任务
			for _, task := range s.tasks {
				if task.ShouldRun() {
					task.Run()
				}
			}
			s.mutex.Unlock()
		case <-s.stop:
			return
		}
	}
}
