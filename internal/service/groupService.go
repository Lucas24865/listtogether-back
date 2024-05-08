package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils"
	"ListTogetherAPI/utils/requests"
	"ListTogetherAPI/utils/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type GroupService interface {
	Create(group *model.Group, ctx *gin.Context) (string, error)
	AddAdmin(request requests.GroupRequest, ctx *gin.Context) error
	RemoveMember(request requests.GroupRequest, ctx *gin.Context) error
	Invite(request requests.GroupRequest, ctx *gin.Context) error
	Get(id string, ctx *gin.Context) (*response.GroupResponse, error)
	GetFullGroups(ids []string, ctx *gin.Context) ([]response.GroupResponse, error)
}

type groupService struct {
	repo                repository.GroupRepository
	notificationService NotificationService
}

func NewGroupService(repo repository.GroupRepository, notificationService NotificationService) GroupService {
	return &groupService{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (s *groupService) Get(id string, ctx *gin.Context) (*response.GroupResponse, error) {
	groupInfo, err := s.repo.GetGroup(id, ctx)
	if err != nil {
		return nil, err
	}
	return groupInfo, nil
}

func (s *groupService) GetFullGroups(ids []string, ctx *gin.Context) ([]response.GroupResponse, error) {
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}

	groupInfo, err := s.repo.GetGroups(ids, ctx)
	if err != nil {
		return nil, err
	}

	return groupInfo, nil
}

func (s *groupService) Create(group *model.Group, ctx *gin.Context) (string, error) {
	return s.repo.Create(group, ctx)
}

func (s *groupService) AddAdmin(request requests.GroupRequest, ctx *gin.Context) error {
	isAdmin, err := s.IsAdmin(request, ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("is not an admin of the group")
	}

	groupInfo, err := s.repo.GetGroup(request.Group, ctx)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s te ha hecho administrador del grupo: %s", request.Admin, groupInfo.Name)
	err = s.notificationService.SendNew(request.User, request.Group, message, ctx)
	if err != nil {
		return err
	}

	return s.repo.AddAdmin(request, ctx)
}

func (s *groupService) RemoveMember(request requests.GroupRequest, ctx *gin.Context) error {
	isAdmin, err := s.IsAdmin(request, ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("admin is not an admin of the group")
	}
	//TODO removeMember
	return nil
	//return s.repo.RemoveMember(request, ctx)
}

func (s *groupService) Invite(request requests.GroupRequest, ctx *gin.Context) error {
	isAdmin, err := s.IsAdmin(request, ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("is not an admin of the group")
	}

	groupInfo, err := s.repo.GetGroup(request.Group, ctx)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s te ha invitado a unirte al grupo: %s", request.Admin, groupInfo.Name)
	err = s.notificationService.SendNew(request.User, request.Group, message, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *groupService) IsAdmin(request requests.GroupRequest, ctx *gin.Context) (bool, error) {
	groupInfo, err := s.repo.GetGroup(request.Group, ctx)
	if err != nil {
		return false, err
	}
	if groupInfo == nil {
		return false, errors.New("invalid group")
	}

	for _, groupAdmin := range groupInfo.Admins {
		if utils.Contains(groupAdmin.Groups, request.Admin) {
			return true, nil
		}
	}

	return false, nil
}
