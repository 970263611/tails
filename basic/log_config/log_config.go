package log_config

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
*
日志支持的配置信息
*/
type logConfig struct {
	//文件路径，精确到文件名,文件路径所涉及文件夹需提前创建 例:/home/tails/logs/all.log
	Filename string
	//单个日志文件大小单位MB
	MaxSize int
	//已过期文件最多保留数量
	MaxBackups int
	//日志保留天数
	MaxAge int
	//是否需要压缩滚动日志, 使用的 gzip 压缩，缺省为 false。
	Compress bool
	//日志级别
	Level string
	//全局对象
	iface.Context
}

/*
*
创建默认日志配置对象
*/
func NewLogConfig(c iface.Context) *logConfig {
	return &logConfig{utils.GetRootPath() + "logs", 50, 10, 90, true, "info", c}
}

/*
*
自定义日志格式化
*/
type CustomFormatter struct {
	iface.Context
}

/*
*
自定义日志格式化
*/
func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	fName := filepath.Base(entry.Caller.File)
	gid, _ := f.FindSystemParams(cons.SYSPARAM_GID)
	if gid != "" {
		return []byte(fmt.Sprintf("[%s] [%s] [%s] [%s:%d %s] %s\n", timestamp, entry.Level, fmt.Sprintf("%v", gid), fName, entry.Caller.Line, entry.Caller.Function, entry.Message)), nil
	} else {
		return []byte(fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n", timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)), nil
	}
}

/*
*
全系统级别日志配置初始化
*/
func Init(cfg *logConfig) {
	log.SetReportCaller(true)
	log.SetFormatter(&CustomFormatter{cfg.Context})
	log.SetLevel(logLevel(cfg.Level))
	log.SetLevel(log.DebugLevel)
	allWriter := createLogWriter(cfg, "all")
	debugWriter := createLogWriter(cfg, "debug")
	infoWriter := createLogWriter(cfg, "info")
	warnWriter := createLogWriter(cfg, "warn")
	errorWriter := createLogWriter(cfg, "error")
	log.AddHook(&writer.Hook{allWriter, log.AllLevels})
	log.AddHook(&writer.Hook{debugWriter, []log.Level{log.TraceLevel, log.DebugLevel}})
	log.AddHook(&writer.Hook{infoWriter, []log.Level{log.InfoLevel}})
	log.AddHook(&writer.Hook{warnWriter, []log.Level{log.WarnLevel}})
	log.AddHook(&writer.Hook{errorWriter, []log.Level{log.ErrorLevel, log.FatalLevel, log.PanicLevel}})
	if _, ok := cfg.FindSystemParams(cons.SYSPARAM_LOG_CONSOLE); ok {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(io.Discard)
	}
}

func createLogWriter(cfg *logConfig, level string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.Filename + "/" + strings.ToLower(level) + ".log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
}

/*
*
日志级别翻译，将字符串类型转换为标准枚举类型
*/
func logLevel(level string) (l log.Level) {
	switch strings.ToUpper(level) {
	case "PANIC":
		l = log.PanicLevel
	case "FATAL":
		l = log.FatalLevel
	case "ERROR":
		l = log.ErrorLevel
	case "WARN":
		l = log.WarnLevel
	case "DEBUG":
		l = log.DebugLevel
	case "TRACE":
		l = log.TraceLevel
	default:
		l = log.InfoLevel
	}
	return
}
