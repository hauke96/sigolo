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
	LOG_DEBUG
	LOG_INFO
	LOG_ERROR
	LOG_FATAL
)

var (
	LogLevel   Level  = LOG_INFO
	DateFormat string = "2006-01-02 15:04:05.000"

	FormatFunctions map[Level]func(*os.File, string, string, int, string, string) = map[Level]func(*os.File, string, string, int, string, string){
		LOG_PLAIN: LogPlain,
		LOG_DEBUG: LogDefault,
		LOG_INFO:  LogDefault,
		LOG_ERROR: LogDefault,
		LOG_FATAL: LogDefault,
	}

	// The current maximum length printed for caller information. This is updated each time something gets printed
	CallerColumnWidth = 0

	levelStrings map[Level]string = map[Level]string{
		LOG_PLAIN: "",
		LOG_DEBUG: "[DEBUG]",
		LOG_INFO:  "[INFO] ",
		LOG_ERROR: "[ERROR]",
		LOG_FATAL: "[FATAL]",
	}

	levelOutputs map[Level]*os.File = map[Level]*os.File{
		LOG_PLAIN: os.Stdout,
		LOG_DEBUG: os.Stdout,
		LOG_INFO:  os.Stdout,
		LOG_ERROR: os.Stderr,
		LOG_FATAL: os.Stderr,
	}
)

func Plain(format string, a ...interface{}) {
	log(LOG_PLAIN, fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
	log(LOG_INFO, fmt.Sprintf(format, a...))
}

func Debug(format string, a ...interface{}) {
	log(LOG_DEBUG, fmt.Sprintf(format, a...))
}

func Error(format string, a ...interface{}) {
	log(LOG_ERROR, fmt.Sprintf(format, a...))
}

func Fatal(format string, a ...interface{}) {
	log(LOG_FATAL, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// FatalCheckf checks if the error exists (!= nil). If so, it'll print the error
// message and fatals with the given format message.
func FatalCheckf(err error, format string, a ...interface{}) {
	if err != nil {
		Error(err.Error())
		Fatal(format, a)
	}
}

// FatalCheck checks if the error exists (!= nil). If so, it'll fatal with the error message.
func FatalCheck(err error) {
	if err != nil {
		Fatal(err.Error())
	}
}

func log(level Level, message string) {
	caller := getCallerDetails()

	updateCallerColumnWidth(caller)

	if LogLevel <= level {
		FormatFunctions[level](levelOutputs[level], time.Now().Format(DateFormat), levelStrings[level], CallerColumnWidth, caller, message)
	}
}

func updateCallerColumnWidth(caller string) {
	if len(caller) > CallerColumnWidth {
		CallerColumnWidth = len(caller)
	}
}

func getCallerDetails() string {
	name := "???"
	line := -1
	ok := false

	// A bit hacky: We know here that the stack contains two calls from inside
	// this file. The third frame comes from the file that initially called a
	// function in this file (e.g. Info())
	_, name, line, ok = runtime.Caller(3)

	if ok {
		name = path.Base(name)
	}

	caller := fmt.Sprintf("%s:%d", name, line)

	return caller
}

func LogDefault(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, "%s %s %-*s | %s\n", time, level, maxLength, caller, message)
}

func LogPlain(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, "%s\n", message)
}
