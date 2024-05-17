package model

import "time"

type List struct {
	Id        string
	Name      string
	Desc      string
	GroupId   string
	Items     []ListItem
	Type      ListType
	CreatedAt time.Time
	CreatedBy string
	Deleted   bool
}
