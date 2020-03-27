package listusers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/homework/hw3/handlers"
	"github.com/homework/hw3/infra"
	user "github.com/homework/hw3/models"
)

// Handle comment
func Handle(db infra.UserDatabase) handlers.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {

		var users []*user.User = make([]*user.User, 0, len(db))
		for _, val := range db {
			if val.Enabled && !val.Expired {
				users = append(users, val)
			}
		}

		sortBy := r.URL.Query().Get("sortBy")

		switch sortBy {
		case "":
			sort.Slice(users, func(i, j int) bool {
				return users[i].ID < users[j].ID
			})
		case "first-name":
			sort.Slice(users, func(i, j int) bool {
				return strings.ToLower(users[i].FirstName) < strings.ToLower(users[j].FirstName)
			})
		case "last-name":
			sort.Slice(users, func(i, j int) bool {
				return strings.ToLower(users[i].LastName) < strings.ToLower(users[j].LastName)
			})
		case "username":
			sort.Slice(users, func(i, j int) bool {
				return strings.ToLower(users[i].Username) < strings.ToLower(users[j].Username)
			})
		case "email":
			sort.Slice(users, func(i, j int) bool {
				return strings.ToLower(users[i].Email) < strings.ToLower(users[j].Email)
			})
		case "role":
			sort.Slice(users, func(i, j int) bool {
				return users[i].Roles > users[j].Roles
			})
		default:
			msg := "invalid sortBy argument"
			return infra.NewHTTPError(errors.New(msg), 400, msg)
		}

		usersJSON, err2 := json.Marshal(users)
		if err2 != nil {
			return infra.NewHTTPError(err2, 500, "json marshall error")
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write(usersJSON)
		return nil
	}

	return ret
}
