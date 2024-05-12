package model

import "time"

type ListItem struct {
	Id                string
	ListId            string
	Name              string
	Quantity          int
	Desc              string
	Completed         bool
	GroupId           string
	CreatedAt         time.Time
	CreatedBy         string
	LimitDate         time.Time
	CompletedByUserId string
}
