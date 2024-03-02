package main

import (
	"api-iad-ams/internal/config"
	"api-iad-ams/internal/restapi"
	pgx "api-iad-ams/pkg/pgx"
	logger "api-iad-ams/pkg/zerolog"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	// -----------------------------------------------------------------------------------------------------------------
	// LOAD APPLICATION CONFIG FROM ENVIRONMENT VARIABLES
	// -----------------------------------------------------------------------------------------------------------------
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("restApi: main failed to load and parse config: %s", err)
		return
	}
	logger := logger.New(cfg).With().
		Str("app", cfg.AppName).
		Str("env", cfg.AppEnv).
		Logger()
	// // -----------------------------------------------------------------------------------------------------------------
	// // INFRASTRUCTURE OBJECTS
	// // -----------------------------------------------------------------------------------------------------------------
	// PgSQL
	sqlDB, sqlDBErr := pgx.NewDB(cfg.PgSQL, ctx)
	// sqlDB.tx
	if sqlDBErr != nil {
		log.Fatalf("pgSql: main failed to construct pgSql %s", sqlDBErr)
		return
	}
	defer func() { sqlDB.Close() }()

	restAPIserver, err := restapi.New(cfg, sqlDB, logger)
	if err != nil {
		log.Fatalf("restapi: main failed to construct server: %s", err)
	}
	restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
}
