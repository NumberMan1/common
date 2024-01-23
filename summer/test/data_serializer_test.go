package test

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer"
	"math"
	"testing"
)

func TestSerializeAndDeserialize(t *testing.T) {
	logger.SLoggerConsole.Debugf("测试%v", 2134)
	anies := make([]any, 0)
	anies = append(anies, int8(math.MinInt8))
	anies = append(anies, int8(math.MaxInt8))
	anies = append(anies, uint8(math.MaxUint8))
	anies = append(anies, int16(math.MinInt16))
	anies = append(anies, int16(math.MaxInt16))
	anies = append(anies, uint16(math.MaxUint16))
	anies = append(anies, int32(math.MinInt32))
	anies = append(anies, int32(math.MaxInt32))
	anies = append(anies, uint32(math.MaxUint32))
	anies = append(anies, int64(math.MinInt64))
	anies = append(anies, int64(math.MaxInt64))
	anies = append(anies, uint64(math.MaxUint64))
	anies = append(anies, uint8(127))
	anies = append(anies, int8(-100))
	anies = append(anies, int16(1024))
	anies = append(anies, int16(-100))
	anies = append(anies, 65535)
	anies = append(anies, 3.1)
	anies = append(anies, float32(3.14))
	anies = append(anies, float64(3.1415926))
	anies = append(anies, false)
	anies = append(anies, "毛主席万岁!")
	anies = append(anies, "秋水共长天一色")
	//anies = append(anies, 34234234.532)
	serialize := summer.Serialize(anies)
	logger.SLoggerConsole.Debug(serialize)
	deserialize := summer.Deserialize(serialize)
	logger.SLoggerConsole.Debug(deserialize)
	for _, v := range deserialize {
		logger.SLoggerConsole.Debugf("%v\n", v)
	}
	encode := summer.VariantEncode(uint64(1000 * float64(123456.789)))
	logger.SLoggerConsole.Debug(encode)
	decode := summer.VariantDecode(encode)
	logger.SLoggerConsole.Debug(float64(decode))
}
