package repository

import (
	"ListTogetherAPI/internal/model"
	"ListTogetherAPI/utils/response"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminRepository interface {
	GetDashStats(ctx *gin.Context) (*response.AdminDashStatsResponse, error)
	GetDashGraphs(from, to time.Time, ctx *gin.Context) (*response.AdminDashGraphResponse, error)
}

type adminRepository struct {
	repo *Repository
}

func NewAdminRepository(repo *Repository) AdminRepository {
	return &adminRepository{
		repo: repo,
	}
}

func (a *adminRepository) GetDashStats(ctx *gin.Context) (*response.AdminDashStatsResponse, error) {
	usersSaved, err := a.repo.Count("users", nil, nil, ctx)
	if err != nil {
		return nil, err
	}

	listsSaved, err := a.repo.GetAll("lists", ctx)
	if err != nil {
		return nil, err
	}

	groupsSaved, err := a.repo.Count("groups", nil, nil, ctx)
	if err != nil {
		return nil, err
	}

	var listsSavedMapped []model.List
	for _, list := range listsSaved {
		listsSavedMapped = append(listsSavedMapped, *mapList(list))
	}

	itemsSaved := int64(0)
	listsTypes := make(map[model.ListType]int64)
	elementsTypes := make(map[model.ListType]int64)

	for _, list := range listsSavedMapped {
		itemsSaved += int64(len(list.Items))
		listsTypes[list.Type]++
		elementsTypes[list.Type] += int64(len(list.Items))
	}

	return &response.AdminDashStatsResponse{
		Users:         usersSaved,
		Groups:        groupsSaved,
		Lists:         int64(len(listsSaved)),
		Items:         itemsSaved,
		ListsTypes:    listsTypes,
		ElementsTypes: elementsTypes,
	}, nil
}

func (a *adminRepository) GetDashGraphs(from, to time.Time, ctx *gin.Context) (*response.AdminDashGraphResponse, error) {
	difference := to.Sub(from)

	if difference.Hours() > 7*24 {
		return nil, errors.New("invalid date range, date range must be less than 7 days")
	}

	usersCreated, err := a.getByDay(from, to, "CreatedAt", "users", ctx)
	if err != nil {
		return nil, err
	}

	logins, err := a.getByDay(from, to, "CreatedAt", "logs", ctx)
	if err != nil {
		return nil, err
	}

	groupsCreated, err := a.getByDay(from, to, "CreatedAt", "groups", ctx)
	if err != nil {
		return nil, err
	}

	listsCreated, err := a.getByDay(from, to, "CreatedAt", "lists", ctx)
	if err != nil {
		return nil, err
	}

	/*itemsCreated,err := a.getByDay(from,to,"CreatedAt","users",ctx)
	if err != nil {
		return nil, err
	}*/

	return &response.AdminDashGraphResponse{
		UsersCreated:  usersCreated,
		Logins:        logins,
		GroupsCreated: groupsCreated,
		ListsCreated:  listsCreated,
		ItemsCreated:  nil,
	}, nil

}

func (a *adminRepository) getByDay(from, to time.Time, prop, collection string, ctx *gin.Context) (map[string]int64, error) {
	results := make(map[string]int64)
	current := from

	for !current.After(to) {
		startOfDay := current
		endOfDay := current.Add(24 * time.Hour)

		from := &Filters{prop, ">=", startOfDay}
		to := &Filters{prop, "<", endOfDay}

		count, err := a.repo.Count(collection, from, to, ctx)
		if err != nil {
			return nil, err
		}

		dateStr := current.Format("02/01/2006")
		results[dateStr] = count

		current = current.Add(24 * time.Hour)
	}

	return results, nil
}
