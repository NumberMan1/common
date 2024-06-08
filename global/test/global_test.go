package test

import (
	"fmt"
	"github.com/NumberMan1/common/global"
	"github.com/NumberMan1/common/global/variable"
	"go.uber.org/zap"
	"testing"
)

func TestInit(t *testing.T) {
	global.Init("tsconfig.json")
	variable.Log.Debug("Initializing", zap.String("debug", "tsconfig.json"))
	variable.Log.Info("Initializing", zap.String("info", "tsconfig.json"))
	variable.Cache.Set([]byte("hello"), []byte("world"), 0)
	value, _ := variable.Cache.Get([]byte("hello"))
	fmt.Println(string(value))

}
