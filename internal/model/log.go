package model

import "time"

type Log struct {
	User      string
	Id        string
	Type      LogType
	CreatedAt time.Time
}
