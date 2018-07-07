# sigolo
Simple golang logging thingy.

# What is this?
This is a simple file (I wouldn't call it _library_ or something), which helps you to log your stuff.

# How to use it
## Log level
Specify the log level by changing `sigolo.LogLevel`. Possible value are `sigolo.LOG_DEBUG`, `sigolo.LOG_INFO` and `sigolo.LOG_ERROR`.

The following things will be printed out:

| log level | method that will produce an output |
|:--:|:--|
| `LOG_DEBUG` | `sigolo.Debug()`<br>`sigolo.Info()`<br>`sigolo.Error()` |
| `LOG_INFO` | `sigolo.Info()`<br>`sigolo.Error()` |
| `LOG_ERROR` | `sigolo.Error()` |

## Simple printing
Just call `sigolo.{Info|Debug|Error}` with a message.

The default printing format is something like this:

```bash
2018-07-21 01:59:05.431 [DEBUG] main.go:21 | Hello world!
```

## Change general output format
The format can be changed by implementing the printing function specified in the `sigolo.FormatFunctions` array.

Exmaple: To specify your own debug-format:
```go
func main() {
	sigolo.FormatFunctions[sigolo.LOG_DEBUG] = simpleDebug

	sigolo.Debug("Hello world!")
}

func simpleDebug(writer *os.File, time, level string, maxLength int, caller, message string) {
	fmt.Fprintf(writer, "Debug: %s\n", message)
}
```
This will print:
```bash
Debug: Hello world!
```

## Change time format
To change only the time format, change the value of the `sigolo.DateFormat` variable. The format of this variable if the format described in the [time package](https://golang.org/pkg/time/).

Example:
```go
func main() {
	sigolo.DateFormat = "02.01.2006 at 15:04:05"

	sigolo.Debug("Hello world!")
}
```
This will produce:
```bash
21.07.2018 at 02:16:41 [DEBUG] main.go:37 | Hello world!
```
