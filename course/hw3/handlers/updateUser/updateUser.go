package updateuser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/homework/hw3/handlers"
	"github.com/homework/hw3/infra"
	user "github.com/homework/hw3/models"
)

// CMD comment
type CMD struct {
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	FirstName string      `json:"first-name"`
	LastName  string      `json:"last-name"`
	Password  string      `json:"password"`
	Expired   *bool       `json:"expired"`
	Roles     *user.Roles `json:"roles"`
}

// Handle comment
func Handle(db infra.UserDatabase) handlers.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		// Parse Query Parameters :
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
		if user == nil {
			return infra.NewHTTPError(err, 404, "could not find user")
		}

		// Parse Body :
		var c CMD
		// TODO: json Decoder doesn't check the size of the request body.
		// Which allows for DDOS attacks.
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&c)
		if err != nil {
			return infra.NewHTTPError(err, 400, "request body is not valid")
		}

		update(db, userID, &c)

		var userJSON []byte
		userJSON, err = json.Marshal(user)
		if err != nil {
			return infra.NewHTTPError(err, 500, "json marshall error")
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write(userJSON)
		return nil
	}

	return ret
}

func update(db infra.UserDatabase, id int, updateCMD *CMD) error {
	u := db.FindUserByID(id)
	if u == nil {
		return fmt.Errorf("Can't find user with id = %d", id)
	}

	if updateCMD.Username != "" {
		u.Username = updateCMD.Username
	}
	if updateCMD.Password != "" {
		u.Password = updateCMD.Password
	}
	if updateCMD.Email != "" {
		u.Email = updateCMD.Email
	}
	if updateCMD.FirstName != "" {
		u.FirstName = updateCMD.FirstName
	}
	if updateCMD.LastName != "" {
		u.LastName = updateCMD.LastName
	}

	if updateCMD.Expired != nil {
		u.Expired = *updateCMD.Expired
	}

	if updateCMD.Roles != nil {
		u.Roles = *updateCMD.Roles
	}

	return nil
}
