package model

import "time"

type Group struct {
	Picture   string   `json:"picture"`
	Name      string   `json:"name"`
	Desc      string   `json:"desc"`
	Id        string   `json:"id"`
	Admins    []string `json:"admins"`
	CreatedAt time.Time
}
