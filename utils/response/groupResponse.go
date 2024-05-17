package response

import (
	"ListTogetherAPI/internal/model"
	"time"
)

type GroupResponse struct {
	Id        string
	Name      string
	Desc      string
	CreatedBy model.User
	CreatedAt time.Time
	Members   []model.User
	Admins    []model.User
	Lists     int
}
