package service

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
)

type NotificationService interface {
	GetAll(user string, ctx *gin.Context) ([]*model.Notification, error)
	Get(id, user string, ctx *gin.Context) (*model.Notification, error)
	Remove(id, user string, ctx *gin.Context) error
	RemoveAll(user string, ctx *gin.Context) error
	Accept(id, user string, ctx *gin.Context) (*model.Notification, error)
	Decline(id, user string, ctx *gin.Context) error
	SendNew(user, group, message string, notType model.NotificationType, ctx *gin.Context) error
	SendNewMultiple(users []string, group, message string, notType model.NotificationType, ctx *gin.Context) error
	Add(not model.Notification, ctx *gin.Context) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func (n notificationService) GetAll(user string, ctx *gin.Context) ([]*model.Notification, error) {
	return n.repo.GetAll(user, ctx)
}

func (n notificationService) Get(id, user string, ctx *gin.Context) (*model.Notification, error) {
	notif, err := n.get(id, user, ctx)
	if err != nil {
		return nil, err
	}

	return notif, nil
}

func (n notificationService) Remove(id, user string, ctx *gin.Context) error {
	notif, err := n.get(id, user, ctx)
	if err != nil {
		return err
	}
	if notif.Deleted {
		return nil
	}

	return n.repo.Remove(id, ctx)
}

func (n notificationService) RemoveAll(user string, ctx *gin.Context) error {
	return n.repo.RemoveAll(user, ctx)
}

func (n notificationService) Accept(id, user string, ctx *gin.Context) (*model.Notification, error) {
	notif, err := n.get(id, user, ctx)
	if err != nil {
		return nil, err
	}
	if notif.Accepted {
		return nil, errors.New("notification already accepted")
	}
	if notif.Read {
		return nil, errors.New("notification already declined")
	}

	return notif, n.repo.Accept(id, ctx)
}

func (n notificationService) Decline(id, user string, ctx *gin.Context) error {
	notif, err := n.get(id, user, ctx)
	if err != nil {
		return err
	}
	if notif.Accepted {
		return errors.New("notification already accepted")
	}
	if notif.Read {
		return errors.New("notification already declined")
	}

	return n.repo.Decline(id, ctx)
}

func (n notificationService) SendNew(user, group, message string, notType model.NotificationType, ctx *gin.Context) error {
	notif := model.NewNotification(user, group, message, notType)

	return n.repo.Add(notif, ctx)
}

func (n notificationService) SendNewMultiple(users []string, group, message string, notType model.NotificationType, ctx *gin.Context) error {
	for _, user := range users {
		err := n.SendNew(user, group, message, notType, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n notificationService) Add(notif model.Notification, ctx *gin.Context) error {
	return n.repo.Add(notif, ctx)
}

func (n notificationService) get(notifId, user string, ctx *gin.Context) (*model.Notification, error) {
	notification, err := n.repo.Get(notifId, ctx)
	if err != nil {
		return nil, err
	}
	if notification.User != user {
		return nil, errors.New("invalid notification")
	}

	return notification, nil
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}
