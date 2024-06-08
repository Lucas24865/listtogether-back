package repository

import (
	"ListTogetherAPI/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

type UserRepository interface {
	Create(b *model.User, ctx *gin.Context) error
	Exits(mail, user string, ctx *gin.Context) (bool, bool, error)
	Delete(b *model.User, ctx *gin.Context) error
	GetByUser(u string, ctx *gin.Context) (*model.User, error)
	GetByMail(u string, ctx *gin.Context) (*model.User, error)
	GetByUserFull(u string, ctx *gin.Context) (*model.User, error)
	GetByMailFull(u string, ctx *gin.Context) (*model.User, error)
	Update(user *model.User, ctx *gin.Context) error
	GetAll(ctx *gin.Context) ([]model.User, error)
}

type userRepository struct {
	repo *Repository
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepository{
		repo: repo,
	}
}

func (u *userRepository) Create(user *model.User, ctx *gin.Context) error {
	user.CreatedAt = time.Now()

	return u.repo.Create("users", user.User, user, ctx)
}

func (u *userRepository) Delete(user *model.User, ctx *gin.Context) error {
	return u.repo.Delete("users", user.User, ctx)
}

func (u *userRepository) Update(user *model.User, ctx *gin.Context) error {
	return u.repo.Update("users", user.User, *user, ctx)
}

func (u *userRepository) GetByUser(userId string, ctx *gin.Context) (*model.User, error) {
	user, err := u.repo.GetById("users", userId, ctx)
	if err != nil {
		return nil, err
	}

	return mapUser(user), nil
}

func (u *userRepository) GetAll(ctx *gin.Context) ([]model.User, error) {
	usersSaved, err := u.repo.FindAll("users", "User", "", "!=", ctx)
	if err != nil {
		return nil, err
	}

	users := []model.User{}

	for _, userSaved := range usersSaved {
		users = append(users, *mapUser(userSaved))
	}

	return users, nil
}

func (u *userRepository) GetByMail(userId string, ctx *gin.Context) (*model.User, error) {
	user, err := u.repo.FindFirst("users", "Mail", userId, "==", ctx)
	if err != nil {
		return nil, err
	}

	return mapUser(user), nil
}

func (u *userRepository) GetByUserFull(userId string, ctx *gin.Context) (*model.User, error) {
	user, err := u.repo.GetById("users", userId, ctx)
	if err != nil {
		return nil, err
	}

	return mapUserComplete(user), nil
}

func (u *userRepository) GetByMailFull(userId string, ctx *gin.Context) (*model.User, error) {
	user, err := u.repo.FindFirst("users", "Mail", userId, "==", ctx)
	if err != nil {
		return nil, err
	}

	return mapUserComplete(user), nil
}

/*func (r *userRepository) Update(b string, m map[string]interface{}) error {
	return r.repo.NewRef("users/"+b).Update(context.Background(), m)
}*/

func (u *userRepository) Exits(mail, user string, ctx *gin.Context) (bool, bool, error) {
	userSaved, err := u.GetByMail(mail, ctx)
	if err != nil {
		return false, false, err
	}
	if userSaved != nil {
		return true, false, nil
	}

	userSaved, err = u.GetByUser(user, ctx)
	if err != nil {
		return false, false, err
	}
	if userSaved != nil {
		return false, true, nil
	}

	return false, false, nil
}

func mapUserComplete(u map[string]interface{}) *model.User {
	if len(u) == 0 {
		return nil
	}

	user := model.User{
		User:      u["User"].(string),
		Pass:      u["Pass"].(string),
		Mail:      u["Mail"].(string),
		Color:     u["Color"].(string),
		Picture:   u["Picture"].(string),
		Name:      u["Name"].(string),
		Admin:     u["Admin"].(bool),
		LastLogin: u["LastLogin"].(time.Time),
		CreatedAt: u["CreatedAt"].(time.Time)}

	return &user
}

func mapUser(u map[string]interface{}) *model.User {
	if len(u) == 0 {
		return nil
	}

	user := mapUserComplete(u)
	user.Pass = ""

	return user
}
