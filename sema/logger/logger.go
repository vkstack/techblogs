package logger

import (
	"fmt"
	"os"
)

type Logger interface {
	Error(...any)
	Info(...any)
}

type stdlogger struct{}

func NewLogger() Logger {
	return &stdlogger{}
}

func (l *stdlogger) Error(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func (l *stdlogger) Info(args ...any) {
	fmt.Fprintln(os.Stdout, args...)
}
