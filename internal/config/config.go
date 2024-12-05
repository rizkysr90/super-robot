package config

import (
	"rizkysr90-pos/internal/auth"

	pgx "github.com/rizkysr90/rizkysr90-go-pkg/pgx"

	"github.com/caarlos0/env/v8"
)

type flatEnv struct {
	APIVersionBaseURL string `env:"API_VERSION_BASE_URL"`
	DBPassword        string `env:"DB_PASSWORD,unset"`
	RestAPIPort       string `env:"REST_API_PORT"`
	DBHost            string `env:"DB_HOST"`
	AppName           string `env:"APP_NAME"`
	DBUsername        string `env:"DB_USERNAME"`
	AppEnv            string `env:"APP_ENV"`
	APIKey            string `env:"API_KEY,unset"`
	DBDatabase        string `env:"DB_DATABASE"`
	LogLevel          string `env:"LOG_LEVEL"`
	SecretKeyJWT      string `env:"SECRET_KEY_JWT,unset"`
	PublicKeyJWT      string `env:"PUBLIC_KEY_JWT,unset"`
	PrivateKeyJWT     string `env:"PRIVATE_KEY_JWT,unset"`
	DBConnMaxIdle     int    `env:"DB_CONN_MAX_IDLE"`
	DBConnMaxOpen     int    `env:"DB_CONN_MAX_OPEN"`
	DBPort            int    `env:"DB_PORT"`
	AuthClientID      string `env:"AUTH_CLIENT_ID"`
	AuthRedirectUri   string `env:"AUTH_REDIRECT_URI"`
	AuthClientSecret  string `env:"AUTH_CLIENT_SECRET"`
	AuthUri           string `env:"AUTH_URI"`
}
type Config struct {
	AppName           string
	AppEnv            string
	RestAPIPort       string
	APIKey            string
	APIVersionBaseURL string
	LogLevel          string
	SecretKeyJWT      string
	PgSQL             pgx.Config
	Auth              *auth.Config
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
		SecretKeyJWT:      envCfg.SecretKeyJWT,
		Auth: &auth.Config{
			BaseURL:      envCfg.AuthUri,
			ClientID:     envCfg.AuthClientID,
			RedirectURI:  envCfg.AuthRedirectUri,
			ClientSecret: envCfg.AuthClientSecret,
		},
	}
}
