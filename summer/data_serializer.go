package summer

import (
	"bytes"
	"encoding/binary"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"golang.org/x/exp/slices"
	"io"
	"reflect"
)

// Serialize 接受uint8,uint16,uint32,uint,uint64
// int8,int16,int,int32,int64,float32,float64,bool,string
func Serialize(args []any) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, arg := range args {
		switch reflect.TypeOf(arg).Kind().String() {
		case "int8":
			buf.WriteByte(1)
			buf.WriteByte(byte(arg.(int8)))
		case "uint8":
			buf.WriteByte(2)
			buf.WriteByte(arg.(uint8))
		case "int16":
			buf.WriteByte(3)
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(arg.(int16)))
			buf.Write(b)
		case "uint16":
			buf.WriteByte(4)
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, arg.(uint16))
			buf.Write(b)
		case "int":
			buf.WriteByte(5)
			encode := VariantEncode(uint64(arg.(int)))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "int32":
			buf.WriteByte(5)
			encode := VariantEncode(uint64(arg.(int32)))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "uint32":
			buf.WriteByte(6)
			encode := VariantEncode(uint64(arg.(uint32)))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "int64":
			buf.WriteByte(7)
			encode := VariantEncode(uint64(arg.(int64)))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "uint64":
			buf.WriteByte(8)
			encode := VariantEncode(arg.(uint64))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "float32":
			buf.WriteByte(9)
			d := float32(1000) * arg.(float32)
			encode := VariantEncode(uint64(d))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "float64":
			buf.WriteByte(10)
			d := float64(1000) * arg.(float64)
			encode := VariantEncode(uint64(d))
			buf.WriteByte(byte(len(encode)))
			buf.Write(encode)
		case "bool":
			if arg.(bool) {
				buf.WriteByte(11)
			} else {
				buf.WriteByte(12)
			}
		case "string":
			bs := []byte(arg.(string))
			if ns.IsLittleEndian() {
				slices.Reverse(bs)
			}
			arr := make([]byte, 4)
			binary.BigEndian.PutUint32(arr, uint32(len(bs)))
			buf.WriteByte(13)
			buf.Write(arr)
			buf.Write(bs)
		default:
			buf.WriteByte(0)
			logger.SLCError("DataSerializer无法处理的类型:%v", reflect.TypeOf(arg))
		}
	}
	return buf.Bytes()
}

// Deserialize 对应Serialize的数据进行反序列化
func Deserialize(data []byte) []any {
	buf := bytes.NewBuffer(data)
	list := make([]any, 0)
	for {
		bt, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		switch bt {
		case 0:
			list = append(list, nil)
		case 1:
			readByte, _ := buf.ReadByte()
			list = append(list, int8(readByte))
		case 2:
			readByte, _ := buf.ReadByte()
			list = append(list, uint8(readByte))
		case 3: //int16
			arr := make([]byte, 2)
			_, _ = buf.Read(arr)
			list = append(list, int16(binary.BigEndian.Uint16(arr)))
		case 4: //uint16
			arr := make([]byte, 2)
			_, _ = buf.Read(arr)
			list = append(list, binary.BigEndian.Uint16(arr))
		case 5: //int32
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, int32(VariantDecode(arr)))
		case 6: //uint32
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, uint32(VariantDecode(arr)))
		case 7: //int64
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, int64(VariantDecode(arr)))
		case 8: //uint64
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, VariantDecode(arr))
		case 9: //float32
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, float32(float32(VariantDecode(arr))/float32(1000)))
		case 10: //float64
			readLen, _ := buf.ReadByte()
			arr := make([]byte, int(readLen))
			_, _ = buf.Read(arr)
			list = append(list, float64(float64(VariantDecode(arr))/float64(1000)))
		case 11: //bool-true
			list = append(list, true)
		case 12: //bool-false
			list = append(list, false)
		case 13: //string
			lenBytes := make([]byte, 4)
			_, _ = buf.Read(lenBytes)
			l := binary.BigEndian.Uint32(lenBytes)
			logger.SLCDebug("字符串长度:%v", l)
			arr := make([]byte, l)
			_, _ = buf.Read(arr)
			if ns.IsLittleEndian() {
				slices.Reverse(arr)
			}
			list = append(list, string(arr))
		default:
			logger.SLCError("DataSerializer无法识别的编码:%v", bt)
		}
	}
	return list
}

func VariantEncode(value uint64) []byte {
	l := make([]byte, 0)
	for value > 0 {
		b := byte(value & 0x7f)
		value >>= 7
		if value > 0 {
			b |= 0x80
		}
		l = append(l, b)
	}
	return l
}

func VariantDecode(buffer []byte) uint64 {
	value := uint64(0)
	shift := int32(0)
	for i, l := 0, len(buffer); i < l; i += 1 {
		b := buffer[i]
		value |= uint64(b&0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
	}
	return value
}
