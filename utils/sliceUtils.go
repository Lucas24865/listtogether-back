package utils

import "ListTogetherAPI/internal/model"

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsUserId(s []model.User, e string) bool {
	for _, a := range s {
		if a.User == e {
			return true
		}
	}
	return false
}
