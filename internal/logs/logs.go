package logs

import (
	"fmt"
	"os"
)

type Logger struct {
	logs chan string
}

func New() *Logger {
	return &Logger{
		logs: make(chan string, 1024),
	}
}

func (l *Logger) Log(message string) {
	l.logs <- message
}

func (l *Logger) Handlelogs() {
	for log := range l.logs {
		fmt.Fprintln(os.Stdout, log)
	}
}
