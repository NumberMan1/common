package ns

type Queue struct {
	data []any
	head int
	tail int
	size int
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: make([]any, size),
		head: 0,
		tail: 0,
		size: size,
	}
}

func (q *Queue) Push(val any) bool {
	if q.Full() {
		return false
	}
	q.data[q.tail] = val
	q.tail = (q.tail + 1) % q.size
	return true
}

func (q *Queue) Pop() any {
	if q.Empty() {
		return nil
	}
	val := q.data[q.head]
	q.head = (q.head + 1) % q.size
	return val
}

func (q *Queue) Empty() bool {
	return q.head == q.tail
}

func (q *Queue) Full() bool {
	return (q.tail+1)%q.size == q.head
}
