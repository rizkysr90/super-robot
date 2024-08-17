package restapi

import (
	"database/sql"

	"auth-service-rizkysr90-pos/internal/config"
	categoryHandler "auth-service-rizkysr90-pos/internal/restapi/handler/category"
	usersHandler "auth-service-rizkysr90-pos/internal/restapi/handler/users"

	categoryService "auth-service-rizkysr90-pos/internal/service/category"
	usersService "auth-service-rizkysr90-pos/internal/service/user"

	"auth-service-rizkysr90-pos/internal/store/pg"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"auth-service-rizkysr90-pos/internal/restapi/middleware"

	"auth-service-rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/rs/zerolog"
)

func New(
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()
	jwt := jwttoken.New(cfg.SecretKeyJWT)
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
	server.Use(middleware.Recovery(logger))
	// Log middleware
	server.Use(middleware.LogMiddleware(logger))
	server.Use(middleware.ResponseBody())
	server.Use(middleware.ErrorHandler(logger))
	server.Use(middleware.RequestBodyMiddleware())
	// corsHandler := cors.AllowAll()
	server.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	}))

	// Users service
	usersStore := pg.NewUserDB(sqlDB)
	usersService := usersService.NewUsersService(sqlDB, usersStore, jwt)
	usersHandler := usersHandler.NewUsersHandler(usersService, cfg)

	// category service
	categoryStore := pg.NewCategory(sqlDB)
	categoryService := categoryService.NewCategoryService(sqlDB, categoryStore)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)

	server.POST("api/v1/login/users", func(ctx *gin.Context) {
		usersHandler.LoginUser(ctx)
	})
	server.POST("api/v1/signup/users", func(ctx *gin.Context) {
		usersHandler.CreateUser(ctx)
	})
	server.POST("/api/v1/categories", func(ctx *gin.Context) {
		categoryHandler.CreateCategory(ctx)
	})
	// server.POST("api/v1/auth/users/login", func(ctx *gin.Context) {
	// 	authHandler.LoginUser(ctx)
	// })
	// server.POST("/api/v1/auth/users/refreshtoken", func(ctx *gin.Context) {
	// 	authHandler.RefreshToken(ctx)
	// })

	// PRIVATE ROUTES
	authGroup := server.Group("")
	authGroup.Use(middleware.AuthRequiredCookies(jwt))
	authGroup.GET("api/v1/privateroutes")

	server.NoRoute(func(c *gin.Context) {
		c.Error(errorHandler.NewNotFound(errorHandler.WithMessage("route not found")))
	})
	return server, nil
}
