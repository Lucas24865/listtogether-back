package controller

import (
	"ListTogetherAPI/internal/service"
	"ListTogetherAPI/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotificationController interface {
	GetAll(ctx *gin.Context)
	GetAllWithDeleted(ctx *gin.Context)
	Remove(ctx *gin.Context)
	AcceptInvite(ctx *gin.Context)
	DeclineInvite(ctx *gin.Context)
	RemoveAll(ctx *gin.Context)
}

type notificationController struct {
	notificationService service.NotificationService
	userService         service.UserService
}

func NewNotificationController(notificationService service.NotificationService, userService service.UserService) NotificationController {
	return &notificationController{
		notificationService: notificationService,
		userService:         userService,
	}
}

func (r *notificationController) Remove(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	err = r.notificationService.Remove(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *notificationController) RemoveAll(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.notificationService.RemoveAll(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *notificationController) AcceptInvite(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	err = r.userService.AcceptInvite(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *notificationController) DeclineInvite(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Param("id")
	err = r.userService.DeclineInvite(id, user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *notificationController) GetAll(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notifications, err := r.notificationService.GetAll(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": notifications})
}

func (r *notificationController) GetAllWithDeleted(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notifications, err := r.notificationService.GetAllWithDeleted(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": notifications})
}
