package logger

import (
	"os"
)

var _ Logger = &NullLogger{}

type Logger interface {
	ErrorWithContext(message string, err error, context map[string]interface{})
	Fatal(message string, err error)
	InfoWithContext(message string, context map[string]interface{})
	Info(message string)
}

type NullLogger struct {
}

func (n *NullLogger) Fatal(_ string, _ error) {
	os.Exit(1)
}

func (n *NullLogger) ErrorWithContext(_ string, _ error, _ map[string]interface{}) {
	// nada!
}

func (n *NullLogger) InfoWithContext(_ string, _ map[string]interface{}) {
	// nada!
}

func (n NullLogger) Info(_ string) {
	// nada!
}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}
