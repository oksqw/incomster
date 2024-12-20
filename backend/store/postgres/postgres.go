package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"incomster/config"
	"log"
)

func Connect(ctx context.Context, config config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}

func CommitOrRollback(tx *sql.Tx, err error) {
	if err != nil {
		rollback(tx)
	}
	commit(tx)
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		// Наверное тут нужен более мягкий способ завершения приложения
		log.Fatalf("failed to rollback transaction: %v", err)
	}
}

func commit(tx *sql.Tx) {
	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
	}
}
