package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Logger struct {
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	path        string
}

func (s *Logger) SugarLogger() *zap.SugaredLogger {
	return s.sugarLogger
}

func (s *Logger) Logger() *zap.Logger {
	return s.logger
}

func (s *Logger) Path() string {
	return s.path
}

func (s *Logger) SetPath(path string) {
	s.path = path
}

func NewLogger(path string) *Logger {
	if len(path) == 0 {
		return nil
	}
	s := &Logger{path: path}
	s.init()
	return s
}

func (s *Logger) init() {
	writeSyncer := s.getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	s.sugarLogger = logger.Sugar()
	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (s *Logger) getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create(s.path + time.Now().Format("2006-01-02") + ".log")
	return zapcore.AddSync(file)
}