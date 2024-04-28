package repository

import (
	"ListTogetherAPI/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type NotificationRepository interface {
	GetAll(user string, ctx *gin.Context) ([]*model.Notification, error)
	Get(id string, ctx *gin.Context) (*model.Notification, error)
	Remove(id string, ctx *gin.Context) error
	Accept(id string, ctx *gin.Context) error
	Decline(id string, ctx *gin.Context) error
	Add(not model.Notification, ctx *gin.Context) error
}

type notificationRepository struct {
	repo *Repository
}

func (r *notificationRepository) Get(id string, ctx *gin.Context) (*model.Notification, error) {
	notification, err := r.repo.GetById("notifications", id, ctx)
	if err != nil {
		return nil, err
	}

	return mapNotification(notification), nil
}

func (r *notificationRepository) GetAll(u string, ctx *gin.Context) ([]*model.Notification, error) {
	notificationsRaw, err := r.repo.FindAll("notifications", "User", u, "==", ctx)
	if err != nil {
		return nil, err
	}
	notifications := make([]*model.Notification, 0)

	for _, notf := range notificationsRaw {
		notifications = append(notifications, mapNotification(notf))
	}

	return notifications, nil
}

func (r *notificationRepository) Remove(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Deleted = true

	return r.repo.Update("notifications", id, notification, ctx)
}

func (r *notificationRepository) Accept(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Accepted = true
	notification.Read = true

	return r.repo.Update("notifications", id, notification, ctx)
}

func (r *notificationRepository) Decline(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Accepted = false
	notification.Read = true

	return r.repo.Update("notifications", id, notification, ctx)
}

func (r *notificationRepository) Add(not model.Notification, ctx *gin.Context) error {
	not.Id = uuid.New().String()
	return r.repo.Create("notifications", not.Id, not, ctx)
}

func NewNotificationRepository(repo *Repository) NotificationRepository {
	return &notificationRepository{
		repo: repo,
	}
}

func mapNotification(u map[string]interface{}) *model.Notification {
	if u == nil {
		return nil
	}
	notification := model.Notification{
		User:      u["User"].(string),
		Body:      u["Body"].(string),
		Deleted:   u["Deleted"].(bool),
		Accepted:  u["Deleted"].(bool),
		Read:      u["Read"].(bool),
		CreatedAt: u["CreatedAt"].(time.Time)}
	return &notification
}
