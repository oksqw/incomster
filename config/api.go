package config

import "fmt"

type ApiConfig struct {
	Host string `yaml:"host" default:"localhost"`
	Port int    `yaml:"port" default:"8080"`
}

func (a ApiConfig) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
