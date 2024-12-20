package model

import "time"

type Notification struct {
	Id        string
	User      string
	Data      string
	Message   string
	UserOwner string
	Group     string
	ListName  string
	Deleted   bool
	Accepted  bool
	Read      bool
	CreatedAt time.Time
	Type      NotificationType
}

func NewNotification(user, data, message, userOwner, group string, notfType NotificationType) Notification {
	return Notification{
		User:      user,
		Data:      data,
		UserOwner: userOwner,
		Group:     group,
		Message:   message,
		Deleted:   false,
		Accepted:  false,
		Read:      false,
		CreatedAt: time.Now(),
		Type:      notfType,
	}
}

func NewNotificationFull(user, data, message, userOwner, group, list string, notfType NotificationType) Notification {
	return Notification{
		User:      user,
		Data:      data,
		Message:   message,
		UserOwner: userOwner,
		Group:     group,
		ListName:  list,
		Deleted:   false,
		Accepted:  false,
		Read:      false,
		CreatedAt: time.Now(),
		Type:      notfType,
	}
}
