package application

import (
	"context"
	"errors"
	"net/http"

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

	"github.com/rs/zerolog/log"
)

type App struct {
	store     store.IStore
	config    *config.Config
	service   *service.Service
	server    *http.Server
	tokenizer *jwt.Tokenizer
	closer    *closer.Closer
}

func New(config *config.Config) *App {
	return &App{config: config, closer: closer.New()}
}

func (a *App) Setup(ctx context.Context) error {
	log.Ctx(ctx).Info().Str("status", "trying").Msg("setup")

	if err := a.setupTokenizer(ctx); err != nil {
		return err
	}
	if err := a.setupStore(ctx); err != nil {
		return err
	}
	if err := a.setupService(ctx); err != nil {
		return err
	}
	if err := a.setupServer(ctx); err != nil {
		return err
	}

	log.Ctx(ctx).Info().Str("status", "OK").Msg("setup")
	return nil
}

func (a *App) Run(ctx context.Context) error {
	log.Ctx(ctx).Info().Str("address", a.config.Api.String()).Msg("running...")

	closed := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Ctx(ctx).Info().Str("status", "trying").Msg("shutdown")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		if err := a.Shutdown(shutdownCtx); err != nil {
			log.Ctx(ctx).Warn().Str("status", "failed").Err(err).Msg("shutdown")
		} else {
			log.Ctx(ctx).Info().Str("status", "gracefully").Msg("shutdown")
		}

		close(closed)
	}()

	if e := a.server.ListenAndServe(); !errors.Is(e, http.ErrServerClosed) {
		return e
	}

	<-closed
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Ctx(ctx).Info().Dur("timeout", a.config.ShutdownTimeout).Msg("shutdown...")
	return a.closer.CloseSequentially(ctx, a.config.ShutdownTimeout)
}

func (a *App) setupTokenizer(ctx context.Context) error {
	log.Ctx(ctx).Debug().Str("status", "trying").Msg("tokenizer setup")

	tokenizer, err := jwt.New(a.config.Jwt.Secret, a.config.Jwt.Duration)
	if err != nil {
		return err
	}
	a.tokenizer = tokenizer

	log.Ctx(ctx).Debug().Str("status", "OK").Msg("tokenizer setup")
	return nil
}

func (a *App) setupStore(ctx context.Context) error {
	log.Ctx(ctx).Debug().Str("status", "trying").Msg("store setup")

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
		err = db.Close()
		if err != nil {
			log.Ctx(ctx).Debug().Str("status", "failed").Err(err).Msg("store shutdown")
		}
		return err
	})

	log.Ctx(ctx).Debug().Str("status", "OK").Msg("store setup")
	return nil
}

func (a *App) setupService(ctx context.Context) error {
	log.Ctx(ctx).Debug().Str("status", "trying").Msg("service setup")

	user := service.NewUserService(a.store.User())
	income := service.NewIncomeService(a.store.Income())
	account := service.NewAccountService(a.store.Session(), a.store.User(), a.tokenizer, a.config)
	security := service.NewSecurityService(a.store.Session(), a.tokenizer)
	a.service = service.NewService(user, income, account, security)

	log.Ctx(ctx).Debug().Str("status", "OK").Msg("service setup")
	return nil
}

func (a *App) setupServer(ctx context.Context) error {
	log.Ctx(ctx).Debug().Str("status", "trying").Msg("server setup")

	v := validation.NewValidator()
	h := api.NewHandler(a.config, a.service, v)
	s, e := oas.NewServer(h, h)
	if e != nil {
		return e
	}

	a.server = &http.Server{
		Addr:    a.config.Api.String(),
		Handler: s,
	}

	a.closer.Add(func(ctx context.Context) error {
		err := a.server.Shutdown(ctx)
		if err != nil {
			log.Ctx(ctx).Debug().Str("status", "failed").Err(err).Msg("server shutdown")
		}
		return err
	})

	log.Ctx(ctx).Debug().Str("status", "OK").Msg("server setup")
	return nil
}
