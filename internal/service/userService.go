package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type UserService interface {
	GetByUsername(ctx *gin.Context, user string) (*model.User, error)
	GetAllGroups(user string, ctx *gin.Context) ([]*model.Group, error)
	AcceptInvite(id, user string, ctx *gin.Context) error
}

type userService struct {
	repo                repository.UserRepository
	notificationService NotificationService
}

func NewUserService(repo repository.UserRepository, notificationService NotificationService) UserService {
	return &userService{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (s *userService) AcceptInvite(id, user string, ctx *gin.Context) error {
	err := s.notificationService.Accept(id, user, ctx)
	if err != nil {
		return err
	}
	notif, err := s.notificationService.Get(id, user, ctx)
	if err != nil {
		return err
	}
	return s.repo.AddGroup(notif.Data, user, ctx)
}

func (r *userService) GetByUsername(ctx *gin.Context, user string) (*model.User, error) {
	userSaved, err := r.repo.GetByUser(strings.ToLower(user), ctx)
	if err != nil {
		return nil, err
	}
	if userSaved == nil {
		return nil, errors.New("invalid user")
	}

	return userSaved, nil
}

func (r *userService) GetByEmail(ctx *gin.Context, user string) (*model.User, error) {
	userSaved, err := r.repo.GetByUser(strings.ToLower(user), ctx)
	if err != nil {
		return nil, err
	}
	if userSaved == nil {
		return nil, errors.New("invalid user")
	}

	return userSaved, nil
}

func (s *userService) GetAllGroups(user string, ctx *gin.Context) ([]*model.Group, error) {
	group, err := s.repo.GetAllGroups(user, ctx)
	if err != nil {
		return nil, err
	}
	return group, nil
}
