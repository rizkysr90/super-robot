package helloworld

import (
	"api-iad-ams/internal/store"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloWorldHandler struct {
	db                  *sql.DB
	empPerformanceStore store.EmpPerformance
}

func NewHandler(
	routes *gin.RouterGroup,
	db *sql.DB,
	empPerformanceStore store.EmpPerformance,
) {
	handler := HelloWorldHandler{
		db:                  db,
		empPerformanceStore: empPerformanceStore,
	}
	routes.GET("/", func(ctx *gin.Context) {
		handler.GetHelloWorld(ctx)
	})
}
func (h *HelloWorldHandler) GetHelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Hello world from GIN",
	})
}
