package server

import (
	"api-iad-ams/internal/config"
	empperformance "api-iad-ams/internal/server/handler/empPerformance"
	empperformanceService "api-iad-ams/internal/service/empPerformance"

	"api-iad-ams/internal/server/handler/helloworld"
	punishmenthistory "api-iad-ams/internal/server/handler/punishmentHistory"
	"api-iad-ams/internal/server/middleware"
	"api-iad-ams/internal/store/mysql"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func New(
	cfg *config.Config,
	sqlDB *sql.DB,
) (*gin.Engine, error) {
	// Setup rest api server and its provided services.
	server := gin.Default()
	server.Use(middleware.AuthRequired(cfg))
	// Employee Performances Services
	empPerformanceAPIPath := fmt.Sprintf("/%s/employee_performances", cfg.ApiVersionBaseURL)
	empPerformanceRoutes := server.Group(empPerformanceAPIPath)
	empPerformanceStore := mysql.NewEmpPerformance(sqlDB)
	empPerformanceService := empperformanceService.NewEmpPerformanceService(empPerformanceStore)
	empperformance.NewHandler(empPerformanceRoutes, empPerformanceService)
	// Punishment Hisotry Services
	punishmentHistoryAPIPath := fmt.Sprintf("/%s/punishment_histories", cfg.ApiVersionBaseURL)
	punishmentHistoryRoutes := server.Group(punishmentHistoryAPIPath)
	punishmentHistoryStore := mysql.NewPunishmentHistory(sqlDB)
	punishmenthistory.NewHandler(punishmentHistoryRoutes, sqlDB, punishmentHistoryStore)
	// Hello World Service Example
	helloWorldAPIPath := fmt.Sprintf("/%s/helloworld", cfg.ApiVersionBaseURL)
	helloWorldRoutes := server.Group(helloWorldAPIPath)
	helloworld.NewHandler(helloWorldRoutes, sqlDB, empPerformanceStore)

	return server, nil
}
