package model

import "time"

type List struct {
	Id        string
	Name      string
	Desc      string
	GroupId   string
	CreatedAt time.Time
	CreatedBy string
}
