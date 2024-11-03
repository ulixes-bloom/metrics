package logger

import (
	"io"

	"github.com/rs/zerolog"
)

func Initialize(level string, w io.Writer) (logger zerolog.Logger, err error) {
	logger = zerolog.New(w)
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return
	}
	logger = logger.Level(lvl).With().Timestamp().Logger()

	return
}
