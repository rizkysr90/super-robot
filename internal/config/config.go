package config

import (
	"api-iad-ams/pkg/mysql"

	"github.com/caarlos0/env/v8"
)

type flatEnv struct {
	AppName           string `env:"APP_NAME"`
	AppEnv            string `env:"APP_ENV"` // local|dev|uat|sit|prod
	RestAPIPort       string `env:"REST_API_PORT"`
	MySQLHost         string `env:"MYSQL_HOST"`
	MySQLPort         int    `env:"MYSQL_PORT"`
	MySQLDatabase     string `env:"MYSQL_DATABASE"`
	MySQLUsername     string `env:"MYSQL_USERNAME"`
	MySQLPassword     string `env:"MYSQL_PASSWORD,unset"` // sensitive, unset from environment after parse!
	MySQLConnMaxOpen  int    `env:"MYSQL_CONN_MAX_OPEN"`
	MySQLConnMaxIdle  int    `env:"MYSQL_CONN_MAX_IDLE"`
	ApiKey            string `env:"API_KEY,unset"`
	ApiVersionBaseURL string `env:"API_VERSION_BASE_URL"`
}
type Config struct {
	AppName           string
	AppEnv            string
	RestAPIPort       string
	MySQL             mysql.Config
	MySQLHost         string
	MySQLPort         int
	MySQLDatabase     string
	MySQLUsername     string
	MySQLPassword     string
	MySQLConnMaxOpen  int
	MySQLConnMaxIdle  int
	APIKey            string
	ApiVersionBaseURL string
}

func LoadFromEnv() (*Config, error) {
	var envCfg flatEnv
	err := env.Parse(&envCfg)
	if err != nil {
		return nil, err
	}
	return newConfig(envCfg), nil
}
func newConfig(envCfg flatEnv) *Config {
	return &Config{
		AppName:     envCfg.AppName,
		AppEnv:      envCfg.AppEnv,
		RestAPIPort: envCfg.RestAPIPort,
		MySQL: mysql.Config{
			Username:    envCfg.MySQLUsername,
			Password:    envCfg.MySQLPassword,
			Database:    envCfg.MySQLDatabase,
			Host:        envCfg.MySQLHost,
			Port:        envCfg.MySQLPort,
			ConnMaxOpen: envCfg.MySQLConnMaxOpen,
			ConnMaxIdle: envCfg.MySQLConnMaxIdle,
		},
		APIKey:            envCfg.ApiKey,
		ApiVersionBaseURL: envCfg.ApiVersionBaseURL,
	}
}
