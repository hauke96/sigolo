package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"

	"github.com/hauke96/sigolo"
)

func someFrameworkFunction() error {
	// This function simulated a library/framework throwing an error but not using the "errors" package. Therefor this
	// pure go error doesn't contain any stack trace information.
	return fmt.Errorf("BOOM!!! some error occurred maybe from within an framework?")
}

func thatFunc() error {
	// The errors package will add stack trace information which is later used by sigolo to print that stack trace
	return errors.Wrap(someFrameworkFunction(), "that func wrapped this error")
}

func thisFunc() error {
	return thatFunc()
}

func main() {
	sigolo.Info("Hello world!")
	sigolo.Debug("Hello world %d times!", 42) // not shown because log-level is on INFO
	sigolo.Error("Hello world!")

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 1 =====\n")
	sigolo.LogLevel = sigolo.LOG_DEBUG

	sigolo.Info("Hello %s!", "world")
	sigolo.Debug("Hello world %d times!", 42) // shown because log-level is on DEBUG
	sigolo.Error("Hello %x!", 123)

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 2 =====\n")
	sigolo.FormatFunctions[sigolo.LOG_INFO] = simpleInfo

	sigolo.Info("Some")
	sigolo.Info("AMAZING")
	sigolo.Info("log")
	sigolo.Info("entries")
	sigolo.Debug("Boring")
	sigolo.Error("Lame")

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 3 =====\n")
	// This will print stack trace information because "thatFunc()" added them:
	sigolo.Stack(thisFunc())
	// This doesn't and just prints the error message:
	sigolo.Stack(someFrameworkFunction())

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 4 =====\n")
	sigolo.DateFormat = "02.01.2006 at 15:04:05"

	sigolo.Info("Hello world!")
	sigolo.Debug("Hello world!")
	sigolo.Error("Hello world!")

	sigolo.FatalCheck(thisFunc())
}

func simpleInfo(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, ">>  My custom Info  ||  %s\n", message)
}
