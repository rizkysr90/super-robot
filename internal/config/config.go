package config

import (
	pgx "api-iad-ams/pkg/pgx"

	"github.com/caarlos0/env/v8"
)

type flatEnv struct {
	AppName           string `env:"APP_NAME"`
	AppEnv            string `env:"APP_ENV"` // local|dev|uat|sit|prod
	RestAPIPort       string `env:"REST_API_PORT"`
	DBHost            string `env:"DB_HOST"`
	DBPort            int    `env:"DB_PORT"`
	DBDatabase        string `env:"DB_DATABASE"`
	DBUsername        string `env:"DB_USERNAME"`
	DBPassword        string `env:"DB_PASSWORD,unset"` // sensitive, unset from environment after parse!
	DBConnMaxOpen     int    `env:"DB_CONN_MAX_OPEN"`
	DBConnMaxIdle     int    `env:"DB_CONN_MAX_IDLE"`
	ApiKey            string `env:"API_KEY,unset"`
	ApiVersionBaseURL string `env:"API_VERSION_BASE_URL"`
	LogLevel          string `env:"LOG_LEVEL"`
}
type Config struct {
	AppName           string
	AppEnv            string
	RestAPIPort       string
	PgSQL             pgx.Config
	APIKey            string
	ApiVersionBaseURL string
	LogLevel          string
}

func LoadFromEnv() (Config, error) {
	var envCfg flatEnv
	err := env.Parse(&envCfg)
	if err != nil {
		return Config{}, err
	}
	return newConfig(envCfg), nil
}
func newConfig(envCfg flatEnv) Config {
	return Config{
		AppName:     envCfg.AppName,
		AppEnv:      envCfg.AppEnv,
		RestAPIPort: envCfg.RestAPIPort,
		PgSQL: pgx.Config{
			Username: envCfg.DBUsername,
			Password: envCfg.DBPassword,
			Database: envCfg.DBDatabase,
			Host:     envCfg.DBHost,
			Port:     envCfg.DBPort,
		},
		APIKey:            envCfg.ApiKey,
		ApiVersionBaseURL: envCfg.ApiVersionBaseURL,
		LogLevel:          envCfg.LogLevel,
	}
}
