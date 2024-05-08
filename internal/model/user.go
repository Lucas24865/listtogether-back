package model

import "time"

type User struct {
	User      string   `json:"user" binding:"required"`
	Pass      string   `json:"pass" binding:"required"`
	Mail      string   `json:"mail"`
	Color     string   `json:"color"`
	Picture   string   `json:"picture"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	CreatedAt time.Time
}
