package repository

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils"
	"ListTogetherAPI/utils/requests"
	"ListTogetherAPI/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type GroupRepository interface {
	Create(group *model.Group, ctx *gin.Context) (string, error)
	AddAdmin(request requests.GroupRequest, ctx *gin.Context) error
	Deactivate(request requests.GroupRequest, ctx *gin.Context) error
	Update(group *model.Group, ctx *gin.Context) error
	GetGroup(group string, ctx *gin.Context) (*response.GroupResponse, error)
	GetGroups(groupNames []string, ctx *gin.Context) ([]response.GroupResponse, error)
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

func (g groupRepository) Create(group *model.Group, ctx *gin.Context) (string, error) {
	group.CreatedAt = time.Now()
	group.Id = uuid.New().String()
	err := g.repo.Create("groups", group.Id, group, ctx)

	return group.Id, err
}

func (g groupRepository) AddAdmin(request requests.GroupRequest, ctx *gin.Context) error {
	groupRaw, err := g.repo.GetById("groups", request.Group, ctx)
	if err != nil {
		return err
	}

	group := mapGroup(groupRaw)
	group.Admins = append(group.Admins, request.User)

	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) Deactivate(request requests.GroupRequest, ctx *gin.Context) error {
	groupRaw, err := g.repo.GetById("groups", request.Group, ctx)
	if err != nil {
		return err
	}

	group := mapGroup(groupRaw)
	group.Admins = append(group.Admins, request.User)

	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) GetGroup(groupName string, ctx *gin.Context) (*response.GroupResponse, error) {
	group, err := g.repo.GetById("groups", groupName, ctx)
	if err != nil {
		return nil, err
	}
	mappedGroup := mapGroup(group)
	members, err := g.repo.FindAll("users", "Groups", mappedGroup.Id, "array-contains", ctx)
	if err != nil {
		return nil, err
	}
	var mappedMembers []model.User
	var admins []model.User

	for _, member := range members {
		mappedMember := *mapUser(member)
		mappedMembers = append(mappedMembers, mappedMember)
		if utils.Contains(mappedGroup.Admins, mappedMember.User) {
			admins = append(admins, mappedMember)
		}
	}

	groupResponse := response.GroupResponse{
		Id:        mappedGroup.Id,
		Name:      mappedGroup.Name,
		Desc:      mappedGroup.Desc,
		CreatedAt: mappedGroup.CreatedAt,
		Members:   mappedMembers,
		Admins:    admins,
	}

	return &groupResponse, nil
}

func (g groupRepository) GetGroups(groupNames []string, ctx *gin.Context) ([]response.GroupResponse, error) {
	groupsSaved, err := g.repo.FindAll("groups", "Id", groupNames, "in", ctx)
	if err != nil {
		return nil, err
	}

	var groupsResponse []response.GroupResponse

	for _, group := range groupsSaved {
		mappedGroup := mapGroup(group)
		members, err := g.repo.FindAll("users", "Groups", mappedGroup.Id, "array-contains", ctx)
		if err != nil {
			return nil, err
		}
		var mappedMembers []model.User
		var admins []model.User

		for _, member := range members {
			mappedMember := *mapUser(member)
			mappedMembers = append(mappedMembers, mappedMember)
			if utils.Contains(mappedGroup.Admins, mappedMember.User) {
				admins = append(admins, mappedMember)
			}
		}

		groupResponse := response.GroupResponse{
			Id:        mappedGroup.Id,
			Name:      mappedGroup.Name,
			Desc:      mappedGroup.Desc,
			CreatedAt: mappedGroup.CreatedAt,
			Members:   mappedMembers,
			Admins:    admins,
		}

		groupsResponse = append(groupsResponse, groupResponse)

	}
	return groupsResponse, nil
}

func mapGroup(u map[string]interface{}) *model.Group {
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
