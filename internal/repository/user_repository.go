package repository

import (
	"ListTogetherAPI/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

type UserRepository interface {
	Create(b *model.User, ctx *gin.Context) error
	Exits(u *model.User, ctx *gin.Context) (bool, bool, error)
	Delete(b *model.User, ctx *gin.Context) error
	GetByUser(u string, ctx *gin.Context) (*model.User, error)
	GetByMail(u string, ctx *gin.Context) (*model.User, error)
	//Update(b string, m map[string]interface{}) error
}

type userRepository struct {
	repo *Repository
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepository{
		repo: repo,
	}
}

func (r *userRepository) Create(u *model.User, ctx *gin.Context) error {
	u.CreatedAt = time.Now()
	if err := r.repo.Create("users", u.User, u, ctx); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(u *model.User, ctx *gin.Context) error {
	err := r.repo.Delete("users", u.User, ctx)
	return err
}

func (r *userRepository) GetByUser(u string, ctx *gin.Context) (*model.User, error) {
	user, err := r.repo.GetById("users", u, ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(user), nil
}

/*func (r *userRepository) Update(b string, m map[string]interface{}) error {
	return r.repo.NewRef("users/"+b).Update(context.Background(), m)
}*/

func (r *userRepository) GetByMail(u string, ctx *gin.Context) (*model.User, error) {
	user, err := r.repo.FindFirst("users", "Mail", u, "==", ctx)
	if err != nil {
		return nil, err
	}
	return mapUser(user), nil
}

// mail,user,error
func (r *userRepository) Exits(u *model.User, ctx *gin.Context) (bool, bool, error) {
	user, err := r.GetByMail(u.Mail, ctx)
	if err != nil {
		return false, false, err
	}
	if user != nil {
		return true, false, nil
	}

	user, err = r.GetByUser(u.User, ctx)
	if err != nil {
		return false, false, err
	}
	if user != nil {
		return false, true, nil
	}

	return false, false, nil
}

func mapUser(u map[string]interface{}) *model.User {
	if u == nil {
		return nil
	}
	user := model.User{
		User:    u["User"].(string),
		Pass:    u["Pass"].(string),
		Mail:    u["Mail"].(string),
		Color:   u["Color"].(string),
		Picture: u["Picture"].(string),
		//Groups:    u["Groups"].([]interface{}),
		Name:      u["Name"].(string),
		CreatedAt: u["CreatedAt"].(time.Time)}
	return &user
}
