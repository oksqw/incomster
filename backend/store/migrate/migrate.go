package migrate

import (
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	Script  string
}

func ParseMigrations(fsys fs.FS) ([]Migration, error) {
	var migrations []Migration

	err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		versionRaw, name, ok := strings.Cut(strings.TrimPrefix(path, "migrations/"), "_")
		if !ok {
			return fmt.Errorf("invalid migration name: %s", path)
		}

		version, err := strconv.Atoi(versionRaw)
		if err != nil {
			return fmt.Errorf("invalid migration version: %s", path)
		}

		script, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read migration: %w", err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Script:  string(script),
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk migrations: %w", err)
	}

	slices.SortFunc(migrations, func(a, b Migration) int {
		return a.Version - b.Version
	})

	return migrations, nil
}
