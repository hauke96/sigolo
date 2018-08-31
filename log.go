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
	LOG_DEBUG Level = iota
	LOG_INFO
	LOG_ERROR
	LOG_FATAL
)

var (
	LogLevel   Level  = LOG_INFO
	DateFormat string = "2006-01-02 15:04:05.000"

	FormatFunctions map[Level]func(*os.File, string, string, int, string, string) = map[Level]func(*os.File, string, string, int, string, string){
		LOG_DEBUG: LogDefault,
		LOG_INFO:  LogDefault,
		LOG_ERROR: LogDefault,
		LOG_FATAL: LogDefault,
	}

	// The current maximum length printed for caller information. This is updated each time something gets printed
	CallerColumnWidth = 0

	levelStrings map[Level]string = map[Level]string{
		LOG_DEBUG: "[DEBUG]",
		LOG_INFO:  "[INFO] ",
		LOG_ERROR: "[ERROR]",
		LOG_FATAL: "[FATAL]",
	}

	levelOutputs map[Level]*os.File = map[Level]*os.File{
		LOG_DEBUG: os.Stdout,
		LOG_INFO:  os.Stdout,
		LOG_ERROR: os.Stderr,
		LOG_FATAL: os.Stderr,
	}
)

func Info(message string) {
	log(LOG_INFO, message)
}

func Debug(message string) {
	log(LOG_DEBUG, message)
}

func Error(message string) {
	log(LOG_ERROR, message)
}

func Fatal(message string) {
	log(LOG_FATAL, message)
	os.Exit(1)
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
