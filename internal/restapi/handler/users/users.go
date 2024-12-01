package users

import (
	"net/http"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/service"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	// db          *pgxpool.Pool
	usersService service.UsersService
	config       config.Config
}

func NewUsersHandler(
	usersService service.UsersService,
	config config.Config,
) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		config:       config,
	}
}

func (u *UsersHandler) LoginUser(ctx *gin.Context) {
	payload := &payload.ReqLoginUsers{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	data, err := u.usersService.LoginUser(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
func (u *UsersHandler) CreateUser(ctx *gin.Context) {
	payload := &payload.ReqCreateUsers{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	data, err := u.usersService.CreateUser(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
