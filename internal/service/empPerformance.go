package service

import (
	"api-iad-ams/internal/store/httppayload"

	"github.com/gin-gonic/gin"
)

type EmpPerformanceServiceAbstraction interface {
	Find(ctx *gin.Context, request *httppayload.EmpPerformanceGetRequest) ([]httppayload.EmpPerformanceGetResponse, error)
}
