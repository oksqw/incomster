package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"incomster/config"
)

type DbConnector struct {
	config config.PostgresConfig
}

func NewDbConnector(config config.PostgresConfig) *DbConnector {
	return &DbConnector{config: config}
}

func (d *DbConnector) Connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", d.config.DSN())
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}
