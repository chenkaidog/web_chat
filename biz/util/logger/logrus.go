package logger

import (
	"context"
	"fmt"
	"io"
	"path"
	"runtime"
	"time"
	traceinfo "web_chat/biz/util/trace_info"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/sirupsen/logrus"
)

const depth = 3

type logrusLogger struct {
	*logrus.Logger
}

func NewLogrusLogger() hlog.FullLogger {
	logger := &logrusLogger{
		Logger: logrus.New(),
	}
	logger.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		},
	)

	return logger
}

func (l *logrusLogger) entryWithLoc() *logrus.Entry {
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return l.Logger.WithFields(
			logrus.Fields{
				"location": fmt.Sprintf("%s:%d", path.Base(file), line),
			})
	}

	return l.Logger.WithFields(logrus.Fields{})
}

// CtxDebugf implements hlog.FullLogger.
func (l *logrusLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Debugf(format, v...)
}

// CtxErrorf implements hlog.FullLogger.
func (l *logrusLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Errorf(format, v...)
}

// CtxFatalf implements hlog.FullLogger.
func (l *logrusLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Fatalf(format, v...)
}

// CtxInfof implements hlog.FullLogger.
func (l *logrusLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Infof(format, v...)
}

// CtxNoticef implements hlog.FullLogger.
func (l *logrusLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Infof(format, v...)
}

// CtxTracef implements hlog.FullLogger.
func (l *logrusLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Tracef(format, v...)
}

// CtxWarnf implements hlog.FullLogger.
func (l *logrusLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.entryWithLoc().WithContext(ctx).Warnf(format, v...)
}

// Debug implements hlog.FullLogger.
func (l *logrusLogger) Debug(v ...interface{}) {
	l.entryWithLoc().Debug(v...)
}

// Debugf implements hlog.FullLogger.
func (l *logrusLogger) Debugf(format string, v ...interface{}) {
	l.entryWithLoc().Debugf(format, v...)
}

// Error implements hlog.FullLogger.
func (l *logrusLogger) Error(v ...interface{}) {
	l.entryWithLoc().Error(v...)
}

// Errorf implements hlog.FullLogger.
func (l *logrusLogger) Errorf(format string, v ...interface{}) {
	l.entryWithLoc().Errorf(format, v...)
}

// Fatal implements hlog.FullLogger.
func (l *logrusLogger) Fatal(v ...interface{}) {
	l.entryWithLoc().Fatal(v...)
}

// Fatalf implements hlog.FullLogger.
func (l *logrusLogger) Fatalf(format string, v ...interface{}) {
	l.entryWithLoc().Fatalf(format, v...)
}

// Info implements hlog.FullLogger.
func (l *logrusLogger) Info(v ...interface{}) {
	l.entryWithLoc().Info(v...)
}

// Infof implements hlog.FullLogger.
func (l *logrusLogger) Infof(format string, v ...interface{}) {
	l.entryWithLoc().Infof(format, v...)
}

// Notice implements hlog.FullLogger.
func (l *logrusLogger) Notice(v ...interface{}) {
	l.entryWithLoc().Info(v...)
}

// Noticef implements hlog.FullLogger.
func (l *logrusLogger) Noticef(format string, v ...interface{}) {
	l.entryWithLoc().Infof(format, v...)
}

// Trace implements hlog.FullLogger.
func (l *logrusLogger) Trace(v ...interface{}) {
	l.entryWithLoc().Trace(v...)
}

// Tracef implements hlog.FullLogger.
func (l *logrusLogger) Tracef(format string, v ...interface{}) {
	l.entryWithLoc().Tracef(format, v...)
}

// Warn implements hlog.FullLogger.
func (l *logrusLogger) Warn(v ...interface{}) {
	l.entryWithLoc().Warn(v...)
}

// Warnf implements hlog.FullLogger.
func (l *logrusLogger) Warnf(format string, v ...interface{}) {
	l.entryWithLoc().Warnf(format, v...)
}

func (logger *logrusLogger) SetLevel(level hlog.Level) {
	switch level {
	case hlog.LevelTrace:
		logger.Logger.SetLevel(logrus.TraceLevel)
	case hlog.LevelDebug:
		logger.Logger.SetLevel(logrus.DebugLevel)
	case hlog.LevelInfo, hlog.LevelNotice:
		logger.Logger.SetLevel(logrus.InfoLevel)
	case hlog.LevelWarn:
		logger.Logger.SetLevel(logrus.WarnLevel)
	case hlog.LevelError:
		logger.Logger.SetLevel(logrus.ErrorLevel)
	case hlog.LevelFatal:
		logger.Logger.SetLevel(logrus.FatalLevel)
	}
}

// SetOutput implements hlog.FullLogger.
func (logger *logrusLogger) SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

type hook struct{}

// Fire implements logrus.Hook.
func (*hook) Fire(entry *logrus.Entry) error {
	traceInfo := traceinfo.GetTraceInfo(entry.Context)
	entry.Data["log_id"] = traceInfo.LogID

	return nil
}

// Levels implements logrus.Hook.
func (*hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
