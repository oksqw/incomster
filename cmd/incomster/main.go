package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"incomster/application"
	"incomster/backend/logging"
	"incomster/config"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load[config.Config]("incomster")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	ctx := logging.NewLoggerContext(context.Background(), &cfg)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	app := application.New(&cfg)

	err = app.Setup(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("setup")
	}

	err = app.Run(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("run")
	}
}
