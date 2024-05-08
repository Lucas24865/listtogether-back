package repository

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils/response"
	"github.com/gin-gonic/gin"
	"time"
)

type ListRepository interface {
	GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error)
}

type listRepository struct {
	repo *Repository
}

func NewListRepository(repo *Repository) ListRepository {
	return &listRepository{
		repo: repo,
	}
}

func (g listRepository) GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error) {

	list := g.repo.FindAll("users", "Groups", mappedGroup.Id, "array-contains", ctx)

	//TODO buscar todas los grupos de un user, buscar todas las listas de esos grupos y completar con los usuarios las respuestas

	return g.repo.Update("groups", group.Id, group, ctx)
}

func mapList(u map[string]interface{}) *model.Group {
	if u == nil {
		return nil
	}
	admins := make([]string, 0)
	for _, admin := range u["Admins"].([]interface{}) {
		admins = append(admins, admin.(string))
	}
	group := model.Group{
		Picture:   u["Picture"].(string),
		Desc:      u["Desc"].(string),
		Name:      u["Name"].(string),
		Admins:    admins,
		Id:        u["Id"].(string),
		CreatedAt: u["CreatedAt"].(time.Time),
		//CreatedBy: u["CreatedBy"].(string),
	}

	return &group
}
