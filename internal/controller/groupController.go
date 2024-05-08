package controller

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/service"
	"ListTogetherAPI/utils/requests"
	"ListTogetherAPI/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GroupController interface {
	GetAll(ctx *gin.Context)
	Create(ctx *gin.Context)
	AddAdmin(ctx *gin.Context)
	Remove(ctx *gin.Context)
	Invite(ctx *gin.Context)
}

type groupController struct {
	groupService        service.GroupService
	userService         service.UserService
	notificationService service.NotificationService
}

func NewGroupController(groupService service.GroupService, userService service.UserService,
	notificationService service.NotificationService) GroupController {
	return &groupController{
		groupService:        groupService,
		userService:         userService,
		notificationService: notificationService,
	}
}

func (r *groupController) GetAll(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groups, err := r.userService.GetAllGroups(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": groups})
}

func (r *groupController) Create(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var content model.Group
	if err := ctx.ShouldBindJSON(&content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	admins := make([]string, 0)
	admins = append(admins, user)
	content.Admins = admins
	content.CreatedBy = user

	groupId, err := r.groupService.Create(&content, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.userService.AddGroup(user, groupId, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *groupController) AddAdmin(ctx *gin.Context) {
	admin, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request requests.GroupRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Admin = admin

	err = r.groupService.AddAdmin(request, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (r *groupController) Remove(ctx *gin.Context) {
	/*admin, err := token.ExtractTokenUsername(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := ctx.Param("user")
		group := ctx.Param("group")

		err = r..Remove(ctx, admin, user, group)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	//TODO
		ctx.JSON(http.StatusOK, gin.H{"name": "user"})*/
}

func (r *groupController) Invite(ctx *gin.Context) {
	admin, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request requests.GroupRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.Admin = admin

	err = r.groupService.Invite(request, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
}
