package response

import (
	"ListTogetherAPI/internal/model"
	"time"
)

type ListResponse struct {
	Id        string
	Name      string
	Desc      string
	Group     GroupResponse
	CreatedBy model.User
	CreatedAt time.Time
	Items     []ListItemResponse
	Type      model.ListType
}
