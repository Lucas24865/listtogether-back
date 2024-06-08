package model

import "time"

type Notification struct {
	Id        string
	User      string
	Data      string
	Message   string
	Deleted   bool
	Accepted  bool
	Read      bool
	CreatedAt time.Time
	Type      NotificationType
}

func NewNotification(user, data, message string, notfType NotificationType) Notification {
	return Notification{
		User:      user,
		Data:      data,
		Message:   message,
		Deleted:   false,
		Accepted:  false,
		Read:      false,
		CreatedAt: time.Now(),
		Type:      notfType,
	}
}
