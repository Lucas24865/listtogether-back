package controller

import (
	"ListTogetherAPI/internal/service"
	"ListTogetherAPI/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListController interface {
	GetAll(ctx *gin.Context)
}

type listController struct {
	listService service.GroupService
}

func NewListController(listService service.ListService) ListController {
	return &listController{
		listService: listService,
	}
}

func (r *listController) GetAll(ctx *gin.Context) {
	user, err := token.ExtractTokenUsername(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lists, err := r.userService.GetAllGroups(user, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": lists})
}
