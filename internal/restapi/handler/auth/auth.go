package auth

import (
	"api-iad-ams/internal/config"
	"api-iad-ams/internal/service"
	"api-iad-ams/pkg/restapierror"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	// db          *pgxpool.Pool
	authService service.AuthService
	config      *config.Config
}

func NewAuthHandler(
	authService service.AuthService,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}
func (a *AuthHandler) AddRoutes(ginEngine *gin.Engine) {
	createUserPath := fmt.Sprintf("%s/auth/users", a.config.ApiVersionBaseURL)
	ginEngine.POST(createUserPath, func(ctx *gin.Context) {
		a.CreateUser(ctx)
	})
}

type testError struct {
	TestArray string `json:"test_array"`
}

func (a *AuthHandler) CreateUser(ctx *gin.Context) {
	// TODO
	// var err error
	// payload := &service.CreateUserRequest{}
	// if err = ctx.Bind(payload); err != nil {
	// 	ctx.jso
	// }
	var payload testError
	if err := ctx.Bind(&payload); err != nil {
		ctx.JSON(500, err)
	}
	// ac := errors.New()
	// if errors.Is(e)
	log.Println("HEREEE", payload.TestArray)
	arr := strings.Split(payload.TestArray, ",")
	// log.Println(arr[3])
	if len(arr) > 3 {
		err := restapierror.NewBadRequest(ctx, restapierror.WithMessage("element array lebih dari 3"))
		// ctx.Error(err)
		ctx.AbortWithStatusJSON(err.Code, err)
	}
	// if len(arr) > 3 {
	// 	ctx.Errors.Errors("testt error")
	// }
	// arr := []int{}
	// log.Println(arr[3])
	// ctx.JSON(200, "success")
	// return nil
}
