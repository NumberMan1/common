package log

import (
	"fmt"
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
	s.logger = logger
	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (s *Logger) getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile(s.path+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}

func (s *Logger) SugarError(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.sugarLogger.Error(sprintf)
}

func (s *Logger) SugarInfo(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.sugarLogger.Info(sprintf)
}

func (s *Logger) LoggerError(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.logger.Error(sprintf)
}

func (s *Logger) LoggerInfo(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.logger.Info(sprintf)
}

func (s *Logger) LoggerFatal(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.logger.Fatal(sprintf)
}

func (s *Logger) SugarFatal(format string, a ...any) {
	sprintf := fmt.Sprintf(format, a)
	s.logger.Fatal(sprintf)
}
