package restapi

import (
	"database/sql"
	"net/http"

	"rizkysr90-pos/internal/auth"
	"rizkysr90-pos/internal/config"
	categoryHandler "rizkysr90-pos/internal/restapi/handler/category"
	producthandler "rizkysr90-pos/internal/restapi/handler/product"
	"rizkysr90-pos/internal/restapi/middleware"
	categoryService "rizkysr90-pos/internal/service/category"
	"rizkysr90-pos/internal/service/productservice"
	"rizkysr90-pos/internal/store/pg"
	"rizkysr90-pos/internal/utility"
	documentgen "rizkysr90-pos/pkg/documentGen"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/rs/zerolog"
)

func New(
	authClient *auth.Client,
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.New()

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

	server.GET("/oauth", func(ctx *gin.Context) {
		stateID, err := utility.GenerateRandomBase64Str()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		// if err = a.authStore.SetState(c, stateID); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// }
		ctx.Redirect(http.StatusFound, authClient.Oauth.AuthCodeURL(stateID))
	})
	// category service
	categoryStore := pg.NewCategory(sqlDB)
	categoryService := categoryService.NewCategoryService(sqlDB, categoryStore)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)

	// product service
	productStore := pg.NewProduct(sqlDB)
	documentGenerator := documentgen.NewDocumentGenerator()
	productService := productservice.NewProductService(sqlDB, productStore, documentGenerator)
	productHandler := producthandler.NewCategoryHandler(*productService)

	// Create a route group for categories
	categoryRoutes := server.Group("/api/v1/categories")
	{
		categoryRoutes.POST("", categoryHandler.CreateCategory)
		categoryRoutes.GET("", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/:category_id", categoryHandler.GetCategoryByID)
		categoryRoutes.PUT("/:category_id", categoryHandler.EditCategoryByID)
		categoryRoutes.DELETE("/:category_id", categoryHandler.DeleteCategory)
	}
	// Create a route group for products
	productRoutes := server.Group("/api/v1/products")
	{
		productRoutes.POST("", productHandler.CreateProduct)
		productRoutes.PUT("/:product_id", productHandler.UpdateProduct)
		productRoutes.GET("/:product_id", productHandler.GetProductByID)
		productRoutes.GET("", productHandler.GetAllProducts)
		productRoutes.DELETE("/:product_id", productHandler.DeleteProductByID)
		productRoutes.POST("/generate-barcode", productHandler.GenerateBarcodePDF)
	}

	server.NoRoute(func(c *gin.Context) {
		if err := c.Error(errorHandler.NewNotFound(errorHandler.WithMessage("route not found"))); err != nil {
			logger.Error().Msg(err.Error())
		}
	})
	return server, nil
}
