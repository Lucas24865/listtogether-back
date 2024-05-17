package model

import "time"

type User struct {
	User      string `json:"User" binding:"required"`
	Pass      string `json:"Pass" binding:"required"`
	Mail      string
	Color     string
	Picture   string
	Name      string
	CreatedAt time.Time
}
