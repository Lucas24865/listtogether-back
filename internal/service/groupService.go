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
	AddMember(user, group string, ctx *gin.Context) error
	Invite(request requests.GroupRequest, ctx *gin.Context) error
	Get(id, user string, ctx *gin.Context) (*model.Group, error)
	GetGroupsFull(user string, ctx *gin.Context) ([]response.GroupResponse, error)
	GetGroupsFullAdmin(ctx *gin.Context) ([]response.GroupResponse, error)
	Edit(group *model.Group, user string, ctx *gin.Context) error
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

func (s *groupService) Get(id, user string, ctx *gin.Context) (*model.Group, error) {
	groupInfo, err := s.repo.GetGroupSimple(id, ctx)
	if err != nil {
		return nil, err
	}
	if !utils.Contains(groupInfo.Users, user) {
		return nil, nil
	}
	return groupInfo, nil
}

func (s *groupService) GetGroupsFull(user string, ctx *gin.Context) ([]response.GroupResponse, error) {
	groupInfo, err := s.repo.GetGroupsFull(user, ctx)
	if err != nil {
		return nil, err
	}

	return groupInfo, nil
}

func (s *groupService) GetGroupsFullAdmin(ctx *gin.Context) ([]response.GroupResponse, error) {
	groupInfo, err := s.repo.GetGroupsFullAdmin(ctx)
	if err != nil {
		return nil, err
	}

	return groupInfo, nil
}

func (s *groupService) Create(group *model.Group, ctx *gin.Context) (string, error) {
	return s.repo.Create(group, ctx)
}

func (s *groupService) Edit(groupRequest *model.Group, user string, ctx *gin.Context) error {
	groupInfo, err := s.repo.GetGroupSimple(groupRequest.Id, ctx)
	if err != nil {
		return err
	}

	if !utils.Contains(groupInfo.Admins, user) {
		return errors.New("is not an admin of the group")
	}

	for _, u := range groupRequest.Users {
		if utils.Contains(groupInfo.Users, u) {
			continue
		} else {
			message := fmt.Sprintf("%s te ha invitado al grupo: %s", user, groupInfo.Name)
			err = s.notificationService.SendNew(u, groupRequest.Id, message, user, groupInfo.Name, model.GroupInvite, ctx)
			if err != nil {
				return err
			}
		}
	}
	for i, u := range groupInfo.Users {
		if utils.Contains(groupRequest.Users, u) {
			continue
		} else {
			groupInfo.Users = append(groupInfo.Users[:i], groupInfo.Users[i+1:]...)
		}
	}

	return s.repo.Update(groupInfo, ctx)
}

func (s *groupService) AddAdmin(request requests.GroupRequest, ctx *gin.Context) error {
	isAdmin, err := s.IsAdmin(request, ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("is not an admin of the group")
	}

	groupInfo, err := s.repo.GetGroupFull(request.Group, ctx)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s te ha hecho administrador del grupo: %s", request.Admin, groupInfo.Name)
	err = s.notificationService.SendNew(request.User, request.Group, message, request.Admin, groupInfo.Name, model.GroupInvite, ctx)
	if err != nil {
		return err
	}

	return s.repo.AddAdmin(request, ctx)
}

func (s *groupService) AddMember(user, groupId string, ctx *gin.Context) error {
	group, err := s.repo.GetGroupSimple(groupId, ctx)
	if err != nil {
		return err
	}
	group.Users = append(group.Users, user)

	return s.repo.Update(group, ctx)
}

func (s *groupService) Invite(request requests.GroupRequest, ctx *gin.Context) error {
	isAdmin, err := s.IsAdmin(request, ctx)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("is not an admin of the group")
	}

	groupInfo, err := s.repo.GetGroupFull(request.Group, ctx)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s te ha invitado a unirte al grupo: %s", request.Admin, groupInfo.Name)
	err = s.notificationService.SendNew(request.User, request.Group, message, request.Admin, groupInfo.Name, model.GroupInvite, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *groupService) IsAdmin(request requests.GroupRequest, ctx *gin.Context) (bool, error) {
	groupInfo, err := s.repo.GetGroupSimple(request.Group, ctx)
	if err != nil {
		return false, err
	}
	if groupInfo == nil {
		return false, errors.New("invalid group")
	}

	if utils.Contains(groupInfo.Admins, request.Admin) {
		return true, nil
	}

	return false, nil
}
