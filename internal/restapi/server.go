package restapi

import (
	"database/sql"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	authHandler "github.com/rizkysr90/go-boilerplate/internal/restapi/handler/auth"
	"github.com/rizkysr90/go-boilerplate/internal/restapi/middleware"
	auth "github.com/rizkysr90/go-boilerplate/internal/service/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store/pg"
	jwttoken "github.com/rizkysr90/go-boilerplate/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func New(
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
	jwtToken jwttoken.JWT,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()
	server.Use(middleware.Recovery(logger))
	server.Use(middleware.ErrorHandler(logger))
	// Auth service
	userStore := pg.NewUserDB(sqlDB)
	authService := auth.NewAuthService(sqlDB, userStore, jwtToken)
	authHandler := authHandler.NewAuthHandler(authService, cfg)
	server.POST("api/v1/auth/users", func(ctx *gin.Context) {
		authHandler.CreateUser(ctx)
	})
	server.POST("api/v1/auth/users/login", func(ctx *gin.Context) {
		authHandler.LoginUser(ctx)
	})

	server.Use(middleware.AuthRequired(cfg, &jwtToken))
	server.GET("api/v1/privateroutes")
	return server, nil
}
