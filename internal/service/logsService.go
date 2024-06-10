package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type LogsService interface {
	Add(log model.Log, ctx *gin.Context) error
	AddLogin(user string, ctx *gin.Context) error
}

type logService struct {
	repo *repository.Repository
}

func NewLogService(repo *repository.Repository) LogsService {
	return &logService{
		repo: repo,
	}
}

func (s *logService) Add(log model.Log, ctx *gin.Context) error {
	id := uuid.New().String()
	log.Id = id
	return s.repo.Create("logs", id, log, ctx)
}

func (s *logService) AddLogin(user string, ctx *gin.Context) error {
	id := uuid.New().String()
	log := model.Log{
		User:      user,
		Id:        id,
		Type:      model.LoginType,
		CreatedAt: time.Now(),
	}
	return s.repo.Create("logs", id, log, ctx)
}
