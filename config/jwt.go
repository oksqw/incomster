package config

import "time"

type JwtConfig struct {
	Secret   string        `json:"secret"   yaml:"secret"   default:"very-secure-secret"`
	Duration time.Duration `json:"duration" yaml:"duration" default:"24h"`
}
