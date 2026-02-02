package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
	JWT        JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env:"HTTP_ADDR" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

type Postgres struct {
	Host        string `yaml:"host" env-default:"localhost"`
	Port        int    `yaml:"port" env-default:"5432"`
	DB          string `yaml:"db" env-default:"postgres"`
	User        string `yaml:"user" env-default:"postgres"`
	Password    string `yaml:"password" env-default:"postgres"`
	MaxAttempts int    `yaml:"max_attempts" env-default:"5"`
}

type JWT struct {
	Secret string        `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	Ttl    time.Duration `yaml:"ttl" env:"JWT_TTL" env-default:"24h"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
