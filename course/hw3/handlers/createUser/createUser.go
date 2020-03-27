package createuser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/homework/hw3/handlers"
	"github.com/homework/hw3/infra"
	user "github.com/homework/hw3/models"
)

// CMD comment
type CMD struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Password  string `json:"password"`
}

// Handle commnet
func Handle(db infra.UserDatabase) handlers.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c CMD
		// TODO: json Decoder doesn't check the size of the request body.
		// Which allows for DDOS attacks.
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&c)
		if err != nil {
			return infra.NewHTTPError(err, 400, "request body is not valid")
		}

		// TODO: perform input validation ?
		err = validateCMD(&c)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		create(db, &c)

		w.Header().Add("Location", fmt.Sprintf("user/%d", len(db)-1))
		w.WriteHeader(201)
		return nil
	}

	return ret
}

func validateCMD(cmd *CMD) error {
	var bulder strings.Builder

	if cmd.Username == "" {
		if _, err := bulder.WriteString("username not set;"); err != nil {
			return err
		}
	}
	if cmd.Email == "" {
		if _, err := bulder.WriteString("email not set;"); err != nil {
			return err
		}
	}
	if cmd.FirstName == "" {
		if _, err := bulder.WriteString("firstName not set;"); err != nil {
			return err
		}
	}
	if cmd.LastName == "" {
		if _, err := bulder.WriteString("lastName not set;"); err != nil {
			return err
		}
	}
	if cmd.Password == "" {
		if _, err := bulder.WriteString("password not set;"); err != nil {
			return err
		}
	}

	if bulder.Len() > 0 {
		return errors.New(bulder.String())
	}

	return nil
}

func create(db infra.UserDatabase, creteCMD *CMD) {
	if db == nil {
		panic("db is nil")
	}

	u := user.NewUser(db.NextID(),
		creteCMD.Username,
		creteCMD.Password,
		creteCMD.Email,
		creteCMD.FirstName,
		creteCMD.LastName)

	db[u.ID] = u
}
