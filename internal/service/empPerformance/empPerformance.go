package empperformance

import (
	"api-iad-ams/internal/constant"
	"api-iad-ams/internal/helper"
	"api-iad-ams/internal/store"
	"api-iad-ams/internal/store/httppayload"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmpPerformanceService struct {
	empPerformanceStore store.EmpPerformance
}

func NewEmpPerformanceService(
	empPerformanceStore store.EmpPerformance,
) *EmpPerformanceService {
	return &EmpPerformanceService{
		empPerformanceStore: empPerformanceStore,
	}
}
func (e *EmpPerformanceService) Find(ctx *gin.Context,
	request *httppayload.EmpPerformanceGetRequest) ([]httppayload.EmpPerformanceGetResponse, error) {

	empPerformanceList, err := e.empPerformanceStore.Find(ctx, request)
	if errors.Is(err, sql.ErrNoRows) {
		// Set the error using c.Error()
		ctx.Error(err)
		// Respond with an error message and status code
		ctx.JSON(http.StatusNotFound,
			helper.RestApiError(http.StatusNotFound, constant.ERR_NOT_FOUND, err.Error()))
		return nil, err
	}
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError,
			helper.RestApiError(http.StatusInternalServerError,
				constant.ERR_INTERNAL_SERVER_ERROR, err.Error()))
		return nil, err
	}
	var res []httppayload.EmpPerformanceGetResponse
	for _, value := range empPerformanceList {
		res = append(res, *value.ToEmpPerformancesGetResponse())
	}
	return res, nil
}
