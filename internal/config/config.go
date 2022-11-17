package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug   bool `envconfig:"DEBUG" default:"true"`
	HTTP    HTTP
	Storage Storage
}

type HTTP struct {
	IP   string `envconfig:"HTTP_IP" default:"0.0.0.0"`
	Port string `envconfig:"HTTP_PORT" default:"8082"`
}

type Storage struct {
	HOST     string `envconfig:"DB_HOST" default:"localhost"`
	PORT     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Name     string `envconfig:"DB_NAME" default:"json-validation-service"`
	Password string `envconfig:"DB_PASSWORD" default:"mysecretpassword"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("JSON_VALIDATION_SERVICE", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
