package main

import (
	"context"
	"incomster/application"
	"incomster/config"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load[config.Config]("incomster")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app := application.New(&cfg)
	err = app.Setup(ctx)
	if err != nil {
		log.Fatalf("setup: %v", err)
	}

	go run(ctx, app)
	shutdown(ctx, app)
}

func run(ctx context.Context, app *application.App) {
	err := app.Run(ctx)
	if err != nil {
		log.Fatalf("run: %v", err)
	}
}

func shutdown(ctx context.Context, app *application.App) {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}

	log.Println("application shut down gracefully")
}
