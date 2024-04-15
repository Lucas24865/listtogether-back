package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/utils/token"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	mailBool, userBool, err := r.repo.Exits(&user, ctx)
	if err != nil {
		return err
	}
	if mailBool {
		return errors.New("mail is already registered")
	}
	if userBool {
		return errors.New("user is already registered")
	}

	err = r.repo.Create(&user, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *authService) Login(ctx *gin.Context, user model.User) (string, error) {
	userSaved, err := r.repo.GetByUser(strings.TrimSpace(strings.ToLower(user.User)), ctx)
	if err != nil {
		return "", err
	}
	if userSaved == nil {
		userSaved, err = r.repo.GetByMail(strings.TrimSpace(strings.ToLower(user.User)), ctx)
		if userSaved == nil {
			return "", errors.New("invalid user or email")
		}
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

	user.User = strings.TrimSpace(strings.ToLower(user.User))

	return user, nil
}
