package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hauke96/sigolo"
)

func main() {
	sigolo.Info("Hello world!")
	sigolo.Debug("Hello world!")
	sigolo.Error("Hello world!")

	time.Sleep(time.Millisecond)
	fmt.Println()
	sigolo.LogLevel = sigolo.LOG_DEBUG

	sigolo.Info("Hello world!")
	sigolo.Debug("Hello world!")
	sigolo.Error("Hello world!")

	time.Sleep(time.Millisecond)
	fmt.Println()
	sigolo.FormatFunctions[sigolo.LOG_INFO] = simpleInfo

	sigolo.Info("Hello world!")
	sigolo.Debug("Hello world!")
	sigolo.Error("Hello world!")
}

func simpleInfo(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, "Info: %s\n", message)
}
