package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/requests"
	"errors"
	"github.com/gin-gonic/gin"
)

type GroupService interface {
	Create(group *model.Group, ctx *gin.Context) error
	AddAdmin(request requests.GroupRequest, ctx *gin.Context) error
	RemoveMember(request requests.GroupRequest, ctx *gin.Context) error
	Invite(request requests.GroupRequest, ctx *gin.Context) error
	AddMember(request requests.GroupRequest, ctx *gin.Context) error
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{
		repo: repo,
	}
}

func (s *groupService) Create(group *model.Group, ctx *gin.Context) error {
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
		return errors.New("admin is not an admin of the group")
	}
	//TODO InviteMember
	return nil
	//return s.repo.Invite(ctx, admin, user, group)
}

func (s *groupService) AddMember(request requests.GroupRequest, ctx *gin.Context) error {
	//TODO ADD MEMBER
	return nil
}

func (s *groupService) IsAdmin(request requests.GroupRequest, ctx *gin.Context) (bool, error) {
	groupInfo, err := s.repo.GetGroup(request.Group, ctx)
	if err != nil {
		return false, err
	}

	for _, groupAdmin := range groupInfo.Admins {
		if groupAdmin == request.Admin {
			return true, nil
		}
	}

	return false, nil
}
