package delegate

import (
	"go/types"
	"testing"
)

func ft1(args ...int) {
	println(args[0] + 1)
}
func TestEvent_AddDelegate(t *testing.T) {
	event1 := Event[int]{}
	event1.AddDelegate(NewDelegate(
		ft1, "a"))
	event1.AddDelegate(NewDelegate(
		func(args ...int) {
			println(args[0] * 2)
		}, "b"))
	event1.AddDelegate(NewDelegate(
		func(args ...int) {
			println(args[0] + 1)
		}, "a"))
	event1.AddDelegate(NewDelegate(
		func(args ...int) {
			println(args[0] - 10)
		}, "c"))
	event1.RemoveDelegates("c")
	event1.Invoke(2)
	event2 := Event[interface{}]{}
	event2.AddDelegate(NewDelegate(func(args ...interface{}) {
		println(args[0].(int))
		println(args[1].(string))
	}, "d"))
	event2.Invoke(1, "hell")
	event3 := Event[types.Nil]{}
	println("有无委托", event3.HasDelegate())
	event3.AddDelegate(NewDelegate(
		func(args ...types.Nil) {
			println("hello world")
		}, "dad"))
	event3.Invoke()
	println("有无委托", event3.HasDelegate())
}
