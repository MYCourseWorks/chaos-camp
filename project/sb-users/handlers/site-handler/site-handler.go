package sitehandler

import (
	"net/http"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	siterepo "github.com/MartinNikolovMarinov/sb-users/repositories/site-repo"
	"github.com/go-playground/validator"
)

// CreateCMD comment
type CreateCMD struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Country string `json:"country" validate:"required,min=2,max=255"`
}

// Create commnet
func Create(sr siterepo.SiteRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c CreateCMD
		var err error
		err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return err
		}

		exists, err := sr.SiteExists(c.Name, c.Country)
		if err != nil {
			return err
		}
		if exists {
			return infra.NewHTTPError(err, 400, "Site Already Exists")
		}

		var id int64
		id, err = sr.Create(c.Name, c.Country)
		if err != nil {
			return err
		}

		resp := map[string]int64{"id": id}
		err = infra.WriteResponseJSON(w, &resp, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
