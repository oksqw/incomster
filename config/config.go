package config

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"strings"
	"time"
)

type Env string

func (e Env) IsDev() bool {
	return e == "dev"
}

type Config struct {
	Env             Env           `yaml:"env"              default:"dev"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"10s"`
	Store           StoreConfig   `yaml:"store"`
	Api             ApiConfig     `yaml:"api"`
	Jwt             JwtConfig     `yaml:"jwt"`
	Log             LogConfig     `yaml:"log"`
}

func Load[T any](name string) (T, error) {
	var cfg T

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix: strings.ReplaceAll(strings.ToUpper(name), "-", "_"),
		Files: []string{
			fmt.Sprintf("./%s.local.yaml", strings.ToLower(name)),
			fmt.Sprintf("./%s.yaml", strings.ToLower(name)),
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}
