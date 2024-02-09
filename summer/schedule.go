package summer

import (
	"context"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/timeunit"
	"reflect"
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
	tasks   []*task
	isStart bool
	stop    context.CancelFunc
	fps     int // 每秒帧数
}

func GetScheduleInstance() *Schedule {
	result, _ := singleton.GetOrDo[*Schedule](&singleSchedule, func() (*Schedule, error) {
		return &Schedule{
			tasks:   make([]*task, 0),
			isStart: false,
			stop:    nil,
			fps:     100,
		}, nil
	})
	return result
}

func (s *Schedule) Start() *Schedule {
	if !s.isStart {
		ctx, cancel := context.WithCancel(context.TODO())
		s.stop = cancel
		s.isStart = true
		go s.run(ctx)
	}
	return s
}

func (s *Schedule) Stop() *Schedule {
	s.stop()
	s.isStart = false
	return s
}

func (s *Schedule) run(ctx context.Context) {
	s.runLoop(ctx)
}

func (s *Schedule) AddTask(action func(), timeUnit, timeValue, repeatCount int) {
	interval, err := timeunit.GetInterval(timeValue, timeUnit)
	if err != nil {
		logger.SLCError("Schedule Failed to AddTask: %s", err.Error())
		return
	}
	startTime := time.Now().UnixMilli()
	t := newTask(action, startTime, interval, repeatCount)
	s.tasks = append(s.tasks, t)
}

func (s *Schedule) RemoveTask(action func()) {
	for i, t := range s.tasks {
		if reflect.ValueOf(t.TaskMethod) == reflect.ValueOf(action) {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			break
		}
	}
}

// runLoop 计时器主循环
func (s *Schedule) runLoop(ctx context.Context) {
	// tick间隔
	interval := 1000 / s.fps
	for {
		select {
		case <-ctx.Done():
			return
		default:
			timeunit.Tick()
			startTime := time.Now().UnixMilli()
			// 把完毕的任务移除
			for i, t := range s.tasks {
				if t.Completed {
					s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
				}
			}
			// 执行任务
			for _, t := range s.tasks {
				if t.ShouldRun() {
					t.Run()
				}
			}
			endTime := time.Now().UnixMilli()
			msTime := int64(interval) - (endTime - startTime)
			if msTime > 0 {
				// Sleep for millisecond
				time.Sleep(time.Millisecond * time.Duration(msTime))
			}
		}
	}
}
