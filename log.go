package sigolo

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type Level int

const (
	LOG_PLAIN Level = iota
	LOG_TRACE
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
)

var (
	nextTraceId = 0
	logLevel    = LOG_INFO
	dateFormat  = "2006-01-02 15:04:05.000"

	// The current maximum length printed for caller information. This is updated each time something gets printed
	CallerColumnWidth = 0

	formatFunctions = DefaultStaticLogFormatFunctions()
	levelStrings    = DefaultLevelStrings()
	levelOutputs    = DefaultLevelOutputs()

	DefaultLogger = GetLoggerWithCurrentDefaults()
)

func DefaultLogFormatFunctions() map[Level]func(*os.File, string, string, int, string, int, string) {
	return map[Level]func(*os.File, string, string, int, string, int, string){
		LOG_PLAIN: LogPlain,
		LOG_TRACE: LogDefault,
		LOG_DEBUG: LogDefault,
		LOG_INFO:  LogDefault,
		LOG_ERROR: LogDefault,
		LOG_FATAL: LogDefault,
	}
}

func DefaultStaticLogFormatFunctions() map[Level]func(*os.File, string, string, int, string, int, string) {
	return map[Level]func(*os.File, string, string, int, string, int, string){
		LOG_PLAIN: LogPlain,
		LOG_TRACE: LogDefaultStatic,
		LOG_DEBUG: LogDefaultStatic,
		LOG_INFO:  LogDefaultStatic,
		LOG_ERROR: LogDefaultStatic,
		LOG_FATAL: LogDefaultStatic,
	}
}

func DefaultLevelStrings() map[Level]string {
	return map[Level]string{
		LOG_PLAIN: "",
		LOG_TRACE: "[TRACE]",
		LOG_DEBUG: "[DEBUG]",
		LOG_INFO:  "[INFO] ",
		LOG_ERROR: "[ERROR]",
		LOG_FATAL: "[FATAL]",
	}
}

func DefaultLevelOutputs() map[Level]*os.File {
	return map[Level]*os.File{
		LOG_PLAIN: os.Stdout,
		LOG_TRACE: os.Stdout,
		LOG_DEBUG: os.Stdout,
		LOG_INFO:  os.Stdout,
		LOG_ERROR: os.Stderr,
		LOG_FATAL: os.Stderr,
	}
}

func GetCurrentLogLevel() Level {
	return logLevel
}

func GetCurrentDateFormat() string {
	return dateFormat
}

func GetCurrentNextTraceId() int {
	return nextTraceId
}

func GetLoggerWithCurrentDefaults() *Logger {
	return &Logger{
		LogTraceId:      nextTraceId,
		LogLevel:        logLevel,
		DateFormat:      dateFormat,
		FormatFunctions: formatFunctions,
		LevelStrings:    levelStrings,
		LevelOutputs:    levelOutputs,
	}
}

func SetDefaultDateFormat(format string) {
	dateFormat = format
	DefaultLogger = GetLoggerWithCurrentDefaults()
}

func SetDefaultLogLevel(level Level) {
	logLevel = level
	DefaultLogger = GetLoggerWithCurrentDefaults()
}

func SetDefaultFormatFunction(level Level, function func(*os.File, string, string, int, string, int, string)) {
	formatFunctions[level] = function
	DefaultLogger = GetLoggerWithCurrentDefaults()
}

func SetDefaultLevelString(level Level, output *os.File) {
	levelOutputs[level] = output
	DefaultLogger = GetLoggerWithCurrentDefaults()
}

func SetDefaultLevelOutput(level Level, prefix string) {
	levelStrings[level] = prefix
	DefaultLogger = GetLoggerWithCurrentDefaults()
}

func Plain(message string) {
	DefaultLogger.Plainb(1, "%s", message)
	increaseTraceId()
}

func Plainf(format string, args ...interface{}) {
	DefaultLogger.Plainb(1, format, args...)
	increaseTraceId()
}

// Plainb is equal to Plainf(...) but can go back in the stack and can therefore show function positions from previous functions.
func Plainb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Plainb(1+framesBackward, format, args...)
	increaseTraceId()
}

func Trace(message string) {
	DefaultLogger.Traceb(1, "%s", message)
	increaseTraceId()
}

func Tracef(format string, args ...interface{}) {
	DefaultLogger.Traceb(1, format, args...)
	increaseTraceId()
}

// Traceb is equal to Tracef(...) but can go back in the stack and can therefore show function positions from previous functions.
func Traceb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Traceb(1+framesBackward, format, args...)
	increaseTraceId()
}

func Debug(message string) {
	DefaultLogger.Debugb(1, "%s", message)
	increaseTraceId()
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugb(1, format, args...)
	increaseTraceId()
}

