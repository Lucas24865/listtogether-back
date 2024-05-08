package model

type ListType string

const (
	NotesListType     = "notes"
	RemindersListType = "reminders"
	AccountsListType  = "accounts"
	ShoppingListType  = "shopping"
	EventsListType    = "events"
)

func (l ListType) GetTypes() []string {
	return []string{NotesListType, RemindersListType, AccountsListType, ShoppingListType, EventsListType}
}
