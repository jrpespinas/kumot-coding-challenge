package listing

import (
	"sort"
	"strings"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/domain"
)

func SortUsers(users []domain.User, field string) {
	var less func(i, j int) bool
	switch field {
	case "name":
		less = func(i, j int) bool {
			return strings.ToLower(*users[i].Name) < strings.ToLower(*users[j].Name)
		}
	case "login":
		less = func(i, j int) bool {
			return strings.ToLower(*users[i].Login) < strings.ToLower(*users[j].Login)
		}
	case "follower":
		less = func(i, j int) bool {
			return *users[i].Followers < *users[j].Followers
		}
	default:
		less = func(i, j int) bool {
			return strings.ToLower(*users[i].Name) < strings.ToLower(*users[j].Name)
		}
	}
	sort.Slice(users, less)
}
