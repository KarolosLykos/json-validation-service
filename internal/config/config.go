package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug bool `envconfig:"DEBUG" default:"false"`
	HTTP  HTTP
}

type HTTP struct {
	IP   string `envconfig:"HTTP_IP" default:"0.0.0.0"`
	Port string `envconfig:"HTTP_PORT" default:"8082"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
