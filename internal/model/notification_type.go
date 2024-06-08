package model

const (
	GroupInvite = "new-group"
	GenericType = "generic"
)

type NotificationType string

func (n NotificationType) GetNotifications() []string {
	return []string{GroupInvite, GenericType}
}

func (n NotificationType) ToString() string {
	return string(n)
}
