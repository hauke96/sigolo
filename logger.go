package sigolo

import (
	"fmt"
	"io"
	"time"
)

type Logger struct {
	LogTraceId      int
	LogLevel        Level
	DateFormat      string
	FormatFunctions map[Level]func(io.Writer, string, string, int, string, int, string)
	LevelStrings    map[Level]string
	LevelOutputs    map[Level]io.Writer
}

func NewLogger() *Logger {
	traceId := nextTraceId
	nextTraceId++
	return &Logger{
		LogTraceId:      traceId,
		LogLevel:        logLevel,
		DateFormat:      dateFormat,
		FormatFunctions: DefaultLogFormatFunctions(),
		LevelStrings:    DefaultLevelStrings(),
		LevelOutputs:    DefaultLevelOutputs(),
	}
}

func NewLoggerl(logLevel Level) *Logger {
	traceId := nextTraceId
	nextTraceId++
	return &Logger{
		LogTraceId:      traceId,
		LogLevel:        logLevel,
		DateFormat:      dateFormat,
		FormatFunctions: DefaultLogFormatFunctions(),
		LevelStrings:    DefaultLevelStrings(),
		LevelOutputs:    DefaultLevelOutputs(),
	}
}

func NewLoggerf(logLevel Level, defaultFormat func(io.Writer, string, string, int, string, int, string)) *Logger {
	traceId := nextTraceId
	nextTraceId++

	formatFunctions := map[Level]func(io.Writer, string, string, int, string, int, string){
		LOG_PLAIN: defaultFormat,
		LOG_TRACE: defaultFormat,
		LOG_DEBUG: defaultFormat,
		LOG_INFO:  defaultFormat,
		LOG_ERROR: defaultFormat,
		LOG_FATAL: defaultFormat,
	}

	return &Logger{
		LogTraceId:      traceId,
		LogLevel:        logLevel,
		DateFormat:      dateFormat,
		FormatFunctions: formatFunctions,
		LevelStrings:    DefaultLevelStrings(),
		LevelOutputs:    DefaultLevelOutputs(),
	}
}

func (l *Logger) Plain(message string) {
	l.Plainb(1, "%s", message)
}

func (l *Logger) Plainf(format string, args ...interface{}) {
	l.Plainb(1, fmt.Sprintf(format, args...))
}

// Plainb is equal to Plainf(...) but can go back in the stack and can therefore show function positions from previous functions.
func (l *Logger) Plainb(framesBackward int, format string, args ...interface{}) {
	if l.LogLevel > LOG_PLAIN {
		return
	}
	l.log(LOG_PLAIN, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Trace(message string) {
	l.Traceb(1, "%s", message)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Traceb(1, fmt.Sprintf(format, args...))
}

// Traceb is equal to Tracef(...) but can go back in the stack and can therefore show function positions from previous functions.
func (l *Logger) Traceb(framesBackward int, format string, args ...interface{}) {
	if l.LogLevel > LOG_TRACE {
		return
	}
	l.log(LOG_TRACE, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(message string) {
	l.Debugb(1, "%s", message)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debugb(1, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugb(framesBackward int, format string, args ...interface{}) {
	if logLevel > LOG_DEBUG {
		return
	}
	l.log(LOG_DEBUG, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(message string) {
	l.Infob(1, "%s", message)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Infob(1, fmt.Sprintf(format, args...))
}

func (l *Logger) Infob(framesBackward int, format string, args ...interface{}) {
	if logLevel > LOG_INFO {
		return
	}
	l.log(LOG_INFO, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(message string) {
	l.Warnb(1, "%s", message)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warnb(1, fmt.Sprintf(format, args...))
}

func (l *Logger) Warnb(framesBackward int, format string, args ...interface{}) {
	if logLevel > LOG_WARN {
		return
	}
	l.log(LOG_WARN, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(message string) {
	l.Errorb(1, "%s", message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Errorb(1, fmt.Sprintf(format, args...))
}

func (l *Logger) Errorb(framesBackward int, format string, args ...interface{}) {
	if logLevel > LOG_ERROR {
		return
	}
	l.log(LOG_ERROR, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(message string) {
	l.Fatalb(1, "%s", message)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatalb(1, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalb(framesBackward int, format string, args ...interface{}) {
	if logLevel > LOG_FATAL {
		return
	}
	l.log(LOG_FATAL, 3+framesBackward, l.LogTraceId, fmt.Sprintf(format, args...))
}

// Stack tries to print the stack trace of the given error using the  %+v  format string. When using the
// https://github.com/pkg/errors package, this will print a full stack trace of the error. If normal errors are used,
// this function will just print the error.
func (l *Logger) Stack(err error) {
	if logLevel > LOG_ERROR {
		return
	}
	// Directly call "log" to avoid extra function call
	l.log(LOG_ERROR, 3, l.LogTraceId, fmt.Sprintf("%+v", err))
}

// Stackb is equal to Stack(...) but can go back in the stack and can therefore show function positions from previous functions.
func (l *Logger) Stackb(framesBackward int, err error) {
	if logLevel > LOG_ERROR {
		return
	}
	// Directly call "log" to avoid extra function call
	l.log(LOG_ERROR, 3+framesBackward, l.LogTraceId, fmt.Sprintf("%+v", err))
}

func (l *Logger) log(level Level, framesBackward int, traceId int, message string) {
	// A bit hacky: We know here that the stack contains two calls from inside
	// this file. The third frame comes from the file that initially called a
	// function in this file (e.g. Infof())
	caller := GetCallerDetails(framesBackward)

	updateCallerColumnWidth(caller)

	l.FormatFunctions[level](l.LevelOutputs[level], time.Now().Format(l.DateFormat), l.LevelStrings[level], CallerColumnWidth, caller, traceId, message)
}
