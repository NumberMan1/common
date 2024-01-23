package test

import (
	"github.com/NumberMan1/common/ns"
	"testing"
)

type Ft1 struct {
}

func (Ft1) Operator(args ...any) {
	println(args[0].(int) + 1)
}
func TestEvent_AddDelegate(t *testing.T) {
	var ft1 Ft1
	event1 := ns.Event[Ft1]{}
	event1.AddDelegate(ns.NewDelegate(
		ft1, "a"))
	event1.RemoveDelegates("c")
	event1.Invoke(2)
}
