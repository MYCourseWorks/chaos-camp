package deleteuser

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/homework/hw3/handlers"
	"github.com/homework/hw3/infra"
)

// Handle comment
func Handle(db infra.UserDatabase) handlers.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		val, ok := vars["userID"]
		if !ok {
			msg := "userID could not be found"
			return infra.NewHTTPError(errors.New(msg), 400, msg)
		}

		userID, err := strconv.Atoi(val)
		if err != nil {
			return infra.NewHTTPError(err, 400, "userID failed to parse as an integer")
		}

		user := db.FindUserByID(userID)
		if user == nil || !user.Enabled || !user.Expired {
			return infra.NewHTTPError(err, 404, "could not find user")
		}

		db[userID].Enabled = false

		w.WriteHeader(200)
		return nil
	}

	return ret
}
