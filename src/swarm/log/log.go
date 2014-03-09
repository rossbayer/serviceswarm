package log

import (
	golog "log"
	"io"
	"os"
)

const logFlags int = golog.Ldate|golog.Ltime | golog.Lshortfile

var logWriters []io.Writer = make([]io.Writer, 0, 10)

func AddLogWriter(out io.Writer) {
	logWriters = append(logWriters, out)
}

func NewLogger(name string) *golog.Logger {
	if len(logWriters) == 0 {
		initDefaultLogWriters()
	}

	return golog.New(io.MultiWriter(logWriters...), name+" ", logFlags)
}

func initDefaultLogWriters() {
	logWriters = append(logWriters, os.Stderr)
}
