package network

import (
	"bytes"
	"encoding/binary"
	"slices"
	"sync"
)

var (
	mutex = sync.Mutex{}
)

type ByteBuffer struct {
	//字节缓存区
	buf []byte
	//读取索引
	readIndex int
	//写入索引
	writeIndex int
	//读取索引标记
	markReadIndex int
	//写入索引标记
	markWriteIndex int
	//缓存区字节数组的长度
	capacity       int
	isLittleEndian bool
}

// NewByteBuffer 要么提供buf, 要么提供capacity来构造
func NewByteBuffer(isLittleEndian bool, buf []byte, capacity int) *ByteBuffer {
	b := &ByteBuffer{isLittleEndian: isLittleEndian}
	if buf != nil {
		b.buf = buf
		b.capacity = len(buf)
	} else {
		b.buf = make([]byte, capacity)
		b.capacity = capacity
	}
	return b
}

/*FixLength
* 根据length长度，确定大于此length的最近的2次方数，如length=7，则返回值为8
 */
func FixLength(length int) int {
	n := 2
	b := 2
	for b < length {
		b = 2 << n
		n++
	}
	return b
}

/*Flip
* 翻转字节数组，如果本地字节序列为低字节序列，则进行翻转以转换为高字节序列
 */
func (bf *ByteBuffer) Flip() {
	slices.Reverse(bf.buf)
}

/*FixSizeAndReset
 * 确定内部字节缓存数组的大小
 */
func (bf *ByteBuffer) FixSizeAndReset(currLen int, futureLen int) int {
	if futureLen > currLen {
		//以原大小的2次方数的两倍确定内部字节缓存区大小
		size := FixLength(currLen) * 2
		if futureLen > size {
			//以将来的大小的2次方的两倍确定内部字节缓存区大小
			size = FixLength(futureLen) * 2
		}
		newBuf := make([]byte, size)
		copy(newBuf, bf.buf[0:currLen])
		bf.buf = newBuf
		bf.capacity = len(newBuf)
	}
	return futureLen
}

/*WriteBytes
 * 将bytes字节数组从startIndex开始的length字节写入到此缓存区
 */
func (bf *ByteBuffer) WriteBytes(bytes []byte, startIndex int, length int) {
	offset := length - startIndex
	if offset <= 0 {
		return
	}
	total := offset + bf.writeIndex
	l := len(bf.buf)
	bf.FixSizeAndReset(l, total)
	i := bf.writeIndex
	j := startIndex
	for i < total {
		bf.buf[i] = bytes[j]
		i += 1
		j += 1
	}
	bf.writeIndex = total
}

func (bf *ByteBuffer) WriteUint16(uint162 uint16) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, uint162)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, uint162)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteUint32(uint322 uint32) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, uint322)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, uint322)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteUInt64(uint642 uint64) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, uint642)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, uint642)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteInt16(int162 int16) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, int162)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, int162)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteInt32(int322 int32) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, int322)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, int322)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteInt64(int642 int64) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, int642)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, int642)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteFloat32(float322 float32) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, float322)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, float322)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) WriteFloat64(float642 float64) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if bf.isLittleEndian {
		binary.Write(bytesBuffer, binary.LittleEndian, float642)
	} else {
		binary.Write(bytesBuffer, binary.BigEndian, float642)
	}
	b := bytesBuffer.Bytes()
	bf.WriteBytes(b, 0, len(b))
}

func (bf *ByteBuffer) Write(buffer *ByteBuffer) {
	if buffer == nil {
		return
	}
	if buffer.ReadableBytes() <= 0 {
		return
	}
	array := buffer.ToArray()
	bf.WriteBytes(array, 0, len(array))
}

/*WriteByte
 * 写入一个byte数据
 */
func (bf *ByteBuffer) WriteByte(value byte) {
	mutex.Lock()
	afterLen := bf.writeIndex + 1
	l := len(bf.buf)
	bf.FixSizeAndReset(l, afterLen)
	bf.buf[bf.writeIndex] = value
	bf.writeIndex = afterLen
	mutex.Unlock()
}

func (bf *ByteBuffer) WriteString(string2 string) {
	if string2 == "" {
		bf.WriteInt32(-1)
	} else {
		bf.WriteInt32(int32(len(string2)))
		bf.WriteBytes([]byte(string2), 0, len(string2))
	}
}

func (bf *ByteBuffer) ReadString() string {
	l := bf.ReadInt32()
	if l < 0 {
		return ""
	}
	bs := bf.Read(int(l))
	return string(bs)
}

/*ReadByte
 * 读取一个字节
 */
func (bf *ByteBuffer) ReadByte() byte {
	b := bf.buf[bf.readIndex]
	bf.readIndex += 1
	return b
}

/*Read
 * 从读取索引位置开始读取len长度的字节数组
 */
