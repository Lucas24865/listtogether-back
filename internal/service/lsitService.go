package service

import (
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/response"
	"github.com/gin-gonic/gin"
)

type ListService interface {
	GetAll(userId string, ctx *gin.Context) (*response.GroupResponse, error)
}

type listService struct {
	repo repository.ListRepository
}

func NewListService(repo repository.GroupRepository, notificationService NotificationService) ListService {
	return &listService{
		repo: repo,
	}
}

func (s *listService) Get(ids []string, ctx *gin.Context) ([]response.GroupResponse, error) {
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}

	groupInfo, err := s.repo.GetGroups(ids, ctx)
	if err != nil {
		return nil, err
	}

	return groupInfo, nil
}
