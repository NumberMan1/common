package core

import (
	"bytes"
	"encoding/binary"
)

type DataStream struct {
	*bytes.Buffer
	isLittleEndian bool
}

var DataStreamPool = make(chan *DataStream, 200)

func newDataStream(size int) *DataStream {
	return &DataStream{
		Buffer:         bytes.NewBuffer(make([]byte, size)),
		isLittleEndian: false, // 统一使用大端模式
	}
}

// AllocateDataStream 申请流, 如果池里有就拿出来,没有就重新创建一个,默认buffer为1024字节
func AllocateDataStream() *DataStream {
	select {
	case stream := <-DataStreamPool:
		return stream
	default:
		return newDataStream(1024)
	}
}

func (ds *DataStream) Write(p []byte) (n int, err error) {
	return ds.Buffer.Write(p)
}

func (ds *DataStream) Read(p []byte) (n int, err error) {
	return ds.Buffer.Read(p)
}

func (ds *DataStream) Release() {
	ds.Buffer.Reset()
	select {
	case DataStreamPool <- ds:
	default:
	}
}

func (ds *DataStream) WriteUInt16(u uint16) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, u)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, u)
	}
	return
}

func (ds *DataStream) WriteUInt32(u uint32) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, u)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, u)
	}
	return
}

func (ds *DataStream) WriteUInt64(u uint64) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, u)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, u)
	}
	return
}

func (ds *DataStream) WriteInt16(i int16) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, i)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, i)
	}
	return
}

func (ds *DataStream) WriteInt32(i int32) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, i)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, i)
	}
	return
}

func (ds *DataStream) WriteInt64(i int64) (err error) {
	if ds.isLittleEndian {
		err = binary.Write(ds.Buffer, binary.LittleEndian, i)
	} else {
		err = binary.Write(ds.Buffer, binary.BigEndian, i)
	}
	return
}

func (ds *DataStream) ReadUInt16() (u uint16, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &u)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &u)
	}
	return
}

func (ds *DataStream) ReadUInt32() (u uint32, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &u)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &u)
	}
	return
}

func (ds *DataStream) ReadUInt64() (u uint64, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &u)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &u)
	}
	return
}

func (ds *DataStream) ReadInt16() (i int16, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &i)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &i)
	}
	return
}

func (ds *DataStream) ReadInt32() (i int32, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &i)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &i)
	}
	return
}

func (ds *DataStream) ReadInt64() (i int64, err error) {
	if ds.isLittleEndian {
		err = binary.Read(ds.Buffer, binary.LittleEndian, &i)
	} else {
		err = binary.Read(ds.Buffer, binary.BigEndian, &i)
	}
	return
}
