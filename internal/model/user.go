package model

type User struct {
	User string `json:"user" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
