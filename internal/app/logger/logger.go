package logger

import (
	"os"
)

var _ Logger = &NullLogger{}

type logContext map[string]interface{}

type Logger interface {
	WarningWithContext(message string, ctx logContext)
	Error(message string, err error)
	ErrorWithContext(message string, err error, ctx logContext)
	Fatal(message string, err error)
	Info(message string)
	InfoWithContext(message string, ctx logContext)
	DebugWithContext(message string, ctx logContext)
}

type NullLogger struct {
}

func (n *NullLogger) Fatal(_ string, _ error) {
	os.Exit(1)
}

func (n *NullLogger) WarningWithContext(_ string, _ logContext) {
	// yoink!
}

func (n *NullLogger) Error(_ string, _ error) {
	// yoink!
}

func (n *NullLogger) ErrorWithContext(_ string, _ error, _ logContext) {
	// yoink!
}

func (n NullLogger) Info(_ string) {
	// nada!
}

func (n NullLogger) InfoWithContext(_ string, _ logContext) {
	// nada!
}

func (n NullLogger) DebugWithContext(_ string, _ logContext) {
	// yoink!
}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}
