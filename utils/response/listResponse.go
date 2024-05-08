package response

import (
	"ListTogetherAPI/internal/model"
	"time"
)

type ListResponse struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Desc      string       `json:"desc"`
	CreatedBy model.User   `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	Members   []model.User `json:"members"`
	Admins    []model.User `json:"admins"`
	Lists     int          `json:"lists"`
}
