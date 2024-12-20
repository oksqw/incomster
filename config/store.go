package config

import (
	"fmt"
	"time"
)

type StoreConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host               string        `yaml:"host"                 default:"localhost"`
	Port               int           `yaml:"port"                 default:"5432"`
	User               string        `yaml:"user"                 default:"incomster"`
	Password           string        `yaml:"password"             default:"incomster"`
	Database           string        `yaml:"database"             default:"incomster"`
	SSLMode            string        `yaml:"ssl_mode"             default:"disable"`
	MaxOpenConnections int           `yaml:"max_open_connections" default:"25"`
	MaxIdleConnections int           `yaml:"max_idle_connections" default:"10"`
	Timeout            time.Duration `yaml:"timeout"              default:"60s"`
}

func (cfg PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
}
