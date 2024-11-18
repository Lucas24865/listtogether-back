package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/requests"
	"ListTogetherAPI/utils/response"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService interface {
	GetByUsernameOrEmail(user string, ctx *gin.Context) (*model.User, error)
	GetByUsernameOrEmailLogin(user string, ctx *gin.Context) (*model.User, error)
	GetAllGroups(user string, ctx *gin.Context) ([]response.GroupResponse, error)
	AcceptInvite(id, user string, ctx *gin.Context) error
	DeclineInvite(id, user string, ctx *gin.Context) error
	Register(user model.User, ctx *gin.Context) error
	Update(user model.User, ctx *gin.Context) error
	Edit(user requests.UserUpdateRequest, ctx *gin.Context) error
	/*	AddGroup(user, group string, ctx *gin.Context) error*/
}

type userService struct {
	repo                repository.UserRepository
	notificationService NotificationService
	groupService        GroupService
}

func NewUserService(repo repository.UserRepository, notificationService NotificationService, groupService GroupService) UserService {
	return &userService{
		repo:                repo,
		notificationService: notificationService,
		groupService:        groupService,
	}
}

func (s *userService) Register(user model.User, ctx *gin.Context) error {
	err := s.exists(user.Mail, user.User, ctx)
	if err != nil {
		return err
	}
	return s.repo.Create(&user, ctx)
}

func (s *userService) Update(user model.User, ctx *gin.Context) error {
	return s.repo.Update(&user, ctx)
}

func (s *userService) Edit(user requests.UserUpdateRequest, ctx *gin.Context) error {
	userSaved, err := s.GetByUsernameOrEmailLogin(user.User, ctx)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userSaved.Pass), []byte(user.OldPass))
	if err != nil {
		return errors.New("invalid pass")
	}

	if user.NewPass != "" {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.NewPass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user.NewPass = string(hashedPassword)
		userSaved.Pass = user.NewPass
	}

	userSaved.Name = user.Name
	userSaved.Color = user.Color

	return s.repo.Update(userSaved, ctx)
}

func (s *userService) GetByUsernameOrEmail(user string, ctx *gin.Context) (*model.User, error) {
	userSaved, err := s.repo.GetByUser(strings.TrimSpace(strings.ToLower(user)), ctx)
	if err != nil {
		return nil, err
	}
	if userSaved == nil {
		userSaved, err = s.repo.GetByMail(strings.TrimSpace(strings.ToLower(user)), ctx)
		if userSaved == nil {
			return nil, err
		}
	}

	return userSaved, nil
}

func (s *userService) GetByUsernameOrEmailLogin(user string, ctx *gin.Context) (*model.User, error) {
	userSaved, err := s.repo.GetByUserFull(strings.TrimSpace(strings.ToLower(user)), ctx)
	if err != nil {
		return nil, err
	}
	if userSaved == nil {
		userSaved, err = s.repo.GetByMailFull(strings.TrimSpace(strings.ToLower(user)), ctx)
		if userSaved == nil {
			return nil, err
		}
	}

	return userSaved, nil
}

/*
	func (s *userService) AddGroup(user, group string, ctx *gin.Context) error {
		return s.repo.AddGroup(group, user, ctx)
	}
*/

func (s *userService) AcceptInvite(id, user string, ctx *gin.Context) error {
	notif, err := s.notificationService.Accept(id, user, ctx)
	if err != nil {
		return err
	}

	return s.groupService.AddMember(user, notif.Data, ctx)
}

func (s *userService) DeclineInvite(id, user string, ctx *gin.Context) error {
	err := s.notificationService.Decline(id, user, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAllGroups(userId string, ctx *gin.Context) ([]response.GroupResponse, error) {
	user, err := s.GetByUsernameOrEmail(userId, ctx)
	if err != nil {
		return nil, err
	}
	groups, err := s.groupService.GetGroupsFull(user.User, ctx)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *userService) exists(mail, user string, ctx *gin.Context) error {
	mailBool, userBool, err := s.repo.Exits(mail, user, ctx)
	if err != nil {
		return err
	}
	if mailBool {
		return errors.New("mail is already registered")
	}
	if userBool {
		return errors.New("user is already registered")
	}

	return nil
}
