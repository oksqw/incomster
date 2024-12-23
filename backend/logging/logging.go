package logging

import (
	"io"
	"os"

	"incomster/config"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

func NewLoggerContext(ctx context.Context, cfg *config.Config) context.Context {
	var (
		w io.Writer     = os.Stderr
		l zerolog.Level = zerolog.InfoLevel
	)

	if cfg.Log.Pretty {
		w = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.00"}
	}
	if cfg.Log.Debug {
		l = zerolog.DebugLevel
	}

	return zerolog.New(w).Level(l).With().Timestamp().Logger().WithContext(ctx)
}
