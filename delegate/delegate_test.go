package delegate

import (
	"testing"
)

type Ft1 struct {
}

func (Ft1) Operator(args ...any) {
	println(args[0].(int) + 1)
}
func TestEvent_AddDelegate(t *testing.T) {
	var ft1 Ft1
	event1 := Event[Ft1]{}
	event1.AddDelegate(NewDelegate(
		ft1, "a"))
	event1.RemoveDelegates("c")
	event1.Invoke(2)
}
