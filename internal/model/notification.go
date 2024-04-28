package model

import "time"

const (
	InvitationNotification = "invite"
	MessageNotification    = "message"
	ReminderNotification   = "reminder"
)

type Notification struct {
	Id        string
	User      string
	Data      string
	Info      string
	Deleted   bool
	Accepted  bool
	Read      bool
	CreatedAt time.Time
}

func NewNotification(user, data, info string) Notification {
	return Notification{
		User:      user,
		Data:      data,
		Info:      info,
		Deleted:   false,
		Accepted:  false,
		Read:      false,
		CreatedAt: time.Now(),
	}
}
