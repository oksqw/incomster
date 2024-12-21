package main

import (
	"context"
	"incomster/application"
	"incomster/config"
	"log"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load[config.Config]("incomster")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	app := application.New(&cfg)
	err = app.Setup(ctx)
	if err != nil {
		log.Fatalf("setup: %v", err)
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("run: %v", err)
	}
}
