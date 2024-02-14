package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFileName = "./log/%s.log"
	// LogFileMaxSize 每个日志文件最大 MB
	logFileMaxSize = 512
	// LogFileMaxBackups 保留日志文件个数
	logFileMaxBackups = 10
	// LogFileMaxAge 保留日志最大天数
	logFileMaxAge = 14

	envLogLevel          = "log_level"
	envLogFileOutput     = "log_file_output"
	envLogOutputFileName = "log_output_filename"
)

func newOutput() io.Writer {
	// fileOutput := os.Getenv(envLogFileOutput)
	// if fileOutput == "true" {
	// 	return io.MultiWriter(
	// 		os.Stdout,
	// 		&lumberjack.Logger{
	// 			Filename:   fmt.Sprintf(logFileName, os.Getenv(envLogOutputFileName)),
	// 			MaxSize:    logFileMaxSize,
	// 			MaxAge:     logFileMaxAge,
	// 			MaxBackups: logFileMaxBackups,
	// 			LocalTime:  true,
	// 			Compress:   false,
	// 		},
	// 	)
	// }

	// return os.Stdout

	return io.MultiWriter(
		&lumberjack.Logger{
			Filename:   fmt.Sprintf(logFileName, os.Getenv(envLogOutputFileName)),
			MaxSize:    logFileMaxSize,
			MaxAge:     logFileMaxAge,
			MaxBackups: logFileMaxBackups,
			LocalTime:  true,
			Compress:   false,
		},
	)
}

func newLevel() hlog.Level {
	switch os.Getenv(envLogLevel) {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	}

	return hlog.LevelTrace
}
