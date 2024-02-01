package service

import (
	"api-iad-ams/internal/store/httppayload"

	"github.com/gin-gonic/gin"
)

type PunishmentHistoriesServiceAbstraction interface {
	Find(ctx *gin.Context,
		request *httppayload.PunishmentHistoriesGetRequest) ([]httppayload.PunishmentHistoriesGetResponse, error)
}
