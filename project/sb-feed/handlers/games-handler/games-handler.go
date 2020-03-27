package gameshandler

import (
	"net/http"
	"strconv"

	gamesrepo "github.com/MartinNikolovMarinov/sb-feed/repositories/games-repo"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

// QueryOptions comment
type QueryOptions struct {
	LeagueID *int64
}

// All comment
func All(gr gamesrepo.GamesRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error

		var leagueID = r.URL.Query().Get("leagueID")
		var opt = QueryOptions{LeagueID: nil}
		if leagueID != "" {
			var lid int64
			lid, err = strconv.ParseInt(leagueID, 10, 64)
			if err != nil {
				return infra.NewHTTPError(err, 400, "Invalid query parameter leagues")
			}
			opt.LeagueID = &lid
		}

		var games []gamesrepo.DaoEvent
		games, err = getGames(gr, &opt)
		if err != nil {
			return infra.NewHTTPError(err, 500, "All games query failed")
		}

		err = infra.WriteResponseJSON(w, games, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

func getGames(gr gamesrepo.GamesRepo, opt *QueryOptions) ([]gamesrepo.DaoEvent, error) {
	if opt.LeagueID != nil {
		return gr.AllPerLeague(*opt.LeagueID)
	} else {
		return gr.All()
	}
}

// Freeze comment
func Freeze(gr gamesrepo.GamesRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error

		var eventID = chi.URLParam(r, "eventID")
		if eventID == "" {
			return infra.NewHTTPError(err, 400, "Invalid query parameter eventID")
		}

		var eid int64
		eid, err = strconv.ParseInt(eventID, 10, 64)
		if err != nil {
			return infra.NewHTTPError(err, 400, "Invalid query parameter eventID")
		}

		err = gr.Freeze(eid)
		if err != nil {
			return infra.NewHTTPError(err, 500, err.Error())
		}

		err = infra.WriteResponseJSON(w, map[string]bool{"ok": true}, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
