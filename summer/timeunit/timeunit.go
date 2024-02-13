package timeunit

import (
	"errors"
	"time"
)

type TimeUnit int

const (
	Milliseconds TimeUnit = iota
	Seconds
	Minutes
	Hours
	Days
)

var (
	StartTime int64
	Time      float64
	DeltaTime float64
	LastTick  int64
)

func init() {
	StartTime = time.Now().UnixMilli()
	Time = 0
	DeltaTime = 0
	LastTick = 0
}

func GetInterval(timeValue int, timeUint TimeUnit) (int64, error) {
	switch timeUint {
	case Milliseconds:
		return int64(timeValue), nil
	case Seconds:
		return int64(timeValue) * 1000, nil
	case Minutes:
		return int64(timeValue) * 1000 * 60, nil
	case Hours:
		return int64(timeValue) * 1000 * 60 * 60, nil
	case Days:
		return int64(timeValue) * 1000 * 60 * 60 * 24, nil
	default:
		return 0, errors.New("invalid timeunit unit")
	}
}

// Tick 由 summer/schedule调用，请不要自行调用，除非你知道自己在做什么！！！
func Tick() {
	now := time.Now().UnixMilli()
	Time = float64(now-StartTime) * 0.001
	if LastTick == 0 {
		LastTick = now
	}
	DeltaTime = float64(now-LastTick) * 0.001
	LastTick = now
}
