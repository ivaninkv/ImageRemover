package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var (
	hostname, _ = os.Hostname()
	Log         = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().
			Timestamp().
			Str("host", hostname).
			Logger()
)
