package restapi

import (
	"database/sql"

	"rizkysr90-pos/internal/auth"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/restapi/handler"
	categoryHandler "rizkysr90-pos/internal/restapi/handler/category"
	producthandler "rizkysr90-pos/internal/restapi/handler/product"
	"rizkysr90-pos/internal/restapi/middleware"
	"rizkysr90-pos/internal/service/admin"
	authService "rizkysr90-pos/internal/service/auth"
	"rizkysr90-pos/internal/service/branches"
	categoryService "rizkysr90-pos/internal/service/category"
	"rizkysr90-pos/internal/service/productservice"
	"rizkysr90-pos/internal/service/roles"

	"rizkysr90-pos/internal/store/pg"
	rds "rizkysr90-pos/internal/store/redis"
	documentgen "rizkysr90-pos/pkg/documentGen"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/rs/zerolog"
)

func New(
	authClient *auth.Client,
	cfg config.Config,
	sqlDB *sql.DB,
	logger zerolog.Logger,
	redis *redis.Client,
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

	// category service
	categoryStore := pg.NewCategory(sqlDB)
	categoryService := categoryService.NewCategoryService(sqlDB, categoryStore)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryService)

	// product service
	productStore := pg.NewProduct(sqlDB)
	documentGenerator := documentgen.NewDocumentGenerator()
	productService := productservice.NewProductService(sqlDB, productStore, documentGenerator)
	productHandler := producthandler.NewCategoryHandler(*productService)

	// Auth service
	authStateStore := pg.NewState(sqlDB)
	tenantStore := pg.NewTenant(sqlDB)
	userStore := pg.NewUser(sqlDB)
	sessionStore := rds.NewSessionRedisManager(redis)
	authService := authService.NewAuthService(
		sqlDB, authClient, authStateStore,
		userStore, tenantStore, sessionStore,
	)
	authHandler := handler.NewAuthHandler(authClient, authService)

	// branch service
	branchStore := pg.NewBranches(sqlDB)
	branchesService := branches.NewBranchService(sqlDB, &cfg, tenantStore, userStore, branchStore)
	branchHandler := handler.NewBranchHandler(branchesService)

	// roles service
	assignmentRoleStore := pg.NewAssignmentRoles(sqlDB)
	tenantPermissionStore := pg.NewTenantPermission(sqlDB)
	workLocationStore := pg.NewWorkLocation(sqlDB)
	tenantRoleStore := pg.NewTenantRole(sqlDB)
	rolesService := roles.NewService(
		sqlDB,
		&cfg,
		tenantStore,
		userStore,
		branchStore,
		assignmentRoleStore,
		tenantPermissionStore,
		workLocationStore,
		tenantRoleStore,
	)
	adminService := admin.NewService(sqlDB, &cfg, tenantStore,
		userStore,
		branchStore,
		assignmentRoleStore,
		tenantPermissionStore,
		workLocationStore,
		tenantRoleStore)

	rolesHandler := handler.NewRolesHandler(rolesService)
	userAdminHandler := handler.NewUserAdmin(adminService)
	// server.GET("/oauth", func(ctx *gin.Context) {
	// 	authClient.HandlerRedirect(ctx, sqlDB, authStateStore)
	// })
	server.GET("/callback", func(ctx *gin.Context) {
		authHandler.Callback(ctx)
	})
	// create a route group for auth
	authRoutes := server.Group("/api/v1/auth")
	{
		authRoutes.GET("/register/owner", authHandler.OwnerRegistration)
		authRoutes.GET("/login/owner", authHandler.OwnerLogin)
	}
	branchRoutes := server.Group("/api/v1/branches")
	{
		branchRoutes.POST("/", branchHandler.Create)
		branchRoutes.GET("/", branchHandler.GetBranches)
	}
	rolesRoutes := server.Group("/api/v1/roles")
	{
		rolesRoutes.POST("/", rolesHandler.Create)
	}
	userAdminRoutes := server.Group("/api/v1/users")
	{
		userAdminRoutes.POST("/", userAdminHandler.Create)
	}
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
