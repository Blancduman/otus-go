package logger

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	logLvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(errors.Wrap(err, "fail to parse log level"))
	}

	logger := zerolog.New(os.Stdout).Level(logLvl)

	return &Logger{
		logger: &logger,
	}
}

func (l Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}
