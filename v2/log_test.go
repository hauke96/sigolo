package sigolo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func prepare(logLevel Level) *os.File {
	SetDefaultLogLevel(LOG_PLAIN)

	readPipe, writePipe, _ := os.Pipe()

	levelOutputs[logLevel] = writePipe

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
		t.Errorf("Data does not match")
		t.Errorf("original : %x\n", originalData)
		t.Errorf("         : %s\n", originalData)
		t.Errorf("written  : %x\n", outputData)
		t.Errorf("         : %s\n", outputData)
		t.Fail()
	}

	// Log-level INFO has an additional space at its end because the string is shorter than others
	if logLevel == LOG_INFO || logLevel == LOG_WARN {
		outputLevel += " "
	}

	if levelStrings[logLevel] != outputLevel {
		t.Errorf("Log-level string does not patch")
		t.Errorf("original : %x\n", levelStrings[logLevel])
		t.Errorf("         : %s\n", levelStrings[logLevel])
		t.Errorf("written  : %x\n", outputLevel)
		t.Errorf("         : %s\n", outputLevel)
		t.Fail()
	}
}

func TestPlain(t *testing.T) {
	pipe := prepare(LOG_PLAIN)

	originalData := "aAzZ1!?_´→"

	Plainf(originalData)

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

	Infof(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_INFO)
}

func TestDebug(t *testing.T) {
	pipe := prepare(LOG_DEBUG)

	originalData := "aAzZ1!?_´→"

	Debugf(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_DEBUG)
}

func TestTrace(t *testing.T) {
	pipe := prepare(LOG_TRACE)

	originalData := "aAzZ1!?_´→"

	Tracef(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_TRACE)
}

func TestWarn(t *testing.T) {
	pipe := prepare(LOG_WARN)

	originalData := "aAzZ1!?_´→"

	Warnf(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_WARN)
}

func TestError(t *testing.T) {
	pipe := prepare(LOG_ERROR)

	originalData := "aAzZ1!?_´→"

	Errorf(originalData)

	checkSimpleWrite(t, pipe, originalData, LOG_ERROR)
}

func TestFatalf(t *testing.T) {
	originalData := "aAzZ1!?_´→"

	if os.Getenv("LOG_FATAL") == "1" {
		Fatalf("%s", originalData)
		return
	}
	readPipe, writePipe, _ := os.Pipe()

	// Starts this test function as separate process to test the "os.Exit(1)" of Fatalf
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf")
	cmd.Env = append(os.Environ(), "LOG_FATAL=1")
	cmd.Stderr = writePipe
	cmd.Stdout = writePipe
	err := cmd.Run()

	checkSimpleWrite(t, readPipe, originalData, LOG_FATAL)

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ExitCode() != 1 {
				t.Errorf("Exit code was not 1 but %s", exiterr)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	} else {
		t.Errorf("Expected command to exit with error")
	}
}

func TestFatal(t *testing.T) {
	originalData := "aAzZ1!?_´→"

	if os.Getenv("LOG_FATAL") == "123" {
		Fatal(originalData)
		return
	}
	readPipe, writePipe, _ := os.Pipe()

	// Starts this test function as separate process to test the "os.Exit(1)" of Fatalf
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "LOG_FATAL=123")
	cmd.Stderr = writePipe
	cmd.Stdout = writePipe
	err := cmd.Run()

	checkSimpleWrite(t, readPipe, originalData, LOG_FATAL)

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.ExitCode() != 1 {
				t.Errorf("Exit code was not 1 but %s", exiterr)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	} else {
		t.Errorf("Expected command to exit with error")
	}
}

func TestPlainFormat(t *testing.T) {
	pipe := prepare(LOG_PLAIN)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Plainf(originalFormat, 123, "bla", "p")

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

	Infof(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_INFO)
}

func TestDebugFormat(t *testing.T) {
	pipe := prepare(LOG_DEBUG)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Debugf(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_DEBUG)
}

func TestTraceFormat(t *testing.T) {
	pipe := prepare(LOG_TRACE)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Tracef(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_TRACE)
}

func TestErrorFormat(t *testing.T) {
	pipe := prepare(LOG_ERROR)

	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	Errorf(originalFormat, 123, "bla", "p")

	checkSimpleWrite(t, pipe, originalData, LOG_ERROR)
}

func TestFatalFormat(t *testing.T) {
	originalData := "foo_123_bla_70"
	originalFormat := "foo_%d_%s_%x"

	if os.Getenv("LOG_FATAL") == "1" {
		Fatalf(originalFormat, 123, "bla", "p")
		return
	}
	readPipe, writePipe, _ := os.Pipe()

	// Starts this test function as separate process to test the "os.Exit(1)" of Fatalf
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalFormat")
	cmd.Env = append(os.Environ(), "LOG_FATAL=1")
	cmd.Stderr = writePipe
	cmd.Stdout = writePipe
	cmd.Run()

	checkSimpleWrite(t, readPipe, originalData, LOG_FATAL)
}

func TestShouldLog(t *testing.T) {
	logLevel = LOG_INFO
	assertTrue(t, ShouldLog(LOG_FATAL))
	assertTrue(t, ShouldLog(LOG_ERROR))
	assertTrue(t, ShouldLog(LOG_WARN))
	assertTrue(t, ShouldLog(LOG_INFO))
	assertFalse(t, ShouldLog(LOG_DEBUG))
	assertFalse(t, ShouldLog(LOG_TRACE))
	assertFalse(t, ShouldLog(LOG_PLAIN))
	assertFalse(t, ShouldLogDebug())
	assertFalse(t, ShouldLogTrace())

	logLevel = LOG_FATAL
	assertTrue(t, ShouldLog(LOG_FATAL))
	assertFalse(t, ShouldLog(LOG_ERROR))
	assertFalse(t, ShouldLog(LOG_WARN))
	assertFalse(t, ShouldLog(LOG_INFO))
	assertFalse(t, ShouldLog(LOG_DEBUG))
	assertFalse(t, ShouldLog(LOG_TRACE))
	assertFalse(t, ShouldLog(LOG_PLAIN))
	assertFalse(t, ShouldLogDebug())
	assertFalse(t, ShouldLogTrace())

	logLevel = LOG_TRACE
	assertTrue(t, ShouldLog(LOG_FATAL))
	assertTrue(t, ShouldLog(LOG_ERROR))
	assertTrue(t, ShouldLog(LOG_WARN))
	assertTrue(t, ShouldLog(LOG_INFO))
	assertTrue(t, ShouldLog(LOG_DEBUG))
	assertTrue(t, ShouldLog(LOG_TRACE))
	assertFalse(t, ShouldLog(LOG_PLAIN))
	assertTrue(t, ShouldLogDebug())
	assertTrue(t, ShouldLogTrace())
}

// TODO more test regarding the caller information (function name and line)

func assertTrue(t *testing.T, b bool) {
	if !b {
		fmt.Println("Expected true but got false")
		t.Fail()
	}
}

func assertFalse(t *testing.T, b bool) {
	if b {
		fmt.Println("Expected false but got true")
		t.Fail()
	}
}
