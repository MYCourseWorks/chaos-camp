package leaguehandler

import (
	"net/http"

	leaguesrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/leagues-repo"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// All comment
func All(lr leaguesrepo.LeaguesRepo) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		var sports []entities.League

		sports, err = lr.All()
		if err != nil {
			return infra.NewHTTPError(err, 500, "All leagues query failed")
		}

		err = infra.WriteResponseJSON(w, sports, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
