package config

type LogConfig struct {
	Pretty bool `json:"pretty" yaml:"pretty" default:"false"`
	Debug  bool `json:"debug"  yaml:"debug"  default:"false"`
}
