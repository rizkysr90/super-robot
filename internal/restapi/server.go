package restapi

import (
	"database/sql"

	"rizkysr90-pos/internal/config"
	categoryHandler "rizkysr90-pos/internal/restapi/handler/category"
	producthandler "rizkysr90-pos/internal/restapi/handler/product"
	usersHandler "rizkysr90-pos/internal/restapi/handler/users"
	"rizkysr90-pos/internal/restapi/middleware"
	categoryService "rizkysr90-pos/internal/service/category"
	"rizkysr90-pos/internal/service/productservice"
	usersService "rizkysr90-pos/internal/service/user"
	"rizkysr90-pos/internal/store/pg"
	"rizkysr90-pos/pkg/errorHandler"
	jwttoken "rizkysr90-pos/pkg/jwt"

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

	server.Use(middleware.Recovery(logger))
	server.Use(middleware.LogMiddleware(logger))
	server.Use(middleware.ResponseBody())
	server.Use(middleware.ErrorHandler(logger))
	server.Use(middleware.RequestBodyMiddleware())
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

	// product service
	productStore := pg.NewProduct(sqlDB)
	productService := productservice.NewProductService(sqlDB, productStore)
	productHandler := producthandler.NewCategoryHandler(*productService)
	server.POST("api/v1/login/users", func(ctx *gin.Context) {
		usersHandler.LoginUser(ctx)
	})
	server.POST("api/v1/signup/users", func(ctx *gin.Context) {
		usersHandler.CreateUser(ctx)
	})
	server.POST("/api/v1/categories", func(ctx *gin.Context) {
		categoryHandler.CreateCategory(ctx)
	})
	server.GET("/api/v1/categories", func(ctx *gin.Context) {
		categoryHandler.GetAllCategories(ctx)
	})
	server.GET("/api/v1/categories/:category_id", func(ctx *gin.Context) {
		categoryHandler.GetCategoryByID(ctx)
	})
	server.PUT("/api/v1/categories/:category_id", func(ctx *gin.Context) {
		categoryHandler.EditCategoryByID(ctx)
	})
	server.DELETE("/api/v1/categories/:category_id", func(ctx *gin.Context) {
		categoryHandler.DeleteCategory(ctx)
	})
	server.POST("/api/v1/products", func(ctx *gin.Context) {
		productHandler.CreateProduct(ctx)
	})
	server.PUT("/api/v1/products/:product_id", func(ctx *gin.Context) {
		productHandler.UpdateProduct(ctx)
	})
	server.GET("/api/v1/products/:product_id", func(ctx *gin.Context) {
		productHandler.GetProductByID(ctx)
	})
	server.GET("/api/v1/products", func(ctx *gin.Context) {
		productHandler.GetAllProducts(ctx)
	})
	server.DELETE("/api/v1/products/:product_id", func(ctx *gin.Context) {
		productHandler.DeleteProductByID(ctx)
	})
	server.POST("/api/v1/products/generate-barcode", func(ctx *gin.Context) {
		productHandler.GenerateBarcodePDF(ctx)
	})
	// server.POST("api/v1/auth/users/login", func(ctx *gin.Context) {
	// 	authHandler.LoginUser(ctx)
	// })
	// server.POST("/api/v1/auth/users/refreshtoken", func(ctx *gin.Context) {
	// 	authHandler.RefreshToken(ctx)
	// })

	// PRIVATE ROUTES
	authGroup := server.Group("")
	authGroup.GET("api/v1/privateroutes")

	server.NoRoute(func(c *gin.Context) {
		if err := c.Error(errorHandler.NewNotFound(errorHandler.WithMessage("route not found"))); err != nil {
			logger.Error().Msg(err.Error())
		}
	})
	return server, nil
}
