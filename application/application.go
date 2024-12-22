package application

import (
	"context"
	"incomster/backend/api/handler"
	"incomster/backend/api/oas"
	"incomster/backend/api/validation"
	"incomster/backend/service"
	"incomster/backend/store"
	"incomster/backend/store/migrate"
	"incomster/backend/store/postgres"
	"incomster/config"
	"incomster/pkg/closer"
	"incomster/pkg/jwt"
	"log"
	"net/http"
	"time"
)

type App struct {
	store     store.IStore
	config    *config.Config
	service   *service.Service
	handler   *api.Handler
	tokenizer *jwt.Tokenizer
	closer    *closer.Closer
}

func New(config *config.Config) *App {
	return &App{config: config, closer: closer.New()}
}

func (a *App) Setup(ctx context.Context) error {
	// TODO: use ctx logger
	log.Print("trying to setup application")

	if err := a.setupTokenizer(ctx); err != nil {
		return err
	}
	if err := a.setupStore(ctx); err != nil {
		return err
	}
	if err := a.setupService(ctx); err != nil {
		return err
	}
	if err := a.setupHandler(ctx); err != nil {
		return err
	}

	log.Print("setup application complete")
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// TODO: use ctx logger

	s, e := oas.NewServer(a.handler, a.handler)
	if e != nil {
		return e
	}

	if e = http.ListenAndServe(a.config.Api.String(), s); e != nil {
		return e
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// TODO: use shutdown timeout from config
	return a.closer.CloseSequentially(ctx, 5*time.Second)
}

func (a *App) setupTokenizer(ctx context.Context) error {
	// TODO: use ctx logger
	log.Print("tokenizer : trying")

	tokenizer, err := jwt.New(a.config.Jwt.Secret, a.config.Jwt.Duration)
	if err != nil {
		return err
	}
	a.tokenizer = tokenizer

	log.Print("tokenizer : ok")
	return nil
}

func (a *App) setupStore(ctx context.Context) error {
	// TODO: use ctx logger
	log.Print("store : trying")

	db, err := postgres.Connect(ctx, a.config.Store.Postgres)
	if err != nil {
		return err
	}

	migrations, err := migrate.ParseMigrations(postgres.FS)
	if err != nil {
		return err
	}

	migrator := postgres.NewMigrator(db, migrations)
	_, _, err = migrator.Up(ctx)
	if err != nil {
		return err
	}

	user := postgres.NewUserStore(db)
	income := postgres.NewIncomeStore(db)
	session := postgres.NewSessionStore(db)
	a.store = postgres.NewStore(user, income, session)

	a.closer.Add(func(ctx context.Context) error {
		// TODO: use ctx logger
		return db.Close()
	})

	log.Print("store : OK")
	return nil
}

func (a *App) setupService(ctx context.Context) error {
	// TODO: use ctx logger
	log.Print("service: trying")

	user := service.NewUserService(a.store.User())
	income := service.NewIncomeService(a.store.Income())
	account := service.NewAuthorizationService(a.store.Session(), a.store.User(), a.tokenizer, a.config)
	security := service.NewSecurityService(a.store.Session(), a.tokenizer)
	a.service = service.NewService(user, income, account, security)

	log.Print("service: OK")
	return nil
}

func (a *App) setupHandler(ctx context.Context) error {
	// TODO: use ctx logger
	log.Print("handler: trying")

	validator := validation.NewValidator()
	a.handler = api.NewHandler(a.config, a.service, validator)

	log.Print("handler: OK")
	return nil
}
