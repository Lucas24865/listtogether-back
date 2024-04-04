package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetByUsername(ctx *gin.Context, user string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (r *userService) GetByUsername(ctx *gin.Context, user string) (*model.User, error) {
	userSaved, err := r.repo.GetByUser(user)
	if err != nil {
		return nil, err
	}
	if userSaved == nil {
		return nil, errors.New("invalid user")
	}

	return userSaved, nil
}
