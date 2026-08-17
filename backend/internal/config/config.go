package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config 集中解析所有环境变量配置。
type Config struct {
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	HTTPPort   string `env:"HTTP_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBName     string `env:"DB_NAME" envDefault:"wmsflow_db"`
	DBUser     string `env:"DB_USER" envDefault:"wmsflow_user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"wmsflow_pwd"`
	JWTSecret  string `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_string"`
	JWTExpire  int    `env:"JWT_EXPIRE_HOURS" envDefault:"72"`
	RedisHost  string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort  string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPass  string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB    int    `env:"REDIS_DB" envDefault:"0"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
}

// Load 从环境变量加载配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config from env: %w", err)
	}
	if len(cfg.JWTSecret) < 16 {
		return nil, fmt.Errorf("config jwt_secret too short: %w", ErrWeakJWTSecret)
	}
	return cfg, nil
}
