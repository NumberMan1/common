package timeunit

import (
	"errors"
	"time"
)

const (
	Milliseconds = iota
	Seconds
	Minutes
	Hours
	Days
)

var (
	DeltaTime float64
	LastTick  int64
)

func init() {
	DeltaTime = 0
	LastTick = 0
}

func GetInterval(timeValue, timeUint int) (int64, error) {
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
	if LastTick == 0 {
		LastTick = now
	}
	DeltaTime = float64(now-LastTick) / 1000.0
	LastTick = now
}
