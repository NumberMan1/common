package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var SLoggerConsole *zap.SugaredLogger

func init() {
	SLoggerConsole = LogInit(false, zap.DebugLevel, "")
}

func SLCDebug(format string, args ...any) {
	SLoggerConsole.Debug(fmt.Sprintf(format, args))
}

func SLCInfo(format string, args ...any) {
	SLoggerConsole.Info(fmt.Sprintf(format, args))
}

func SLCWarn(format string, args ...any) {
	SLoggerConsole.Warn(fmt.Sprintf(format, args))
}

func SLCError(format string, args ...any) {
	SLoggerConsole.Error(fmt.Sprintf(format, args))
}

func SLCFatal(format string, args ...any) {
	SLoggerConsole.Fatal(fmt.Sprintf(format, args))
}

func SLCPanic(format string, args ...any) {
	SLoggerConsole.Panic(fmt.Sprintf(format, args))
}

// LogInit isJson决定文件输出的是否为json格式, level决定输出的最小等级,
// filePath为路径名比如/log/test,最终为/log/test日期.log
// filePath不输出到文件直接传空字符串
func LogInit(isJson bool, level zapcore.Level, filePath string) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	var fileEncoder zapcore.Encoder
	if isJson {
		fileEncoder = zapcore.NewJSONEncoder(pe)
	} else {
		fileEncoder = zapcore.NewConsoleEncoder(pe)
	}
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	var core zapcore.Core
	if len(filePath) > 0 {
		file, _ := os.OpenFile(filePath+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
	}
	l := zap.New(core)
	return l.Sugar()
}

//type Logger struct {
//	sugarLogger *zap.SugaredLogger
//	logger      *zap.Logger
//	path        string
//}
//
//func (s *Logger) SugarLogger() *zap.SugaredLogger {
//	return s.sugarLogger
//}
//
//func (s *Logger) Logger() *zap.Logger {
//	return s.logger
//}
//
//func (s *Logger) Path() string {
//	return s.path
//}
//
//func (s *Logger) SetPath(path string) {
//	s.path = path
//}
//
//func NewLogger(path string) *Logger {
//	if len(path) == 0 {
//		return nil
//	}
//	s := &Logger{path: path}
//	s.init()
//	return s
//}
//
//func (s *Logger) init() {
//	writeSyncer := s.getLogWriter()
//	encoder := getEncoder()
//	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
//	logger := zap.New(core)
//	s.sugarLogger = logger.Sugar()
//	s.logger = logger
//	zap.ReplaceGlobals(logger)
//}
//
//func getEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
//}
//
//func (s *Logger) getLogWriter() zapcore.WriteSyncer {
//	file, _ := os.OpenFile(s.path+time.Now().Format("2006-01-02")+".log", os.O_APPEND|os.O_WRONLY, 0644)
//	return zapcore.AddSync(file)
//}
//func (s *Logger) SugarError(format string, a ...any) {
//	s.sugarLogger.Errorf(format, a)
//}
//
//func (s *Logger) SugarInfo(format string, a ...any) {
//	s.sugarLogger.Infof(format, a)
//}
