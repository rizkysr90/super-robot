package punishmenthistory

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

type PunishmentHistory struct {
	db                     *sql.DB
	punishmentHistoryStore store.PunishmentHistory
}

func NewHandler(
	routes *gin.RouterGroup,
	db *sql.DB,
	punishmentHistoryStore store.PunishmentHistory,
) {
	handler := &PunishmentHistory{
		db:                     db,
		punishmentHistoryStore: punishmentHistoryStore,
	}
	routes.GET("/", func(ctx *gin.Context) {
		handler.GetAllPunishmentHistoryByNik(ctx)
	})
}

func (p *PunishmentHistory) GetAllPunishmentHistoryByNik(ctx *gin.Context) {
	punishmentHistoryList, err := p.punishmentHistoryStore.FindByID(ctx, ctx.Query("nik"))
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound,
			helper.RestApiError(http.StatusNotFound, constant.ERR_NOT_FOUND, err.Error()))
		return
	}
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError,
			helper.RestApiError(http.StatusInternalServerError,
				constant.ERR_INTERNAL_SERVER_ERROR, err.Error()))
		return
	}
	var res []*httppayload.PunishmentHistoriesGetResponse
	for _, value := range punishmentHistoryList {
		res = append(res, value.ToGetAllPunishmentHistoryByNik())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": res,
	})
}
