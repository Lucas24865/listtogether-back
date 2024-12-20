package model

import (
	"reflect"
	"time"
)

type ListItem struct {
	Name        string
	Quantity    string
	Desc        string
	Completed   bool
	CreatedAt   time.Time
	CreatedBy   string
	LimitDate   time.Time
	CompletedBy string
}

func (item ListItem) Compare(other ListItem) bool {
	val1 := reflect.ValueOf(item)
	val2 := reflect.ValueOf(other)

	for i := 0; i < val1.NumField(); i++ {
		value1 := val1.Field(i).Interface()
		value2 := val2.Field(i).Interface()

		if value1 != value2 {
			return true
		}
	}

	return false
}
