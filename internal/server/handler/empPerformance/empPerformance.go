package empperformance

import (
	"api-iad-ams/internal/service"
	"api-iad-ams/internal/store/httppayload"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmpPerformance struct {
	empPerformanceService service.EmpPerformanceServiceAbstraction
}

func NewHandler(
	routes *gin.RouterGroup,
	empPerformanceService service.EmpPerformanceServiceAbstraction,
) {
	handler := &EmpPerformance{
		empPerformanceService: empPerformanceService,
	}
	routes.GET("/", func(ctx *gin.Context) {
		handler.GetEmpPerformances(ctx)
	})
}

func (e *EmpPerformance) GetEmpPerformances(ctx *gin.Context) {
	request := httppayload.EmpPerformanceGetRequest{
		Nik: ctx.Query("nik"),
	}
	empPerformancesList, err := e.empPerformanceService.Find(ctx, &request)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": empPerformancesList,
	})
}
