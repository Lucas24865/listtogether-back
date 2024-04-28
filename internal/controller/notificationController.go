package controller

import (
	"ListTogetherAPI/internal/service"
	"ListTogetherAPI/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotificationController interface {
	GetAll(ctx *gin.Context)
	Remove(ctx *gin.Context)
	Accept(ctx *gin.Context)
	Decline(ctx *gin.Context)
}

type notificationController struct {
	service service.NotificationService
}

func NewNotificationController(service service.NotificationService) NotificationController {
	return &notificationController{
		service: service,
	}
}

func (r *notificationController) Remove(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	err = r.service.Remove(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (r *notificationController) Accept(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	err = r.service.Accept(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (r *notificationController) Decline(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	err = r.service.Decline(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (r *notificationController) GetAll(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notifications, err := r.service.GetAll(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": notifications})
}
