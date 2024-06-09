package logger

import (
	"fmt"
	"github.com/NumberMan1/common/global/variable"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var SLoggerConsole *zap.SugaredLogger
var level zapcore.Level

func init() {
	SLoggerConsole = LogInit(false, zap.DebugLevel, "")
}

func SLCDebug(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Debug(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Debug(format)
	}
}

func SLCInfo(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Info(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Info(format)
	}
}

func SLCWarn(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Warn(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Warn(format)
	}
}

func SLCError(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Error(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Error(format)
	}
}

func SLCFatal(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Fatal(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Fatal(format)
	}
}

func SLCPanic(format string, args ...any) {
	if len(args) > 0 {
		SLoggerConsole.Panic(fmt.Sprintf(format, args...))
	} else {
		SLoggerConsole.Panic(format)
	}
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
	l := zap.New(core, zap.AddStacktrace(level))
	return l.Sugar()
}

func Zap() (logger *zap.Logger) {
	if ok, _ := PathExists(variable.Config.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", variable.Config.Zap.Director)
		_ = os.Mkdir(variable.Config.Zap.Director, os.ModePerm)
	}

	switch variable.Config.Zap.LogLevel { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore())
	}
	logger = logger.WithOptions(zap.AddCaller())
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "protocol",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case variable.Config.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case variable.Config.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case variable.Config.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case variable.Config.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(variable.Config.Zap.LogPrefix + "2006-01-02 15:04:05.000"))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(variable.Config.Zap.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err

}
