package sportshandler

import (
	"net/http"

	sportsrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/sports-repo"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// All get a list of all valid sports
func All(sr sportsrepo.SportsRepo) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		var sports []entities.Sport

		sports, err = sr.All()
		if err != nil {
			return infra.NewHTTPError(err, 500, "All sports query failed")
		}

		err = infra.WriteResponseJSON(w, sports, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
