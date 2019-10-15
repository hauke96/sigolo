package sigolo

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func prepare(logLevel Level) *os.File {
	LogLevel = LOG_PLAIN

	readPipe, writePipe, _ := os.Pipe()

	LevelOutputs[logLevel] = writePipe

	return readPipe
}

func cutOutput(f *os.File) (string, string) {
	data := make([]byte, 2<<10)
	f.Read(data)

	writtenOutput := string(data)
	writtenParts := strings.Split(writtenOutput, " ")

	writtenPayload := writtenParts[len(writtenParts)-1]
	writtenPayload = strings.Trim(writtenPayload, "\000")
	writtenPayload = strings.Trim(writtenPayload, "\n")

	return writtenParts[2], writtenPayload
}

func checkSimpleWrite(t *testing.T, pipe *os.File, originalData string, logLevel Level) {
	outputLevel, outputData := cutOutput(pipe)

	if originalData != outputData {
		t.Errorf("Payload does not match")
		t.Errorf("original : %x\n", originalData)
		t.Errorf("         : %s\n", originalData)
		t.Errorf("written  : %x\n", outputData)
		t.Errorf("         : %s\n", outputData)
		t.Fail()
	}

	// Log-level INFO has an additional space at its end because the string is shorter than others
	if logLevel == LOG_INFO {
		outputLevel += " "
	}

	if LevelStrings[logLevel] != outputLevel {
		t.Errorf("Log-level string does not patch")
		t.Errorf("original : %x\n", LevelStrings[logLevel])
		t.Errorf("         : %s\n", LevelStrings[logLevel])
		t.Errorf("written  : %x\n", outputLevel)
		t.Errorf("         : %s\n", outputLevel)
		t.Fail()
	}
}

func TestPlain(t *testing.T) {
	pipe := prepare(LOG_PLAIN)

	originalData := "aAzZ1!?_´→"

	Plain(originalData)

	data := make([]byte, 2<<10)
	pipe.Read(data)

	writtenOutput := string(data)
	writtenOutput = strings.Trim(writtenOutput, "\000")
	writtenOutput = strings.Trim(writtenOutput, "\n")

	if originalData != writtenOutput {
		t.Errorf("Payload does not match")
		t.Errorf("original : %x\n", originalData)
		t.Errorf("written  : %x\n", writtenOutput)
		t.Fail()
	}
}

func TestInfo(t *testing.T) {
	pipe := prepare(LOG_INFO)

	originalData := "aAzZ1!?_´→"

	Info(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_INFO)
}

func TestDebug(t *testing.T) {
	pipe := prepare(LOG_DEBUG)

	originalData := "aAzZ1!?_´→"

	Debug(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_DEBUG)
}

func TestError(t *testing.T) {
	pipe := prepare(LOG_ERROR)

	originalData := "aAzZ1!?_´→"

	Error(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_ERROR)
}

func TestFatal(t *testing.T) {
	originalData := "aAzZ1!?_´→"

	if os.Getenv("LOG_FATAL") == "1" {
		Fatal(originalData)
		return
	}
	readPipe, writePipe, _ := os.Pipe()

	// Starts this test function as separate process to test the "os.Exit(1)" of Fatal
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "LOG_FATAL=1")
	cmd.Stderr = writePipe
	cmd.Stdout = writePipe
	cmd.Run()

	checkSimpleWrite(t, readPipe, originalData, LOG_FATAL)
}

func TestPlainFormat(t *testing.T) {
	pipe := prepare(LOG_PLAIN)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Plain(originalFormat, 123, "bla", "p")

	data := make([]byte, 2<<10)
	pipe.Read(data)

	writtenOutput := string(data)
	writtenOutput = strings.Trim(writtenOutput, "\000")
	writtenOutput = strings.Trim(writtenOutput, "\n")

	if originalData != writtenOutput {
		t.Errorf("Payload does not match")
		t.Errorf("original : %x\n", originalData)
		t.Errorf("written  : %x\n", writtenOutput)
		t.Fail()
	}
}

func TestInfoFormat(t *testing.T) {
	pipe := prepare(LOG_INFO)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Info(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_INFO)
}

func TestDebugFormat(t *testing.T) {
	pipe := prepare(LOG_DEBUG)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Debug(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_DEBUG)
}

func TestErrorFormat(t *testing.T) {
	pipe := prepare(LOG_ERROR)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Error(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_ERROR)
}

func TestFatalFormat(t *testing.T) {
	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	if os.Getenv("LOG_FATAL") == "1" {
		Fatal(originalFormat, 123, "bla", "p")
		return
	}
	readPipe, writePipe, _ := os.Pipe()

	// Starts this test function as separate process to test the "os.Exit(1)" of Fatal
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalFormat")
	cmd.Env = append(os.Environ(), "LOG_FATAL=1")
	cmd.Stderr = writePipe
	cmd.Stdout = writePipe
	cmd.Run()

	checkSimpleWrite(t, readPipe, originalData, LOG_FATAL)
}

// TODO more test regarding the caller information (function name and line)