func (bf *ByteBuffer) Read(len int) []byte {
	newBytes := make([]byte, len)
	copy(newBytes, bf.buf[bf.readIndex:bf.readIndex+len])
	bf.readIndex += len
	return newBytes
}
func (bf *ByteBuffer) ReadFloat64() float64 {
	bytesBuffer := bytes.NewBuffer(bf.Read(8))
	var x float64
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}
func (bf *ByteBuffer) ReadFloat32() float32 {
	bytesBuffer := bytes.NewBuffer(bf.Read(4))
	var x float32
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}
func (bf *ByteBuffer) ReadUInt64() uint64 {
	bytesBuffer := bytes.NewBuffer(bf.Read(8))
	var x uint64
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}

func (bf *ByteBuffer) ReadUInt32() uint32 {
	bytesBuffer := bytes.NewBuffer(bf.Read(4))
	var x uint32
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}

func (bf *ByteBuffer) ReadUInt16() uint16 {
	bytesBuffer := bytes.NewBuffer(bf.Read(2))
	var x uint16
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}
func (bf *ByteBuffer) ReadInt64() int64 {
	bytesBuffer := bytes.NewBuffer(bf.Read(8))
	var x int64
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}
func (bf *ByteBuffer) ReadInt32() int32 {
	bytesBuffer := bytes.NewBuffer(bf.Read(4))
	var x int32
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}
func (bf *ByteBuffer) ReadInt16() int16 {
	bytesBuffer := bytes.NewBuffer(bf.Read(2))
	var x int16
	if bf.isLittleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &x)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &x)
	}
	return x
}

/*ReadBytes
 * 从读取索引位置开始读取len长度的字节到disbytes目标字节数组中
 * @params disstart 目标字节数组的写入索引
 */
func (bf *ByteBuffer) ReadBytes(disbytes []byte, disstart int, len int) {
	size := disstart + len
	i := disstart
	for i < size {
		disbytes[i] = bf.ReadByte()
		i += 1
	}
}

/*DiscardReadBytes
 * 清除已读字节并重建缓存区
 */
func (bf *ByteBuffer) DiscardReadBytes() {
	if bf.readIndex <= 0 {
		return
	}
	l := len(bf.buf) - bf.readIndex
	newBuf := make([]byte, l)
	copy(newBuf, bf.buf[bf.readIndex:bf.readIndex+l])
	bf.buf = newBuf
	bf.writeIndex -= bf.readIndex
	bf.markReadIndex -= bf.readIndex
	if bf.markReadIndex < 0 {
		bf.markReadIndex = bf.readIndex
	}
	bf.markWriteIndex -= bf.readIndex
	if bf.markWriteIndex < 0 || bf.markWriteIndex < bf.readIndex || bf.markWriteIndex < bf.markReadIndex {
		bf.markWriteIndex = bf.writeIndex
	}
	bf.readIndex = 0
}

/*Clear
 * 清空此对象
 */
func (bf *ByteBuffer) Clear() {
	bf.buf = make([]byte, len(bf.buf))
	bf.readIndex = 0
	bf.writeIndex = 0
	bf.markReadIndex = 0
	bf.markWriteIndex = 0
}

/*SetReaderIndex
 * 设置开始读取的索引
 */
func (bf *ByteBuffer) SetReaderIndex(index int) {
	if index < 0 {
		return
	}
	bf.readIndex = index
}

/*MarkReaderIndex
 * 标记读取的索引位置
 */
func (bf *ByteBuffer) MarkReaderIndex() {
	bf.markReadIndex = bf.readIndex
}

/*MarkWriterIndex
 * 标记写入的索引位置
 */
func (bf *ByteBuffer) MarkWriterIndex() {
	bf.markWriteIndex = bf.writeIndex
}

/*ResetReaderIndex
 * 将读取的索引位置重置为标记的读取索引位置
 */
func (bf *ByteBuffer) ResetReaderIndex() {
	bf.readIndex = bf.markReadIndex
}

/*ResetWriterIndex
 * 将写入的索引位置重置为标记的写入索引位置
 */
func (bf *ByteBuffer) ResetWriterIndex() {
	bf.writeIndex = bf.markWriteIndex
}

/*ReadableBytes
 * 可读的有效字节数
 */
func (bf *ByteBuffer) ReadableBytes() int {
	return bf.writeIndex - bf.readIndex
}

/*ToArray
 * 获取可读的字节数组
 */
func (bf *ByteBuffer) ToArray() []byte {
	bs := make([]byte, bf.writeIndex)
	copy(bs, bf.buf[0:len(bs)])
	return bs
}

/*GetCapacity
 * 获取缓存区大小
 */
func (bf *ByteBuffer) GetCapacity() int {
	return bf.capacity
}
