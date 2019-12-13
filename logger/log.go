package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type SimpleLog struct {
	mu  sync.Mutex
	out io.Writer
}

var simpleLog = SimpleLog{
	out: os.Stderr,
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()
	simpleLog.out = w
}

func Info(format string, a ...interface{}) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()

	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(simpleLog.out, "\x1b[92minfo:\x1b[0m", out)
}

func Warn(format string, a ...interface{}) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()

	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(simpleLog.out, "\x1b[93mwarn:\x1b[0m", out)
}

func Error(format string, a ...interface{}) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()

	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(simpleLog.out, "\x1b[91merror:\x1b[0m", out)
}

func FatalError(format string, a ...interface{}) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()

	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(simpleLog.out, "\x1b[91merror:\x1b[0m", out)
	os.Exit(1)
}

func Event(format string, a ...interface{}) {
	simpleLog.mu.Lock()
	defer simpleLog.mu.Unlock()

	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(simpleLog.out, "\x1b[94mevent:\x1b[0m", out)
}
