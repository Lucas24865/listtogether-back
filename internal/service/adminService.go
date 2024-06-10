package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/response"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type AdminService interface {
	GetAll(user string, ctx *gin.Context) ([]model.User, error)
	GetDashStats(user string, ctx *gin.Context) (*response.AdminDashStatsResponse, error)
	GetDashGraphs(user string, from, to time.Time, ctx *gin.Context) (*response.AdminDashGraphResponse, error)
}

type adminService struct {
	userRepo     repository.UserRepository
	adminRepo    repository.AdminRepository
	groupService GroupService
}

func NewAdminService(userRepo repository.UserRepository, adminRepo repository.AdminRepository, groupService GroupService) AdminService {
	return &adminService{
		userRepo:     userRepo,
		adminRepo:    adminRepo,
		groupService: groupService,
	}
}

func (s *adminService) GetAll(user string, ctx *gin.Context) ([]model.User, error) {
	userSaved, err := s.userRepo.GetByUser(strings.TrimSpace(strings.ToLower(user)), ctx)
	if err != nil {
		return nil, err
	}
	if !userSaved.Admin {
		return nil, errors.New("not admin")
	}

	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *adminService) GetDashStats(user string, ctx *gin.Context) (*response.AdminDashStatsResponse, error) {
	userSaved, err := s.userRepo.GetByUser(strings.TrimSpace(strings.ToLower(user)), ctx)
	if err != nil {
		return nil, err
	}
	if !userSaved.Admin {
		return nil, errors.New("not admin")
	}

	response, err := s.adminRepo.GetDashStats(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *adminService) GetDashGraphs(user string, from, to time.Time, ctx *gin.Context) (*response.AdminDashGraphResponse, error) {
	userSaved, err := s.userRepo.GetByUser(strings.TrimSpace(strings.ToLower(user)), ctx)
	if err != nil {
		return nil, err
	}
	if !userSaved.Admin {
		return nil, errors.New("not admin")
	}

	response, err := s.adminRepo.GetDashGraphs(from, to, ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}
