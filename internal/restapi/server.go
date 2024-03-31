package restapi

import (
	"database/sql"
	"net/http"
	"time"

	"auth-service-rizkysr90-pos/internal/config"
	"auth-service-rizkysr90-pos/internal/constant"
	authHandler "auth-service-rizkysr90-pos/internal/restapi/handler/auth"
	employeeHandler "auth-service-rizkysr90-pos/internal/restapi/handler/employee"
	storeHandler "auth-service-rizkysr90-pos/internal/restapi/handler/store"

	"auth-service-rizkysr90-pos/internal/restapi/middleware"
	auth "auth-service-rizkysr90-pos/internal/service/auth"
	"auth-service-rizkysr90-pos/internal/service/employee"
	storeService "auth-service-rizkysr90-pos/internal/service/store"

	"auth-service-rizkysr90-pos/internal/store/pg"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/gin-gonic/gin"
	restapimiddleware "github.com/rizkysr90/rizkysr90-go-pkg/restapi/middleware"
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
	// server.Use(gin.RecoveryWithWriter(logger))
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
	// Employee Service
	employeeStore := pg.NewEmployeeDB(sqlDB)
	employeeService := employee.NewEmployeeService(sqlDB, employeeStore, jwtToken)
	employeeHandler := employeeHandler.NewEmployeeHandler(employeeService)
	server.POST("api/v1/auth/users", func(ctx *gin.Context) {
		authHandler.CreateUser(ctx)
	})
	server.POST("api/v1/auth/users/login", func(ctx *gin.Context) {
		authHandler.LoginUser(ctx)
	})
	server.POST("/api/v1/auth/users/refreshtoken", func(ctx *gin.Context) {
		authHandler.RefreshToken(ctx)
	})
	server.POST("/api/v1/auth/employees/login", func(ctx *gin.Context) {
		employeeHandler.LoginUser(ctx)
	})

	// Store service
	storeStore := pg.NewStoreDB(sqlDB)
	storeService := storeService.NewStoreService(sqlDB, storeStore)
	storeHander := storeHandler.NewAuthHandler(storeService)

	// PRIVATE ROUTES
	authGroup := server.Group("")
	authGroup.Use(middleware.AuthRequiredCookies(jwtToken))
	authGroup.GET("api/v1/privateroutes")

	// EMPLOYEE ENDPOINT
	authGroup.POST("/api/v1/employees", middleware.RBACMiddleware(constant.RBAC_LEVEL_SUPERVISOR), func(ctx *gin.Context) {
		employeeHandler.CreateStore(ctx)
	})
	authGroup.POST("/api/v1/stores", middleware.RBACMiddleware(constant.RBAC_LEVEL_OWNER),
		func(ctx *gin.Context) {
			storeHander.CreateStore(ctx)
		})
	authGroup.GET("/api/v1/stores", func(ctx *gin.Context) {
		storeHander.GetAllStore(ctx)
	})
	server.NoRoute(func(c *gin.Context) {
		c.Error(restapierror.NewNotFound(restapierror.WithMessage("route not found")))
	})
	return server, nil
}
