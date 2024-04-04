package repository

import (
	"ListTogetherAPI/internal/model"
	"context"
)

type UserRepository interface {
	Create(b *model.User) error
	Delete(b *model.User) error
	GetByUser(u string) (*model.User, error)
	Update(b string, m map[string]interface{}) error
}

type userRepository struct {
	repo *Repository
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepository{
		repo: repo,
	}
}

func (r *userRepository) Create(u *model.User) error {
	if err := r.repo.NewRef("users/"+u.User).Set(context.Background(), u); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(u *model.User) error {
	return r.repo.NewRef("users/" + u.User).Delete(context.Background())
}

func (r *userRepository) GetByUser(u string) (*model.User, error) {
	user := &model.User{}
	if err := r.repo.NewRef("users/"+u).Get(context.Background(), user); err != nil {
		return nil, err
	}
	if user.User == "" {
		return nil, nil
	}

	return user, nil
}

func (r *userRepository) Update(b string, m map[string]interface{}) error {
	return r.repo.NewRef("users/"+b).Update(context.Background(), m)
}
