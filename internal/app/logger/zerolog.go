package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"net"
	"os"
	"time"
)

const (
	localCacheSize  = 1000
	elkPollInterval = 10 * time.Millisecond
)

var _ Logger = &ZeroLogger{}

type ZeroLogger struct {
	logger zerolog.Logger
}

func (z *ZeroLogger) Info(message string) {
	z.logger.Info().Msg(message)
}

func (z *ZeroLogger) InfoWithContext(message string, ctx LogContext) {
	z.logger.Info().Fields(ctx).Msg(message)
}

func (z *ZeroLogger) Debug(message string) {
	z.logger.Info().Msg(message)
}

func (z *ZeroLogger) DebugWithContext(message string, ctx LogContext) {
	z.logger.Info().Fields(ctx).Msg(message)
}

func (z *ZeroLogger) WarningWithContext(message string, ctx LogContext) {
	z.logger.Warn().Fields(ctx).Msg(message)
}

func (z *ZeroLogger) Fatal(message string, err error) {
	z.logger.Fatal().Err(err).Msg(message)
}

func (z *ZeroLogger) Error(message string, err error) {
	z.logger.Error().Err(err).Msg(message)
}

func (z *ZeroLogger) ErrorWithContext(message string, err error, ctx LogContext) {
	z.logger.Error().Err(err).Fields(ctx).Msg(message)
}

func NewZeroLogger(logger zerolog.Logger) *ZeroLogger {
	return &ZeroLogger{logger: logger}
}

func NewConsoleELKZeroLogger(lvl, protocol, address string) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(lvl)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	conn, err := net.Dial(protocol, address)
	if err != nil {
		conn = nil
	}

	consoleWriter := newDiodeConsoleWriter(elkPollInterval)

	if err != nil {
		logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
		logger.Error().Err(err).Msg("cant establish connection to ELK, falling back to console logger")
		return logger
	}

	multiWriter := io.MultiWriter(conn, consoleWriter)

	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

	return logger
}

func NewConsoleZeroLogger(lvl string) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(lvl)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	consoleWriter := newDiodeConsoleWriter(time.Second)

	return zerolog.New(consoleWriter).With().Timestamp().Logger()
}

func newDiodeConsoleWriter(pollInterval time.Duration) diode.Writer {
	consoleWriter := diode.NewWriter(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
		localCacheSize,
		pollInterval,
		func(missed int) {
			fmt.Printf("dropped %d messages", missed) //nolint:forbidigo
		})

	return consoleWriter
}
