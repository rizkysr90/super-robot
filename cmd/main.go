package main

import (
	"context"
	"log"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	"github.com/rizkysr90/go-boilerplate/internal/restapi"
	pgx "github.com/rizkysr90/go-boilerplate/pkg/pgx"
	logger "github.com/rizkysr90/go-boilerplate/pkg/zerolog"
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
	sqlDB, sqlDBErr := pgx.NewDB(ctx, cfg.PgSQL)
	// sqlDB.tx
	if sqlDBErr != nil {
		logger.Error().Err(sqlDBErr).Msgf("pgSql: main failed to construct pgSql %s", sqlDBErr)
		return
	}
	defer func() { sqlDB.Close() }()

	restAPIserver, err := restapi.New(cfg, sqlDB, logger)
	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to construct server: %s", err)
	}
	err = restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to run server: %s", err)
	}
}
