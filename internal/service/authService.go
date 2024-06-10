package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthService interface {
	Register(user model.User, ctx *gin.Context) error
	Login(user model.User, ctx *gin.Context) (string, error)
	AdminLogin(user model.User, ctx *gin.Context) (string, error)
}

type authService struct {
	userService UserService
	logService  LogsService
}

func NewAuthService(userService UserService, logService LogsService) AuthService {
	return &authService{
		userService: userService,
		logService:  logService,
	}
}

func (r *authService) Register(user model.User, ctx *gin.Context) error {
	userProcessed, err := beforeSave(user)
	if err != nil {
		return err
	}

	err = r.userService.Register(userProcessed, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *authService) Login(user model.User, ctx *gin.Context) (string, error) {
	user.User = strings.TrimSpace(strings.ToLower(user.User))

	userSaved, err := r.userService.GetByUsernameOrEmailLogin(user.User, ctx)
	if err != nil {
		return "", err
	}
	if userSaved == nil {
		return "", errors.New("invalid user or email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userSaved.Pass), []byte(user.Pass))
	if err != nil {
		return "", errors.New("invalid pass")
	}

	userSaved.LastLogin = time.Now()
	err = r.userService.Update(*userSaved, ctx)
	if err != nil {
		return "", err
	}

	err = r.logService.AddLogin(userSaved.User, ctx)
	if err != nil {
		return "", err
	}

	generatedToken, err := token.GenerateToken(userSaved.User)
	if err != nil {
		return "", err
	}

	return generatedToken, nil
}

func (r *authService) AdminLogin(user model.User, ctx *gin.Context) (string, error) {
	user.User = strings.TrimSpace(strings.ToLower(user.User))

	userSaved, err := r.userService.GetByUsernameOrEmailLogin(user.User, ctx)
	if err != nil {
		return "", err
	}
	if userSaved == nil {
		return "", errors.New("invalid user or email")
	}
	if !userSaved.Admin {
		return "", errors.New("user not admin")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userSaved.Pass), []byte(user.Pass))
	if err != nil {
		return "", errors.New("invalid pass")
	}

	userSaved.LastLogin = time.Now()
	err = r.userService.Update(*userSaved, ctx)
	if err != nil {
		return "", err
	}

	generatedToken, err := token.GenerateToken(userSaved.User)
	if err != nil {
		return "", err
	}

	return generatedToken, nil
}

func beforeSave(user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Pass = string(hashedPassword)
	user.User = strings.TrimSpace(strings.ToLower(user.User))

	return user, nil
}
