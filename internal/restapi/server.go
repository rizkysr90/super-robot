package restapi

import (
	"database/sql"

	"auth-service-rizkysr90-pos/internal/config"
	authHandler "auth-service-rizkysr90-pos/internal/restapi/handler/auth"
	"auth-service-rizkysr90-pos/internal/restapi/middleware"
	auth "auth-service-rizkysr90-pos/internal/service/auth"
	"auth-service-rizkysr90-pos/internal/store/pg"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/gin-gonic/gin"
	restapimiddleware "github.com/rizkysr90/rizkysr90-go-pkg/restapi/middleware"
	"github.com/rs/zerolog"
)

func New(
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
	jwtToken *jwttoken.JWT,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()
	server.Use(restapimiddleware.Recovery(logger))
	server.Use(restapimiddleware.ErrorHandler(logger))
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

	server.Use(middleware.AuthRequired(cfg, jwtToken))
	server.GET("api/v1/privateroutes")
	return server, nil
}