// Debugb is equal to Debugf(...) but can go back in the stack and can therefore show function positions from previous functions.
func Debugb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Debugb(1+framesBackward, format, args...)
	increaseTraceId()
}

func Info(message string) {
	DefaultLogger.Infob(1, "%s", message)
	increaseTraceId()
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infob(1, format, args...)
	increaseTraceId()
}

// Infob is equal to Infof(...) but can go back in the stack and can therefore show function positions from previous functions.
func Infob(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Infob(1+framesBackward, format, args...)
	increaseTraceId()
}

func Warn(message string) {
	DefaultLogger.Warnb(1, "%s", message)
	increaseTraceId()
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnb(1, format, args...)
	increaseTraceId()
}

// Warnb is equal to Warnf(...) but can go back in the stack and can therefore show function positions from previous functions.
func Warnb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Warnb(1+framesBackward, format, args...)
	increaseTraceId()
}

func Error(message string) {
	DefaultLogger.Errorb(1, "%s", message)
	increaseTraceId()
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorb(1, format, args...)
	increaseTraceId()
}

// Errorb is equal to Errorf(...) but can go back in the stack and can therefore show function positions from previous functions.
func Errorb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Errorb(1+framesBackward, format, args...)
	increaseTraceId()
}

func Fatal(message string) {
	DefaultLogger.Fatalb(1, "%s", message)
	increaseTraceId()
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalb(1, format, args...)
	increaseTraceId()
	os.Exit(1)
}

// Fatalb is equal to Fatalf(...) but can go back in the stack and can therefore show function positions from previous functions.
func Fatalb(framesBackward int, format string, args ...interface{}) {
	DefaultLogger.Fatalb(1+framesBackward, format, args...)
	increaseTraceId()
}

// Stack tries to print the stack trace of the given error using the  %+v  format string. When using the
// https://github.com/pkg/errors package, this will print a full stack trace of the error. If normal errors are used,
// this function will just print the error.
func Stack(err error) {
	DefaultLogger.Stackb(1, err)
	increaseTraceId()
}

// Stackb is equal to Stack(...) but can go back in the stack and can therefore show function positions from previous functions.
func Stackb(framesBackward int, err error) {
	DefaultLogger.Stackb(1+framesBackward, err)
	increaseTraceId()
}

// FatalCheckf checks if the error exists (!= nil). If so, it'll print the error
// message and fatals with the given format message.
func FatalCheckf(err error, traceId int, format string, args ...interface{}) {
	if err != nil {
		Stackb(1, err)
		if args != nil {
			internalFatalf(traceId, format, args...)
		} else {
			internalFatalf(traceId, format)
		}
	}
}

// FatalCheck checks if the error exists (!= nil). If so, it'll fatal with the error message.
func FatalCheck(err error) {
	if err != nil {
		Stackb(1, err)
		os.Exit(1)
	}
}

func internalFatalf(traceId int, format string, args ...interface{}) {
	internalLog(LOG_FATAL, traceId, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func internalLog(level Level, traceId int, message string) {
	// A bit hacky: We know here that the stack contains three calls from inside
	// this file. The third frame comes from the file that initially called a
	// function in this file (e.g. Infof())
	caller := getCallerDetails(4)

	updateCallerColumnWidth(caller)

	formatFunctions[level](levelOutputs[level], time.Now().Format(dateFormat), levelStrings[level], CallerColumnWidth, caller, traceId, message)
}

func updateCallerColumnWidth(caller string) {
	if len(caller) > CallerColumnWidth {
		CallerColumnWidth = len(caller)
	}
}

func getCallerDetails(framesBackwards int) string {
	name := "???"
	line := -1
	ok := false

	_, name, line, ok = runtime.Caller(framesBackwards)

	if ok {
		name = path.Base(name)
	}

	caller := fmt.Sprintf("%s:%d", name, line)

	return caller
}

// increaseTraceId increases the trace ID of the default logger. This has the effect, that the caller doesn't know that
// in the background the same DefaultLogger instance is "recycled".
func increaseTraceId() {
	DefaultLogger.LogTraceId++
	nextTraceId++
}

func LogDefault(writer *os.File, time, level string, maxLength int, caller string, traceId int, message string) {
	fmt.Fprintf(writer, "%s %s %-*s | #%x | %s\n", time, level, maxLength, caller, traceId, message)
}

func LogDefaultStatic(writer *os.File, time, level string, maxLength int, caller string, traceId int, message string) {
	fmt.Fprintf(writer, "%s %s %-*s | #%x | %s\n", time, level, maxLength, caller, traceId, message)
}

func LogPlain(writer *os.File, time, level string, maxLength int, caller string, traceId int, message string) {
	fmt.Fprintf(writer, "%s\n", message)
}
