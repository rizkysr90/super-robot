package restapi

import (
	"database/sql"
	"net/http"
	"time"

	"auth-service-rizkysr90-pos/internal/config"
	authHandler "auth-service-rizkysr90-pos/internal/restapi/handler/auth"
	"auth-service-rizkysr90-pos/internal/restapi/middleware"
	auth "auth-service-rizkysr90-pos/internal/service/auth"
	"auth-service-rizkysr90-pos/internal/store/pg"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/gin-gonic/gin"
	restapimiddleware "github.com/rizkysr90/rizkysr90-go-pkg/restapi/middleware"
	cors "github.com/rs/cors/wrapper/gin"
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
	server.Use(func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		// Log request duration and response status
		if ctx.Writer.Status() >= http.StatusBadRequest {
			logger.Error().
				Str("method", ctx.Request.Method).
				Str("url", ctx.Request.URL.String()).
				Str("client_ip", ctx.ClientIP()).
				Str("user_agent", ctx.GetHeader("User-Agent")).
				Dur("duration", time.Since(startTime)).
				Int("status", ctx.Writer.Status()).
				Msg(ctx.Errors[0].Err.Error())
		} else {
			logger.Info().
				Str("method", ctx.Request.Method).
				Str("url", ctx.Request.URL.String()).
				Str("client_ip", ctx.ClientIP()).
				Str("user_agent", ctx.GetHeader("User-Agent")).
				Dur("duration", time.Since(startTime)).
				Int("status", ctx.Writer.Status()).
				Msg("Request processed")
		}
	})
	server.Use(restapimiddleware.Recovery(logger))
	server.Use(restapimiddleware.ErrorHandler())
	// corsHandler := cors.AllowAll()
	server.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	}))

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
	server.POST("/api/v1/auth/users/refreshtoken", func(ctx *gin.Context) {
		authHandler.RefreshToken(ctx)
	})
	authGroup := server.Group("")
	authGroup.Use(middleware.AuthRequiredCookies(jwtToken))
	authGroup.GET("api/v1/privateroutes")
	server.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	// server.Use(func(ctx *gin.Context) {
	// 	log.Println(ctx)
	// 	logger.Info().
	// 		Str("path", ctx.FullPath()).
	// 		Msg("incoming process")
	// })
	return server, nil
}
