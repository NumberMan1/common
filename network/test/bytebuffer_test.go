package test

import (
	"fmt"
	"github.com/NumberMan1/common/network"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	buffer := network.NewByteBufferByCapacity(false, 1024)
	buffer.WriteFloat64(3.3)
	buffer.WriteFloat32(3.2)
	buffer.WriteInt64(-142)
	buffer.WriteUInt64(125)
	buffer.WriteString("你好 14dsa")
	f64 := buffer.ReadFloat64()
	f32 := buffer.ReadFloat32()
	readInt64 := buffer.ReadInt64()
	uInt64 := buffer.ReadUInt64()
	readString := buffer.ReadString()
	fmt.Println(f64)
	fmt.Println(f32)
	fmt.Println(readInt64)
	fmt.Println(uInt64)
	fmt.Println(readString)
}
