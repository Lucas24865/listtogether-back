package response

import (
	"ListTogetherAPI/internal/model"
	"time"
)

type ListItemResponse struct {
	Id          string
	Name        string
	Quantity    string
	Desc        string
	Completed   bool
	CreatedAt   time.Time
	CreatedBy   model.User
	LimitDate   time.Time
	CompletedBy model.User
}
