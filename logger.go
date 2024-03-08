package sigolo

import (
	"fmt"
)

var (
	nextTraceId          = 0
	defaultLoggingFormat = "#%x | %s"
)

type Logger struct {
	LogTraceId int
	Format     string
}

func NewLogger() *Logger {
	traceId := nextTraceId
	nextTraceId++
	return &Logger{
		LogTraceId: traceId,
		Format:     defaultLoggingFormat,
	}
}

func NewLoggerf(format string) *Logger {
	traceId := nextTraceId
	nextTraceId++
	return &Logger{
		LogTraceId: traceId,
		Format:     format,
	}
}

func (l *Logger) Debug(message string) {
	Debugb(1, defaultLoggingFormat, l.LogTraceId, message)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	Debugb(1, defaultLoggingFormat, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(message string) {
	Infob(1, defaultLoggingFormat, l.LogTraceId, message)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	Infob(1, defaultLoggingFormat, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(message string) {
	Warnb(1, defaultLoggingFormat, l.LogTraceId, message)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	Warnb(1, defaultLoggingFormat, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(message string) {
	Errorb(1, defaultLoggingFormat, l.LogTraceId, message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	Errorb(1, defaultLoggingFormat, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(message string) {
	Fatalb(1, defaultLoggingFormat, l.LogTraceId, message)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	Fatalb(1, defaultLoggingFormat, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Stack(err error) {
	Stackb(1, err)
}
