package test

import (
	"fmt"
	"github.com/NumberMan1/common/global"
	"github.com/NumberMan1/common/global/variable"
	"go.uber.org/zap"
	"testing"
)

func TestInit(t *testing.T) {
	global.Init("tsconfig.yaml")
	variable.Log.Debug("Initializing", zap.String("debug", "tsconfig.yaml"))
	variable.Log.Info("Initializing", zap.String("info", "tsconfig.yaml"))
	variable.Cache.Set([]byte("hello"), []byte("world"), 0)
	value, _ := variable.Cache.Get([]byte("hello"))
	fmt.Println(string(value))

}
