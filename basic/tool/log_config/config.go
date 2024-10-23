package log_config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

type LogConfig struct {
	//文件路径，精确到文件名,文件路径所涉及文件夹需提前创建 例:/home/tails/logs/all.log
	Filename string
	//单个日志文件大小单位MB
	MaxSize int `default:50`
	//已过期文件最多保留数量
	MaxBackups int `default:10`
	//日志保留天数
	MaxAge int `default:90`
	//是否需要压缩滚动日志, 使用的 gzip 压缩，缺省为 false。
	Compress bool `default:"true"`
	//日志级别
	Level log.Level `default:log.InfoLevel`
	//0仅控制台输出 1仅日志文件输出 2控制台和日志文件输出
	OutType int `default:1`
}

/*
*
自定义日志格式化
*/
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	fName := filepath.Base(entry.Caller.File)
	//日期 ; 日志级别 ; 包名 行数 方法名 日志内容
	return []byte(fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n", timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)), nil
}

func (cfg *LogConfig) Init() {
	log.SetReportCaller(true)
	log.SetFormatter(&CustomFormatter{})
	log.SetLevel(cfg.Level)
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
