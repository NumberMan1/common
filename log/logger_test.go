package log

import "testing"

func TestLogger(t *testing.T) {
	logger := NewLogger("C:\\Users\\hp\\Desktop\\Project\\fo_go\\self-work\\common\\log\\test")
	logger.SugarLogger().Info("hello world!")
}
