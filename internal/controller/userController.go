package controller

import (
	"ListTogetherAPI/internal/service"
	"ListTogetherAPI/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	Get(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userController{
		service: service,
	}
}

func (r *userController) Get(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = r.service.GetByUsername(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": user})
}