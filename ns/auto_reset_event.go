package ns

// AutoResetEvent 用于线程同步
type AutoResetEvent struct {
	c chan struct{}
}

func NewAutoResetEvent() *AutoResetEvent {
	return &AutoResetEvent{make(chan struct{}, 1)}
}

func (e *AutoResetEvent) Wait() {
	<-e.c
}

func (e *AutoResetEvent) Set() {
	select {
	case e.c <- struct{}{}:
	default:
	}
}
