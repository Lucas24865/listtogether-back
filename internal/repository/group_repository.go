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
	GetGroupFull(group string, ctx *gin.Context) (*response.GroupResponse, error)
	GetGroupsFull(user string, ctx *gin.Context) ([]response.GroupResponse, error)
	GetGroupSimple(groupId string, ctx *gin.Context) (*model.Group, error)
	GetGroupsSimple(user string, ctx *gin.Context) ([]model.Group, error)
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
	return g.repo.Update("groups", group.Id, *group, ctx)
}

func (g groupRepository) Create(group *model.Group, ctx *gin.Context) (string, error) {
	group.CreatedAt = time.Now()
	group.Id = uuid.New().String()
	group.Deactivated = false
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
	group.Deactivated = true

	return g.repo.Update("groups", group.Id, group, ctx)
}

func (g groupRepository) GetGroupSimple(groupId string, ctx *gin.Context) (*model.Group, error) {
	group, err := g.repo.GetById("groups", groupId, ctx)
	if err != nil {
		return nil, err
	}
	mappedGroup := mapGroup(group)
	return mappedGroup, nil
}

func (g groupRepository) GetGroupsSimple(user string, ctx *gin.Context) ([]model.Group, error) {
	groups, err := g.repo.FindAll("groups", "Users", user, "array-contains", ctx)
	if err != nil {
		return nil, err
	}

	var mappedGroups []model.Group
	for _, group := range groups {
		mappedGroups = append(mappedGroups, *mapGroup(group))
	}

	return mappedGroups, nil
}

func (g groupRepository) GetGroupFull(groupId string, ctx *gin.Context) (*response.GroupResponse, error) {
	group, err := g.repo.GetById("groups", groupId, ctx)
	if err != nil {
		return nil, err
	}
	mappedGroup := mapGroup(group)
	members, err := g.repo.FindAll("users", "User", mappedGroup.Users, "in", ctx)
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

func (g groupRepository) GetGroupsFull(user string, ctx *gin.Context) ([]response.GroupResponse, error) {
	groupsSaved, err := g.repo.FindAll("groups", "Users", user, "array-contains", ctx)
	if err != nil {
		return nil, err
	}

	var groupsResponse []response.GroupResponse

	for _, group := range groupsSaved {
		mappedGroup := mapGroup(group)
		members, err := g.repo.FindAll("users", "User", mappedGroup.Users, "in", ctx)
		if err != nil {
			return nil, err
		}
		var mappedMembers []model.User
		var admins []model.User

		for _, member := range members {
			mappedMember := *mapUser(member)
			if utils.Contains(mappedGroup.Admins, mappedMember.User) {
				admins = append(admins, mappedMember)
			} else {
				mappedMembers = append(mappedMembers, mappedMember)
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

	users := make([]string, 0)
	for _, user := range u["Users"].([]interface{}) {
		users = append(users, user.(string))
	}

	group := model.Group{
		Picture:     u["Picture"].(string),
		Desc:        u["Desc"].(string),
		Name:        u["Name"].(string),
		Admins:      admins,
		Users:       users,
		Deactivated: u["Deactivated"].(bool),
		Id:          u["Id"].(string),
		CreatedAt:   u["CreatedAt"].(time.Time),
		//CreatedBy: u["CreatedBy"].(string),
	}

	return &group
}
