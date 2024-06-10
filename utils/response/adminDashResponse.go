package response

import "ListTogetherAPI/internal/model"

type AdminDashStatsResponse struct {
	Users         int64
	Groups        int64
	Lists         int64
	Items         int64
	ListsTypes    map[model.ListType]int64
	ElementsTypes map[model.ListType]int64
}

type AdminDashGraphResponse struct {
	UsersCreated  map[string]int64
	Logins        map[string]int64
	GroupsCreated map[string]int64
	ListsCreated  map[string]int64
	ItemsCreated  map[string]int64
}
