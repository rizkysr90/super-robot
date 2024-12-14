//nolint:godot
package main

import (
	"context"
	"log"

	"rizkysr90-pos/docs" // This is where Swag will generate its docs.go file
	"rizkysr90-pos/internal/auth"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/restapi"

	"github.com/redis/go-redis/v9"
	pgx "github.com/rizkysr90/rizkysr90-go-pkg/pgx"
	logger "github.com/rizkysr90/rizkysr90-go-pkg/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Point of Sale API
// @version 1.0
// @description API for Point of Sale System
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
	// RedisClient init
	redisClient := redis.NewClient(&cfg.RedisConfig)
	if err = redisClient.Ping(ctx).Err(); err != nil {
		logger.Error().Err(sqlDBErr).Msgf("redis: main failed to construct redis %s", err.Error())
		return
	}
	defer func() { sqlDB.Close() }()

	authClient, err := auth.New(ctx, &auth.Config{
		BaseURL:      cfg.Auth.BaseURL,
		ClientID:     cfg.Auth.ClientID,
		RedirectURI:  cfg.Auth.RedirectURI,
		ClientSecret: cfg.Auth.ClientSecret,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("restapi:  main failed to create auth client: %s", err)
		return
	}
	restAPIserver, err := restapi.New(authClient, cfg, sqlDB, logger, redisClient)
	docs.SwaggerInfo.BasePath = "/api/v1"
	restAPIserver.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err != nil {
		logger.Error().Err(err).Msgf("restapi: main failed to construct server: %s", err)
	}
	err = restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
	if err != nil {
		logger.Error().Err(err).Msgf("restapi: main failed to run server: %s", err)
	}
}
