package logging

import (
	"github.com/rs/zerolog/log"
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

	// TODO: Не уверен что это корректный подход.
	// По-сути мы тут глобальному логгеру присваиваем только что созданный логгер.
	// Мне это выглядит нормально, но я не уверен что это не вызовет каких-то проблем.
	log.Logger = zerolog.New(w).Level(l).With().Timestamp().Logger()
	return log.Logger.WithContext(ctx)
}
