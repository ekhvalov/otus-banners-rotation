package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type LogLevel string

const (
	LevelDebug   LogLevel = "debug"
	LevelInfo    LogLevel = "info"
	LevelWarning LogLevel = "warning"
	LevelError   LogLevel = "error"
)

func NewLogger(cfg Config, w io.Writer) *Logger {
	l := zerolog.New(w).With().Timestamp().Logger()
	switch cfg.GetLevel() {
	case LevelDebug:
		l = l.Level(zerolog.DebugLevel)
	case LevelWarning:
		l = l.Level(zerolog.WarnLevel)
	case LevelError:
		l = l.Level(zerolog.ErrorLevel)
	case LevelInfo:
		fallthrough
	default:
		l = l.Level(zerolog.InfoLevel)
	}
	return &Logger{logger: l}
}

type Logger struct {
	logger zerolog.Logger
}

func (l Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}
