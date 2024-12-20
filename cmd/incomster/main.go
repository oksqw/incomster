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

	//connector := postgres.NewDbConnector(cfg.Store.Postgres)
	//db, err := connector.Connect(ctx)
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}
	//
	//migrations, err := migrate.ParseMigrations(postgres.FS)
	//if err != nil {
	//	log.Fatalf("parse migrations: %s", err)
	//}
	//
	//migrator, err := postgres.NewMigrator(db, migrations)
	//if err != nil {
	//	log.Fatalf("Failed to initialize migrator: %v", err)
	//}

	//fmt.Printf("config: %+v\n", cfg)
	//
	//ctx := context.Background()
	//ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	//defer cancel()
	//
	//db, err := sql.Open("postgres", cfg.Store.Postgres.DSN())
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}
	//defer db.Close()
	//
	//migrations, err := migrate.ParseMigrations(postgres.FS)
	//if err != nil {
	//	log.Fatalf("parse migrations: %s", err)
	//}
	//
	//migrator, err := postgres.NewMigrator(db, migrations)
	//if err != nil {
	//	log.Fatalf("Failed to initialize migrator: %v", err)
	//}
	//
	//from, to, err := migrator.Up(ctx)
	//if err != nil {
	//	log.Fatalf("migrate up: %s", err)
	//}
	//
	//log.Printf("from : %d | to : %d | store postgres migrate up", from, to)
}
