package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	d := LogInit(true, zap.InfoLevel, "C:\\Users\\hp\\Desktop\\Project\\fo_go\\self-work\\mmo-server\\common\\logger\\test2")
	d.Debug("debug1")
	d.Debug("debug2")
	d.Info("debug3")
	d.Error("debug4")
	i := LogInit(false, zap.ErrorLevel, "C:\\Users\\hp\\Desktop\\Project\\fo_go\\self-work\\mmo-server\\common\\logger\\test3")
	i.Info("info1")
	i.Info("info2")
	i.Debug("info3")
	i.Error("info4")
	defer d.Sync()
	defer i.Sync()
}
