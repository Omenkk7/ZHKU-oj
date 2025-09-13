package logger

import (
	"io"
	"os"
	"path/filepath"

	"zhku-oj/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

// Init 初始化日志
func Init(cfg config.LoggingConfig) {
	log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// 设置输出
	var writers []io.Writer

	// 标准输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if cfg.File.Enabled {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.File.Path)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.WithError(err).Error("创建日志目录失败")
		} else {
			fileWriter := &lumberjack.Logger{
				Filename:   cfg.File.Path,
				MaxSize:    cfg.File.MaxSize,
				MaxBackups: cfg.File.MaxBackups,
				MaxAge:     cfg.File.MaxAge,
				LocalTime:  true,
				Compress:   true,
			}
			writers = append(writers, fileWriter)
		}
	}

	if len(writers) > 1 {
		log.SetOutput(io.MultiWriter(writers...))
	} else if len(writers) == 1 {
		log.SetOutput(writers[0])
	}
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
	}
	return log
}

// Debug 调试日志
func Debug(msg string, fields ...interface{}) {
	entry := log.WithFields(parseFields(fields...))
	entry.Debug(msg)
}

// Info 信息日志
func Info(msg string, fields ...interface{}) {
	entry := log.WithFields(parseFields(fields...))
	entry.Info(msg)
}

// Warn 警告日志
func Warn(msg string, fields ...interface{}) {
	entry := log.WithFields(parseFields(fields...))
	entry.Warn(msg)
}

// Error 错误日志
func Error(msg string, fields ...interface{}) {
	entry := log.WithFields(parseFields(fields...))
	entry.Error(msg)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...interface{}) {
	entry := log.WithFields(parseFields(fields...))
	entry.Fatal(msg)
}

// parseFields 解析字段参数
func parseFields(fields ...interface{}) logrus.Fields {
	result := logrus.Fields{}
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}
	return result
}
