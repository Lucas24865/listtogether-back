package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/response"
	"github.com/gin-gonic/gin"
)

type ListService interface {
	GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error)
	Get(userId, listId string, ctx *gin.Context) (*response.ListResponse, error)
	Create(request model.List, ctx *gin.Context) error
	Update(request model.List, user string, ctx *gin.Context) error
	Delete(userId, listId string, ctx *gin.Context) error
}

type listService struct {
	repo repository.ListRepository
}

func NewListService(repo repository.ListRepository) ListService {
	return &listService{
		repo: repo,
	}
}

func (s *listService) Create(request model.List, ctx *gin.Context) error {
	return s.repo.Create(request, ctx)
}

func (s *listService) Update(request model.List, user string, ctx *gin.Context) error {
	return s.repo.Update(request, user, ctx)
}

func (s *listService) GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error) {
	return s.repo.GetAll(userId, ctx)
}

func (s *listService) Get(userId, listId string, ctx *gin.Context) (*response.ListResponse, error) {
	return s.repo.Get(userId, listId, ctx)
}

func (s *listService) Delete(userId, listId string, ctx *gin.Context) error {
	return s.repo.Delete(userId, listId, ctx)
}
