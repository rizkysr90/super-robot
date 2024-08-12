package main

import (
	"context"
	"log"

	"auth-service-rizkysr90-pos/internal/config"
	"auth-service-rizkysr90-pos/internal/restapi"

	pgx "github.com/rizkysr90/rizkysr90-go-pkg/pgx"
	logger "github.com/rizkysr90/rizkysr90-go-pkg/zerolog"
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
	logger := logger.New().With().
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
	log.Println("HOREE : ",cfg.RestAPIPort)
	restAPIserver, err := restapi.New(cfg, sqlDB, logger)
	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to construct server: %s", err)
	}
	log.Println("HOREE : ",cfg.RestAPIPort)
	err = restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to run server: %s", err)
	}
}
