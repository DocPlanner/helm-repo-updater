package logger

import (
	"os"
)

var _ Logger = &NullLogger{}

type LogContext map[string]interface{}

type Logger interface {
	WarningWithContext(message string, ctx LogContext)
	Error(message string, err error)
	ErrorWithContext(message string, err error, ctx LogContext)
	Fatal(message string, err error)
	InfoWithContext(message string, ctx LogContext)
	Debug(message string)
	DebugWithContext(message string, ctx LogContext)
}

type NullLogger struct {
}

func (n *NullLogger) Fatal(_ string, _ error) {
	os.Exit(1)
}

func (n *NullLogger) WarningWithContext(_ string, _ LogContext) {
	// yoink!
}

func (n *NullLogger) Error(_ string, _ error) {
	// yoink!
}

func (n *NullLogger) ErrorWithContext(_ string, _ error, _ LogContext) {
	// yoink!
}

func (n NullLogger) InfoWithContext(_ string, _ LogContext) {
	// nada!
}

func (n NullLogger) Debug(_ string) {
	// yoink!
}

func (n NullLogger) DebugWithContext(_ string, _ LogContext) {
	// yoink!
}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}
