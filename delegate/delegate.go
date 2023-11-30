package delegate

//type Delegate[T any] func(args ...T)

type Delegate[T any] struct {
	Op func(args ...T)
	// 用于识别不同的Delegate
	Tag string
}

// NewDelegate 若op不需要参数类型直接传types.Nil
func NewDelegate[T any](op func(args ...T), tag string) *Delegate[T] {
	return &Delegate[T]{Op: op, Tag: tag}
}

type Event[T any] struct {
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

func (e *Event[T]) Invoke(args ...T) {
	for _, d := range e.delegates {
		d.Op(args...)
	}
}

func (e *Event[T]) HasDelegate() bool {
	return len(e.delegates) > 0
}
