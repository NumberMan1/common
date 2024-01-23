package ns

//type Delegate[T any] func(args ...T)

type Func interface {
	Operator(args ...any)
}

type Delegate[T Func] struct {
	//Op func(args ...T)
	Op T
	// 用于识别不同的Delegate
	Tag string
}

// NewDelegate 创建用于执行回调的委托器
func NewDelegate[T Func](op T, tag string) *Delegate[T] {
	return &Delegate[T]{Op: op, Tag: tag}
}

type Event[T Func] struct {
	delegates []*Delegate[T]
}

func (e *Event[T]) AddDelegate(d *Delegate[T]) {
	e.delegates = append(e.delegates, d)
}

func (e *Event[T]) RemoveDelegates(tag string) {
	for i, d := range e.delegates {
		if d.Tag == tag {
			e.delegates = append(e.delegates[:i], e.delegates[i+1:]...)
			break
		}
	}
}

func (e *Event[T]) Invoke(args ...any) {
	for _, d := range e.delegates {
		d.Op.Operator(args...)
	}
}

func (e *Event[T]) HasDelegate() bool {
	return len(e.delegates) > 0
}

func (e *Event[T]) Length() int {
	return len(e.delegates)
}
