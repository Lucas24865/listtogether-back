package model

import "time"

type List struct {
	Picture   string
	Name      string
	Desc      string
	Id        string
	Admins    []string
	CreatedAt time.Time
	CreatedBy string
}
