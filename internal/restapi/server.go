package restapi

import (
	"database/sql"

	"auth-service-rizkysr90-pos/internal/config"
	usersHandler "auth-service-rizkysr90-pos/internal/restapi/handler/users"
	"auth-service-rizkysr90-pos/internal/service/users"

	"auth-service-rizkysr90-pos/internal/store/pg"

	"auth-service-rizkysr90-pos/internal/restapi/middleware"

	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
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
	// server.Use(func(ctx *gin.Context) {
	// 	startTime := time.Now()
	// 	ctx.Next()
	// 	// Log request duration and response status
	// 	if ctx.Writer.Status() >= http.StatusBadRequest {
	// 		logger.Error().
	// 			Str("method", ctx.Request.Method).
	// 			Str("url", ctx.Request.URL.String()).
	// 			Str("client_ip", ctx.ClientIP()).
	// 			Str("user_agent", ctx.GetHeader("User-Agent")).
	// 			Dur("duration", time.Since(startTime)).
	// 			Int("status", ctx.Writer.Status()).
	// 			Msg(ctx.Errors[0].Err.Error())
	// 	} else {
	// 		logger.Info().
	// 			Str("method", ctx.Request.Method).
	// 			Str("url", ctx.Request.URL.String()).
	// 			Str("client_ip", ctx.ClientIP()).
	// 			Str("user_agent", ctx.GetHeader("User-Agent")).
	// 			Dur("duration", time.Since(startTime)).
	// 			Int("status", ctx.Writer.Status()).
	// 			Msg("Request processed")
	// 	}
	// })
	// server.Use(restapimiddleware.Recovery(logger))
	server.Use(middleware.RequestBodyMiddleware())
	server.Use(middleware.ResponseBody())
	server.Use(middleware.Recovery(logger))
	// Log middleware
	server.Use(middleware.ErrorHandler(logger))
	// corsHandler := cors.AllowAll()
	server.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	}))

	// Users service
	usersStore := pg.NewUserDB(sqlDB)
	usersService := users.NewUsersService(sqlDB, usersStore, jwtToken)
	user := usersHandler.NewAuthHandler(usersService, cfg)

	server.POST("api/v1/users", func(ctx *gin.Context) {
		user.CreateUser(ctx)
	})
	// server.POST("api/v1/auth/users/login", func(ctx *gin.Context) {
	// 	authHandler.LoginUser(ctx)
	// })
	// server.POST("/api/v1/auth/users/refreshtoken", func(ctx *gin.Context) {
	// 	authHandler.RefreshToken(ctx)
	// })

	// PRIVATE ROUTES
	authGroup := server.Group("")
	authGroup.Use(middleware.AuthRequiredCookies(jwtToken))
	authGroup.GET("api/v1/privateroutes")

	server.NoRoute(func(c *gin.Context) {
		c.Error(restapierror.NewNotFound(restapierror.WithMessage("route not found")))
	})
	return server, nil
}
