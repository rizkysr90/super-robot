package restapi

import (
	"database/sql"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	authHandler "github.com/rizkysr90/go-boilerplate/internal/restapi/handler/auth"
	"github.com/rizkysr90/go-boilerplate/internal/restapi/middleware"
	auth "github.com/rizkysr90/go-boilerplate/internal/service/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store/pg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func New(
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()
	server.Use(middleware.Recovery(logger))
	server.Use(middleware.ErrorHandler(logger))

	// Auth service
	userStore := pg.NewUserDB(sqlDB)
	authService := auth.NewAuthService(sqlDB, userStore)
	// server.Use(middleware.AuthRequired(cfg))
	authHandler := authHandler.NewAuthHandler(authService, cfg)

	authHandler.AddRoutes(server)
	return server, nil
}
