package main

import (
	"api-iad-ams/internal/config"
	"api-iad-ams/internal/server"
	"api-iad-ams/pkg/mysql"
	"log"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// LOAD APPLICATION CONFIG FROM ENVIRONMENT VARIABLES
	// -----------------------------------------------------------------------------------------------------------------
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("restapi: main failed to load and parse config: %s", err)
		return
	}
	// -----------------------------------------------------------------------------------------------------------------
	// INFRASTRUCTURE OBJECTS
	// -----------------------------------------------------------------------------------------------------------------
	// MYSQL
	sqlDB, sqlDBErr := mysql.NewDB(cfg.MySQL)
	if sqlDBErr != nil {
		log.Fatalf("restapi: main failed to construct mysql %s", sqlDBErr)
		return
	}
	defer func() { _ = sqlDB.Close() }()
	restAPIserver, err := server.New(cfg, sqlDB)
	if err != nil {
		log.Fatalf("restapi: main failed to construct server: %s", err)
	}
	restAPIserver.Run(cfg.RestAPIPort) // listen and serve on 0.0.0.0:8080
}
