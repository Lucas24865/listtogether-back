package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthService interface {
	Register(user model.User, ctx *gin.Context) error
	Login(user model.User, ctx *gin.Context) (string, error)
}

type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{
		userService: userService,
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
