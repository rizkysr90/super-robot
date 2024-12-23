//nolint:godot
package main

import (
	"context"
	"log"

	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/docs" // This is where Swag will generate its docs.go file
	"rizkysr90-pos/internal/restapi"

	pgx "github.com/rizkysr90/rizkysr90-go-pkg/pgx"
	logger "github.com/rizkysr90/rizkysr90-go-pkg/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Rizki Plastik API
// @version 1.0
// @description rizki plastik point of sale api server.
// @host localhost:8080
// @BasePath /api/v1
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
	restAPIserver, err := restapi.New(cfg, sqlDB, logger)
	docs.SwaggerInfo.BasePath = "/api/v1"
	restAPIserver.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to construct server: %s", err)
	}
	err = restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		logger.Error().Err(sqlDBErr).Msgf("restapi: main failed to run server: %s", err)
	}
}
