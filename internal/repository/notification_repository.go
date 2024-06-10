package repository

import (
	"ListTogetherAPI/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type NotificationRepository struct {
	repo *Repository
}

func (r *NotificationRepository) Get(id string, ctx *gin.Context) (*model.Notification, error) {
	notification, err := r.repo.GetById("notifications", id, ctx)
	if err != nil {
		return nil, err
	}

	return mapNotification(notification), nil
}

func (r *NotificationRepository) GetAll(u string, ctx *gin.Context) ([]*model.Notification, error) {
	notificationsRaw, err := r.repo.FindAllTwoProps("notifications", "User", u, "==", "Deleted", false, "==", ctx)
	if err != nil {
		return nil, err
	}
	notifications := make([]*model.Notification, 0)

	for _, notf := range notificationsRaw {
		notifications = append(notifications, mapNotification(notf))
	}

	return notifications, nil
}

func (r *NotificationRepository) Remove(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Deleted = true

	return r.repo.Update("notifications", id, *notification, ctx)
}

func (r *NotificationRepository) RemoveAll(user string, ctx *gin.Context) error {
	notifRead, err := r.repo.FindAllTwoProps("notifications", "User", user, "==", "Deleted", false, "==", ctx)
	if err != nil {
		return err
	}
	var ids []string
	var deletedProp []map[string]interface{}
	for _, noti := range notifRead {
		if noti["Read"].(bool) && noti["Type"].(string) == model.GroupInvite || noti["Type"].(string) == model.GenericType {
			ids = append(ids, noti["Id"].(string))
			deletedProp = append(deletedProp, map[string]interface{}{"deleted": true})
		}
	}

	return r.repo.UpdateBatch("notifications", ids, deletedProp, ctx)
}

func (r *NotificationRepository) Accept(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Accepted = true
	notification.Read = true

	return r.repo.Update("notifications", id, *notification, ctx)
}

func (r *NotificationRepository) Decline(id string, ctx *gin.Context) error {
	notification, err := r.Get(id, ctx)
	if err != nil {
		return err
	}

	notification.Accepted = false
	notification.Read = true

	return r.repo.Update("notifications", id, *notification, ctx)
}

func (r *NotificationRepository) Add(not model.Notification, ctx *gin.Context) error {
	not.Id = uuid.New().String()
	return r.repo.Create("notifications", not.Id, not, ctx)
}

func (r *NotificationRepository) AddMultipleGeneric(message string, to []string, ctx *gin.Context) error {
	not := model.NewNotification("", "", message, model.GenericType)
	for _, user := range to {
		not.User = user
		err := r.Add(not, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewNotificationRepository(repo *Repository) NotificationRepository {
	return NotificationRepository{
		repo: repo,
	}
}

func mapNotification(u map[string]interface{}) *model.Notification {
	if u == nil {
		return nil
	}
	notification := model.Notification{
		Id:        u["Id"].(string),
		User:      u["User"].(string),
		Data:      u["Data"].(string),
		Deleted:   u["Deleted"].(bool),
		Message:   u["Message"].(string),
		Accepted:  u["Accepted"].(bool),
		Read:      u["Read"].(bool),
		CreatedAt: u["CreatedAt"].(time.Time),
		Type:      model.NotificationType(u["Type"].(string))}
	return &notification
}
