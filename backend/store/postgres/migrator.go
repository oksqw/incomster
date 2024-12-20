package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"incomster/backend/store/migrate"
)

//go:embed migrations/*.sql
var FS embed.FS

type Migrator struct {
	db         *sql.DB
	migrations []migrate.Migration
}

func NewMigrator(db *sql.DB, migrations []migrate.Migration) *Migrator {
	return &Migrator{db: db, migrations: migrations}
}

func (m *Migrator) Up(ctx context.Context) (from, to int, err error) {
	if err = m.init(ctx); err != nil {
		return 0, 0, fmt.Errorf("init: %w", err)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	current, err := m.version(ctx, tx)
	if err != nil {
		return 0, 0, fmt.Errorf("get version: %w", err)
	}

	to, err = m.apply(ctx, tx, current)
	if err != nil {
		return 0, 0, fmt.Errorf("apply migrations: %w", err)
	}

	return current, to, nil
}

func (m *Migrator) Down(ctx context.Context, targetVersion int) (from, to int, err error) {
	if err = m.init(ctx); err != nil {
		return 0, 0, fmt.Errorf("init: %w", err)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	current, err := m.version(ctx, tx)
	if err != nil {
		return 0, 0, fmt.Errorf("get version: %w", err)
	}

	if targetVersion >= current {
		return current, targetVersion, fmt.Errorf("target version %d is not less than current version %d", targetVersion, current)
	}

	to, err = m.rollback(ctx, tx, current, targetVersion)
	if err != nil {
		return 0, 0, fmt.Errorf("rollback migrations: %w", err)
	}

	return current, to, nil
}

func (m *Migrator) Migrations() ([]migrate.Migration, error) {
	return m.migrations, nil
}

func (m *Migrator) init(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS _migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	return nil
}

func (m *Migrator) version(ctx context.Context, tx *sql.Tx) (int, error) {
	var version int

	if err := tx.QueryRowContext(ctx, `SELECT version FROM _migrations ORDER BY version DESC LIMIT 1`).Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("select version: %w", err)
	}

	return version, nil
}

func (m *Migrator) apply(ctx context.Context, tx *sql.Tx, from int) (int, error) {
	current := from

	for _, migration := range m.migrations {
		if migration.Version <= current {
			continue
		}

		if _, err := tx.ExecContext(ctx, migration.Script); err != nil {
			return 0, fmt.Errorf("exec migration: %w", err)
		}

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO _migrations (version, name) VALUES ($1, $2)`,
			migration.Version,
			migration.Name,
		); err != nil {
			return 0, fmt.Errorf("insert migration: %w", err)
		}

		current++
	}

	return current, nil
}

func (m *Migrator) rollback(ctx context.Context, tx *sql.Tx, from, to int) (int, error) {
	current := from

	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]

		if migration.Version > current || migration.Version <= to {
			continue
		}

		if _, err := tx.ExecContext(ctx, migration.Script); err != nil {
			return 0, fmt.Errorf("exec migration down: %w", err)
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM _migrations WHERE version = $1`, migration.Version); err != nil {
			return 0, fmt.Errorf("delete migration record: %w", err)
		}

		current--
	}

	return current, nil
}
