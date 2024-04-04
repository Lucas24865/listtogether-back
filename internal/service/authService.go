package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"strings"
)

type AuthService interface {
	Register(ctx *gin.Context, user model.User) error
	Login(ctx *gin.Context, user model.User) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (r *authService) Register(ctx *gin.Context, user model.User) error {
	user, err := beforeSave(user)
	if err != nil {
		return err
	}

	err = r.repo.Create(&user)
	if err != nil {
		return err
	}
	return nil
}

func (r *authService) Login(ctx *gin.Context, user model.User) (string, error) {
	userSaved, err := r.repo.GetByUser(user.User)
	if err != nil {
		return "", err
	}
	if userSaved == nil {
		return "", errors.New("invalid user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userSaved.Pass), []byte(user.Pass))
	if err != nil {
		return "", errors.New("invalid pass")
	}

	generatedToken, err := token.GenerateToken(user.User)
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

	user.User = html.EscapeString(strings.TrimSpace(user.User))

	return user, nil
}
