package model

import "time"

type ListItem struct {
	Name        string
	Quantity    string
	Desc        string
	Completed   bool
	CreatedAt   time.Time
	CreatedBy   string
	LimitDate   time.Time
	CompletedBy string
}
