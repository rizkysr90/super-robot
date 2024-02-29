package restapi

import (
	"api-iad-ams/internal/config"
	authHandler "api-iad-ams/internal/restapi/handler/auth"
	auth "api-iad-ams/internal/service/auth"
	"api-iad-ams/internal/store/pg"

	"api-iad-ams/internal/restapi/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(
	cfg *config.Config,
	sqlDB *pgxpool.Pool,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()
	server.Use(middleware.AuthRequired(cfg))
	server.Use(middleware.Recovery())
	server.Use(middleware.ErrorHandler())
	// server.Use(gin.Recovery())
	// Auth service
	authStore := pg.NewAuthDB(sqlDB)
	authService := auth.NewAuthService(authStore)
	authHandler := authHandler.NewAuthHandler(authService, cfg)
	authHandler.AddRoutes(server)
	return server, nil
}
