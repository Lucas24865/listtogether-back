package controller

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	AdminLogin(ctx *gin.Context)
}

type authController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}

func (r *authController) Register(ctx *gin.Context) {
	var input model.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := r.service.Register(input, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "validated!"})
}

func (r *authController) Login(ctx *gin.Context) {
	var input model.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := r.service.Login(input, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Token": token})
}

func (r *authController) AdminLogin(ctx *gin.Context) {
	var input model.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := r.service.AdminLogin(input, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Token": token})
}
