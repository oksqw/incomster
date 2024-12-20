package main

import (
	"golang.org/x/net/context"
	"incomster/backend/store/migrate"
	"incomster/backend/store/postgres"
	"incomster/config"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load[config.Config]("incomster")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := postgres.Connect(ctx, cfg.Store.Postgres)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	migrations, err := migrate.ParseMigrations(postgres.FS)
	if err != nil {
		log.Fatalf("migrations: %v", err)
	}

	migrator := postgres.NewMigrator(db, migrations)
	from, to, err := migrator.Up(ctx)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Printf("migrate : %d => %d", from, to)
}
