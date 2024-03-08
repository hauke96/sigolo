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
	sigolo.Plainf("Hello world!") // will not be printed as log level is to restrictive
	sigolo.Infof("Hello world!")
	sigolo.Debugf("Hello world %d times!", 42) // not shown because log-level is on INFO
	sigolo.Errorf("Hello world!")

	sigolo.LogLevel = sigolo.LOG_PLAIN
	sigolo.Plainf("Some plain text") // now the log level is ok

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 1 =====\n")
	sigolo.LogLevel = sigolo.LOG_DEBUG

	sigolo.Infof("Hello %s!", "world")
	sigolo.Debugf("Hello world %d times!", 42) // shown because log-level is on DEBUG
	sigolo.Errorf("Hello %x!", 123)

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 2 =====\n")
	sigolo.FormatFunctions[sigolo.LOG_INFO] = simpleInfo

	sigolo.Infof("Some")
	sigolo.Infof("AMAZING")
	sigolo.Infof("log")
	sigolo.Infof("entries")
	sigolo.Debugf("Boring")
	sigolo.Errorf("Lame")

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 3 =====\n")
	// This will print stack trace information because "thatFunc()" added them:
	sigolo.Stack(thisFunc())
	// This doesn't and just prints the error message:
	sigolo.Stack(someFrameworkFunction())

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 4 =====\n")
	sigolo.DateFormat = "02.01.2006 at 15:04:05"

	sigolo.Infof("Hello world!")
	sigolo.Debugf("Hello world!")
	sigolo.Errorf("Hello world!")

	sigolo.FatalCheck(thisFunc())
}

func simpleInfo(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, ">>  My custom Infof  ||  %s\n", message)
}
