package log_config

import (
	"basic"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	//1仅日志文件输出 2控制台和日志文件输出 3仅控制台输出
	OutType int
}

/*
*
创建默认日志配置对象
*/
func NewLogConfig() *logConfig {
	return &logConfig{"./logs/all.log", 50, 10, 90, true, "info", 1}
}

/*
*
自定义日志格式化
*/
type CustomFormatter struct{}

/*
*
自定义日志格式化
*/
func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	fName := filepath.Base(entry.Caller.File)
	gid, _ := basic.GetCache(basic.GID)
	if gid != nil {
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
	log.SetFormatter(&CustomFormatter{})
	log.SetLevel(logLevel(cfg.Level))
	logger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	switch cfg.OutType {
	case 0:
		log.SetOutput(os.Stdout)
	case 1:
		log.SetOutput(logger)
	default:
		log.SetOutput(io.MultiWriter(logger, os.Stdout))
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
