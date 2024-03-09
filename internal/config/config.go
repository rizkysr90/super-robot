package config

import (
	pgx "github.com/rizkysr90/go-boilerplate/pkg/pgx"

	"github.com/caarlos0/env/v8"
)

type flatEnv struct {
	AppName           string `env:"APP_NAME"`
	AppEnv            string `env:"APP_ENV"` // local|dev|uat|sit|prod
	RestAPIPort       string `env:"REST_API_PORT"`
	DBHost            string `env:"DB_HOST"`
	DBDatabase        string `env:"DB_DATABASE"`
	DBUsername        string `env:"DB_USERNAME"`
	DBPassword        string `env:"DB_PASSWORD,unset"` // sensitive, unset from environment after parse!
	APIKey            string `env:"API_KEY,unset"`
	APIVersionBaseURL string `env:"API_VERSION_BASE_URL"`
	LogLevel          string `env:"LOG_LEVEL"`
	DBPort            int    `env:"DB_PORT"`
	DBConnMaxOpen     int    `env:"DB_CONN_MAX_OPEN"`
	DBConnMaxIdle     int    `env:"DB_CONN_MAX_IDLE"`
	PrivateKeyJWT     string `env:"PRIVATE_KEY_JWT,unset"`
	PublicKeyJWT      string `env:"PUBLIC_KEY_JWT,unset"`
}
type Config struct {
	AppName           string
	AppEnv            string
	RestAPIPort       string
	APIKey            string
	APIVersionBaseURL string
	LogLevel          string
	PgSQL             pgx.Config
	PrivateKeyJWT     string
	PublicKeyJWT      string
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
		APIKey:            envCfg.APIKey,
		APIVersionBaseURL: envCfg.APIVersionBaseURL,
		LogLevel:          envCfg.LogLevel,
		PrivateKeyJWT:     envCfg.PrivateKeyJWT,
		PublicKeyJWT:      envCfg.PublicKeyJWT,
	}
}
