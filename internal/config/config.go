package config

import (
	"fmt"
	"log"
	"rizkysr90-pos/internal/auth"
	"strconv"

	"github.com/redis/go-redis/v9"
	pgx "github.com/rizkysr90/rizkysr90-go-pkg/pgx"

	"github.com/caarlos0/env/v8"
)

type flatEnv struct {
	PrivateKeyJWT     string `env:"PRIVATE_KEY_JWT,unset"`
	AuthURI           string `env:"AUTH_URI"`
	APIVersionBaseURL string `env:"API_VERSION_BASE_URL"`
	DBHost            string `env:"DB_HOST"`
	RedisDatabase     string `env:"REDIS_DATABASE"`
	DBUsername        string `env:"DB_USERNAME"`
	AppEnv            string `env:"APP_ENV"`
	RedisPort         string `env:"REDIS_PORT"`
	DBDatabase        string `env:"DB_DATABASE"`
	LogLevel          string `env:"LOG_LEVEL"`
	SecretKeyJWT      string `env:"SECRET_KEY_JWT,unset"`
	PublicKeyJWT      string `env:"PUBLIC_KEY_JWT,unset"`
	RestAPIPort       string `env:"REST_API_PORT"`
	AppName           string `env:"APP_NAME"`
	APIKey            string `env:"API_KEY,unset"`
	RedisHost         string `env:"REDIS_HOST"`
	AuthClientID      string `env:"AUTH_CLIENT_ID"`
	AuthRedirectURI   string `env:"AUTH_REDIRECT_URI"`
	AuthClientSecret  string `env:"AUTH_CLIENT_SECRET"`
	DBPassword        string `env:"DB_PASSWORD,unset"`
	RedisUsername     string `env:"REDIS_USERNAME"`
	RedisPassword     string `env:"REDIS_PASSWORD"`
	DBPort            int    `env:"DB_PORT"`
	DBConnMaxOpen     int    `env:"DB_CONN_MAX_OPEN"`
	DBConnMaxIdle     int    `env:"DB_CONN_MAX_IDLE"`
}
type Config struct {
	Auth              *auth.Config
	AppName           string
	AppEnv            string
	RestAPIPort       string
	APIKey            string
	APIVersionBaseURL string
	LogLevel          string
	SecretKeyJWT      string
	RedisConfig       redis.Options
	PgSQL             pgx.Config
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
	redisDB, err := strconv.Atoi(envCfg.RedisDatabase)
	if err != nil {
		log.Fatal("Database redis invalid : ", err)
	}
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
			BaseURL:      envCfg.AuthURI,
			ClientID:     envCfg.AuthClientID,
			RedirectURI:  envCfg.AuthRedirectURI,
			ClientSecret: envCfg.AuthClientSecret,
		},
		RedisConfig: redis.Options{
			Addr:     fmt.Sprintf("%s:%s", envCfg.RedisHost, envCfg.RedisPort),
			Username: envCfg.RedisUsername,
			Password: envCfg.RedisPassword,
			DB:       redisDB,
		},
	}
}
