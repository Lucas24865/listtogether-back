package repository

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type GroupRepository interface {
	Create(group *model.Group, ctx *gin.Context) error
	AddAdmin(request requests.GroupRequest, ctx *gin.Context) error
	Deactivate(request requests.GroupRequest, ctx *gin.Context) error
	Update(group *model.Group, ctx *gin.Context) error
	GetGroup(group string, ctx *gin.Context) (*model.Group, error)
}

type groupRepository struct {
	repo *Repository
}

func NewGroupRepository(repo *Repository) GroupRepository {
	return &groupRepository{
		repo: repo,
	}
}

func (g groupRepository) Update(group *model.Group, ctx *gin.Context) error {
	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) Create(group *model.Group, ctx *gin.Context) error {
	group.CreatedAt = time.Now()
	group.Id = uuid.New().String()

	return g.repo.Create("groups", group.Id, group, ctx)
}

func (g groupRepository) AddAdmin(request requests.GroupRequest, ctx *gin.Context) error {
	groupRaw, err := g.repo.GetById("group", request.Group, ctx)
	if err != nil {
		return err
	}

	group := mapGroup(groupRaw)
	group.Admins = append(group.Admins, request.User)

	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) Deactivate(request requests.GroupRequest, ctx *gin.Context) error {
	groupRaw, err := g.repo.GetById("group", request.Group, ctx)
	if err != nil {
		return err
	}

	group := mapGroup(groupRaw)
	group.Admins = append(group.Admins, request.User)

	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) GetGroup(groupName string, ctx *gin.Context) (*model.Group, error) {
	group, err := g.repo.GetById("group", groupName, ctx)
	if err != nil {
		return nil, err
	}
	return mapGroup(group), nil
}

func mapGroup(u map[string]interface{}) *model.Group {
	if u == nil {
		return nil
	}
	group := model.Group{
		Picture:   u["Picture"].(string),
		Desc:      u["Desc"].(string),
		Name:      u["Name"].(string),
		Admins:    u["Admins"].([]string),
		Id:        u["Id"].(string),
		CreatedAt: u["CreatedAt"].(time.Time)}
	return &group
}
