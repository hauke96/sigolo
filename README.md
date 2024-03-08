# sigolo

**Si**mple **go**lang **lo**gging helper with lots of ways to customize the logging output.

# How to use it

tl;dr Just call `sigolo.{Plain|Debug|Info|Warn|Error|Fatal}` with a message.

## Simple calls

```go
sigolo.Info("Hello world!")
sigolo.Debugf("Coordinate: %d, %d", x, y)
```

The default printing format is something like this:

```bash
2018-07-21 01:59:05.431 [INFO]  main.go:21 | Hello world!
2018-07-21 01:59:05.432 [DEBUG] main.go:22 | Coordinate: 42, 13
```

Only the `sigolo.Plainf` function does not produce leading information (date, log-level, etc.) and just acts like `fmt.Printf` does.

## Error handling

I recommend the [pkg/errors](https://github.com/pkg/errors) package to create and wrap your errors.
Why? It enables you to see stack traces ;)

To exit on an error, there's the `sigolo.FatalCheck` function:

```go
err := someFunctionCall()
sigolo.FatalCheck(err)
```

When `err` is *not* `nil`, then the error including stack trace will be printed and your application exists with exit code 1.

## Log level

Specify the log level by changing `sigolo.LogLevel`.
Possible value are `sigolo.LOG_PLAIN`, `sigolo.LOG_DEBUG`, `sigolo.LOG_INFO`, `sigolo.LOG_WARN`, `sigolo.LOG_ERROR` and `sigolo.LOG_FATAL`.
The levels are ordered, choosing one "mutes" the previous ones.
Example: When choosing `sigolo.LOG_INFO` then the plain and debug relates method do not print anything.

The fatal methods print without any formatting on stderr and then exit with `os.Exit(1)`.

## Function suffixes / Variants

Some functions have a suffix with slightly different behavior.

* `b`: Acts like the normal function, but in order to print the correct caller you can go **b**ack in the stack by a given number of frames.
* `f`: Acts like `fmt.Printf`.

## Change general output format

The format can be changed by implementing the printing function specified in the `sigolo.FormatFunctions` array.

Exmaple: To specify your own debug-format:

```go
func main() {
    // Whenever sigolo.Debug is called, our simpleDebug method is used to produce the output.
    sigolo.FormatFunctions[sigolo.LOG_DEBUG] = simpleDebug
    
    sigolo.Debug("Hello world!")
    }
    
    func simpleDebug(writer *os.File, time, level string, maxLength int, caller, message string) {
    // Don't forget the \n at the end ;)
    fmt.Fprintf(writer, "Debug: %s\n", message)
}
```

This example will print:

```bash
Debug: Hello world!
```

## Change time format

To change only the time format, change the value of the `sigolo.DateFormat` variable. The format of this variable if the
format described in the [time package](https://golang.org/pkg/time/).

Example:

```go
func main() {
    // Use the go time formatting string as described in https://pkg.go.dev/time
    sigolo.DateFormat = "02.01.2006 at 15:04:05"
    
    sigolo.Debug("Hello world!")
}
```

This will produce:

```bash
21.07.2018 at 02:16:41 [DEBUG] main.go:37 | Hello world!
```
