package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server" env-prefix:"HTTP_"`
	Postgres   Postgres   `yaml:"postgres" env-prefix:"POSTGRES_"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env:"HTTP_ADDR" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

type Postgres struct {
	Host        string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port        int    `yaml:"port" env:"PORT" env-default:"5432"`
	DB          string `yaml:"db" env:"DB" env-default:"postgres"`
	User        string `yaml:"user" env:"USER" env-default:"root"`
	Password    string `yaml:"password" env:"PASSWORD" env-default:"1234"`
	MaxAttempts int    `yaml:"max_attempts" env:"MAX_ATTEMPTS" env-default:"5"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
