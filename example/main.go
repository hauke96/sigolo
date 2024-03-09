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
	sigolo.Info("Hello world!")
	sigolo.Infof("Hello world! With formatting %d", 123)
	sigolo.Debugf("Hello world %d times!", 42) // not shown because log-level is on INFO
	sigolo.Errorf("Hello world!")

	sigolo.SetDefaultLogLevel(sigolo.LOG_PLAIN)
	sigolo.Plain("Some plain text")                          // now the log level is ok
	sigolo.Plainf("Some plain text with formatting %d", 123) // now the log level is ok

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 1 =====\n")
	sigolo.SetDefaultLogLevel(sigolo.LOG_DEBUG)

	sigolo.Info("Hello world!")
	sigolo.Infof("Hello formatted %s!", "world")
	sigolo.Infob(0, "Backward hello formatted %s!", "world")
	sigolo.Debugf("Hello world %d times!", 42) // shown because log-level is on DEBUG
	sigolo.Errorf("Hello %x!", 123)

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== 2 =====\n")
	sigolo.SetDefaultFormatFunction(sigolo.LOG_INFO, simpleInfo)

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
	sigolo.SetDefaultDateFormat("02.01.2006 at 15:04:05")

	sigolo.Infof("Hello world!")
	sigolo.Debugf("Hello world!")
	sigolo.Errorf("Hello world!")

	time.Sleep(time.Millisecond)
	fmt.Println("\n===== Logger struct - default =====\n")
	logger := sigolo.NewLoggerf(sigolo.LOG_INFO, sigolo.LogDefault)
	logger.Info("Normal info")
	logger.Infof("Formatted info %d", 123)
	logger.Infob(0, "Backward info %d", 123)
	logger.Debugf("Not visible %d", 123)

	fmt.Println("\n===== FatalCheck =====\n")
	sigolo.FatalCheck(thisFunc())
}

func simpleInfo(writer *os.File, time, level string, maxLength int, caller string, traceId int, message string) {
	fmt.Fprintf(writer, ">>  My custom Infof  ||  %s\n", message)
}
