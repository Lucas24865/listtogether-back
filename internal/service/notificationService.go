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
	Accept(id, user string, ctx *gin.Context) error
	Decline(id, user string, ctx *gin.Context) error
	Invite(user, group string) error
	Add(not model.Notification) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func (n notificationService) GetAll(user string, ctx *gin.Context) ([]*model.Notification, error) {
	return n.repo.GetAll(user, ctx)
}

func (n notificationService) Get(id, user string, ctx *gin.Context) (*model.Notification, error) {
	err := n.havePermission(id, user, ctx)
	if err != nil {
		return nil, err
	}
	return n.repo.Get(id, ctx)
}

func (n notificationService) Remove(id, user string, ctx *gin.Context) error {
	err := n.havePermission(id, user, ctx)
	if err != nil {
		return err
	}
	return n.repo.Remove(id, ctx)
}

func (n notificationService) Accept(id, user string, ctx *gin.Context) error {
	err := n.havePermission(id, user, ctx)
	if err != nil {
		return err
	}
	return n.repo.Accept(id, ctx)
}

func (n notificationService) Decline(id, user string, ctx *gin.Context) error {
	err := n.havePermission(id, user, ctx)
	if err != nil {
		return err
	}
	return n.repo.Decline(id, ctx)
}

func (n notificationService) Invite(user, group string) error {
	//TODO implement me
	panic("implement me")
}

func (n notificationService) Add(not model.Notification) error {
	//TODO implement me
	panic("implement me")
}

func (n notificationService) havePermission(notifId, user string, ctx *gin.Context) error {
	notification, err := n.repo.Get(notifId, ctx)
	if err != nil {
		return err
	}
	if notification.User != user {
		return errors.New("Invalid notification")
	}
	return nil
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}
