package model

import "time"

type Group struct {
	Picture     string
	Name        string
	Desc        string
	Id          string
	Deactivated bool
	Admins      []string
	Users       []string
	CreatedAt   time.Time
	CreatedBy   string
}
