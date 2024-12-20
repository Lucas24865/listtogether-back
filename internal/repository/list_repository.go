package repository

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils"
	"ListTogetherAPI/utils/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
)

type ListRepository interface {
	GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error)
	Get(userId, listId string, ctx *gin.Context) (*response.ListResponse, error)
	Create(request model.List, ctx *gin.Context) error
	Update(request model.List, user string, ctx *gin.Context) error
	Delete(userId, listId string, ctx *gin.Context) error
}

type listRepository struct {
	repo             *Repository
	groupRepo        GroupRepository
	notificationRepo NotificationRepository
}

func (l listRepository) Create(request model.List, ctx *gin.Context) error {
	group, err := l.groupRepo.GetGroupSimple(request.GroupId, ctx)
	if err != nil {
		return err
	}
	if !utils.Contains(group.Users, request.CreatedBy) {
		return errors.New("invalid group")
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.Deleted = false

	err = l.repo.Create("lists", request.Id, request, ctx)
	if err != nil {
		return err
	}

	notMessage := fmt.Sprintf("%s ha agregado la lista: %s, en el grupo: %s", request.CreatedBy,
		request.Name, group.Name)

	err = l.notificationRepo.AddMultipleGeneric(notMessage, "", request.CreatedBy, group.Name, request.Name, group.Users, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (l listRepository) Update(request model.List, userId string, ctx *gin.Context) error {
	list, err := l.getList(request.Id, ctx)
	if err != nil {
		return err
	}

	if list.CreatedBy != request.CreatedBy || list.CreatedAt != request.CreatedAt ||
		list.GroupId != request.GroupId || list.Type != request.Type {
		return errors.New("invalid list")
	}

	for _, item := range list.Items {
		if item.CreatedBy != item.CreatedBy || item.CreatedAt != item.CreatedAt {
			return errors.New("invalid list")
		}
	}

	group, err := l.groupRepo.GetGroupSimple(list.GroupId, ctx)
	if err != nil {
		return err
	}

	changes := "Se modificaron los elementos: "
	modifiedItems := []string{}

	for i, item := range request.Items {
		if list.Items[i].Compare(item) {
			modifiedItems = append(modifiedItems, item.Name)
		}
	}

	changes += strings.Join(modifiedItems, ", ")

	if len(modifiedItems) == 0 {
		changes = "No se modificaron elementos."
	}

	if !utils.Contains(group.Users, userId) {
		return errors.New("invalid group")
	}

	err = l.repo.Update("lists", request.Id, request, ctx)
	if err != nil {
		return err
	}

	notMessage := fmt.Sprintf("%s ha actualizado la lista: %s, en el grupo: %s", userId,
		request.Name, group.Name)

	if len(modifiedItems) == 0 {
		changes = ""
	}

	err = l.notificationRepo.AddMultipleGeneric(notMessage, changes, userId, group.Name, request.Name, group.Users, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (l listRepository) Delete(userId, listId string, ctx *gin.Context) error {
	list, err := l.getList(listId, ctx)
	if err != nil {
		return err
	}

	group, err := l.groupRepo.GetGroupSimple(list.GroupId, ctx)
	if err != nil {
		return err
	}
	if !utils.Contains(group.Users, userId) {
		return errors.New("invalid group")
	}

	list.Deleted = true

	err = l.repo.Update("lists", listId, *list, ctx)
	if err != nil {
		return err
	}

	notMessage := fmt.Sprintf("%s ha agredado la lista: %s, en el grupo: %s", userId,
		list.Name, group.Name)

	err = l.notificationRepo.AddMultipleGeneric(notMessage, "", userId, group.Name, list.Name, group.Users, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (l listRepository) getList(listId string, ctx *gin.Context) (*model.List, error) {
	list, err := l.repo.GetById("lists", listId, ctx)
	if err != nil {
		return nil, err
	}

	return mapList(list), nil
}

func NewListRepository(repo *Repository, groupRepo GroupRepository, notificationRepo NotificationRepository) ListRepository {
	return &listRepository{
		repo:             repo,
		groupRepo:        groupRepo,
		notificationRepo: notificationRepo,
	}
}

func (l listRepository) GetAll(userId string, ctx *gin.Context) ([]response.ListResponse, error) {

	groups, err := l.groupRepo.GetGroupsFull(userId, ctx)
	if err != nil {
		return nil, err
	}

	var groupsId []string
	mapGroups := make(map[string]response.GroupResponse)
	mapUsers := make(map[string]model.User)

	for _, group := range groups {
		groupsId = append(groupsId, group.Id)

		for _, member := range group.Members {
			mapUsers[member.User] = member
		}

		for _, member := range group.Admins {
			mapUsers[member.User] = member
		}

		mapGroups[group.Id] = group
	}

	lists, err := l.repo.FindAll("lists", "GroupId", groupsId, "in", ctx)
	if err != nil {
		return nil, err
	}

	var listResponses []response.ListResponse
	for _, list := range lists {
		if list["Deleted"] != nil && list["Deleted"].(bool) {
			continue
		}

		listMapped := *mapList(list)

		listResponse := response.ListResponse{
			Id:        listMapped.Id,
			Name:      listMapped.Name,
			Desc:      listMapped.Desc,
			Group:     mapGroups[listMapped.GroupId],
			CreatedBy: mapUsers[listMapped.CreatedBy],
			CreatedAt: listMapped.CreatedAt,
			Items:     []response.ListItemResponse{},
			Type:      listMapped.Type,
		}

		for _, item := range listMapped.Items {
			itemResponse := response.ListItemResponse{
				Name:        item.Name,
				Quantity:    item.Quantity,
				Desc:        item.Desc,
				Completed:   item.Completed,
				CreatedAt:   item.CreatedAt,
				CreatedBy:   mapUsers[item.CreatedBy],
				LimitDate:   item.LimitDate,
				CompletedBy: mapUsers[item.CompletedBy],
			}
			listResponse.Items = append(listResponse.Items, itemResponse)
		}

		listResponses = append(listResponses, listResponse)
	}

	return listResponses, nil
}

func (l listRepository) Get(userId string, listId string, ctx *gin.Context) (*response.ListResponse, error) {

	list, err := l.repo.FindFirst("lists", "Id", listId, "==", ctx)
	if err != nil {
		return nil, err
	}

	listMapped := *mapList(list)

	group, err := l.groupRepo.GetGroupFull(listMapped.GroupId, ctx)
	if err != nil {
		return nil, err
	}

	mapUsers := make(map[string]model.User)
	for _, member := range group.Members {
		mapUsers[member.User] = member
	}
	for _, member := range group.Admins {
		mapUsers[member.User] = member
	}

	if _, exists := mapUsers[userId]; !exists {
		return nil, errors.New("invalid list")
	}

	listResponse := response.ListResponse{
		Id:        listMapped.Id,
		Name:      listMapped.Name,
		Desc:      listMapped.Desc,
		Group:     *group,
		CreatedBy: mapUsers[listMapped.CreatedBy],
		CreatedAt: listMapped.CreatedAt,
		Items:     []response.ListItemResponse{},
		Type:      listMapped.Type,
	}

	for _, item := range listMapped.Items {
		itemResponse := response.ListItemResponse{
			Name:        item.Name,
			Quantity:    item.Quantity,
			Desc:        item.Desc,
			Completed:   item.Completed,
			CreatedAt:   item.CreatedAt,
			CreatedBy:   mapUsers[item.CreatedBy],
			LimitDate:   item.LimitDate,
			CompletedBy: mapUsers[item.CompletedBy],
		}
		listResponse.Items = append(listResponse.Items, itemResponse)
	}

	return &listResponse, nil
}

func mapList(u map[string]interface{}) *model.List {
	if len(u) == 0 {
		return nil
	}
	var items []model.ListItem

	for _, itemRaw := range u["Items"].([]interface{}) {
		i := itemRaw.(map[string]interface{})
		listItem := model.ListItem{
			Name:        i["Name"].(string),
			Quantity:    i["Quantity"].(string),
			Desc:        i["Desc"].(string),
			Completed:   i["Completed"].(bool),
			CreatedAt:   i["CreatedAt"].(time.Time),
			CreatedBy:   i["CreatedBy"].(string),
			LimitDate:   i["LimitDate"].(time.Time),
			CompletedBy: i["CompletedBy"].(string),
		}
		items = append(items, listItem)
	}

	list := model.List{
		Id:        u["Id"].(string),
		Name:      u["Name"].(string),
		Desc:      u["Desc"].(string),
		GroupId:   u["GroupId"].(string),
		Items:     items,
		Type:      model.ListType(u["Type"].(string)),
		CreatedAt: u["CreatedAt"].(time.Time),
		CreatedBy: u["CreatedBy"].(string),
	}

	return &list
}
