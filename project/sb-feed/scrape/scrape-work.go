package scrape

import (
	"errors"

	"github.com/MartinNikolovMarinov/sb-feed/scrape/persistance"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

func persistGames(dl *persistance.DataLayer, games []*entities.Game) {
	var err error

	for i := 0; i < len(games); i++ {
		g := games[i]

		err = findOrCreateSport(dl, g.Sport)
		if err != nil {
			infra.Error("%s\n", err.Error())
			continue
		}

		err = findOrCreateLeague(dl, g.League)
		if err != nil {
			infra.Error("%s\n", err.Error())
			continue
		}

		for _, t := range g.Teams {
			err = findOrCreateTeam(dl, t)
			if err != nil {
				infra.Error("%s\n", err.Error())
				continue
			}
		}

		err = findOrCreateGame(dl, g)
		if err != nil {
			infra.Error("%s\n", err.Error())
			continue
		}

		for _, e := range g.Events {
			err = findOrCreateEvent(dl, e)
			if err != nil {
				infra.Error("%s\n", err.Error())
				continue
			}

			for _, l := range e.Lines {
				err = findOrCreateLine(dl, l)
				if err != nil {
					infra.Error("%s\n", err.Error())
					continue
				}

				for _, o := range l.Odds {
					err = findOrCreateOdd(dl, o)
					if err != nil {
						infra.Error("%s\n", err.Error())
						continue
					}
				}
			}
		}
	}
}

const errScrappedSportNotValid = "Scrapped sport is not valid"

func findOrCreateSport(dl *persistance.DataLayer, sport *entities.Sport) error {
	if sport == nil || sport.Name == "" {
		return errors.New(errScrappedSportNotValid)
	}

	s, err := dl.FindSport(sport.Name)
	if err != nil {
		return err
	}

	if s == nil {
		err := dl.CreateSport(sport)
		if err != nil {
			return err
		}
	} else {
		sport.ID = s.ID
	}

	return nil
}

const errScrappedLeagueNotValid = "Scrapped league is not valid"

func findOrCreateLeague(dl *persistance.DataLayer, league *entities.League) error {
	if league == nil || league.Name == "" {
		return errors.New(errScrappedLeagueNotValid)
	}

	l, err := dl.FindLeague(league.Name, league.Sport.ID, league.Country)
	if err != nil {
		return err
	}

	if l == nil {
		err := dl.CreateLeague(league)
		if err != nil {
			return err
		}
	} else {
		league.ID = l.ID
	}

	return nil
}

const errScrappedTeamNotValid = "Scrapped team is not valid"

func findOrCreateTeam(dl *persistance.DataLayer, team *entities.Team) error {
	if team == nil || team.Name == "" || team.League == nil {
		return errors.New(errScrappedTeamNotValid)
	}

	t, err := dl.FindTeam(team.Name, team.League.ID)
	if err != nil {
		return err
	}

	if t == nil {
		err := dl.CreateTeam(team)
		if err != nil {
			return err
		}
	} else {
		team.ID = t.ID
	}

	return nil
}

const errScrappedGame = "Scrapped game is not valid"

func findOrCreateGame(dl *persistance.DataLayer, game *entities.Game) error {
	if game == nil || game.Name == "" || game.League == nil || game.Sport == nil {
		return errors.New(errScrappedGame)
	}

	t, err := dl.FindGame(game.Name, game.Sport.ID, game.League.ID)
	if err != nil {
		return err
	}

	if t == nil {
		err := dl.CreateGame(game)
		if err != nil {
			return err
		}
	} else {
		game.ID = t.ID
	}

	return nil
}

const errScrappedEvent = "Scrapped event is not valid"

func findOrCreateEvent(dl *persistance.DataLayer, event *entities.GameEvent) error {
	if event == nil || event.RelatedGame == nil {
		return errors.New(errScrappedEvent)
	}

	ev, err := dl.FindEvent(&event.Date, event.EventType, event.RelatedGame.ID)
	if err != nil {
		return err
	}

	if ev == nil {
		err := dl.CreateEvent(event)
		if err != nil {
			return err
		}
	} else {
		event.ID = ev.ID
	}

	return nil
}

const errScrappedLine = "Scrapped line is not valid"

func findOrCreateLine(dl *persistance.DataLayer, line *entities.Line) error {
	if line == nil || line.Description == "" {
		return errors.New(errScrappedLine)
	}

	ev, err := dl.FindLine(line.Description, line.Event.ID)
	if err != nil {
		return err
	}

	if ev == nil {
		err := dl.CreateLine(line)
		if err != nil {
			return err
		}
	} else {
		line.ID = ev.ID
	}

	return nil
}

const errScrappedOdd = "Scrapped line is not valid"

func findOrCreateOdd(dl *persistance.DataLayer, odd *entities.Odd) error {
	if odd == nil || odd.Source == "" || odd.Values == nil || len(odd.Values) == 0 {
		return errors.New(errScrappedOdd)
	}

	ev, err := dl.FindOdd(odd.Source, odd.Line.ID)
	if err != nil {
		return err
	}

	if ev == nil {
		err := dl.CreateOdd(odd)
		if err != nil {
			return err
		}
	} else {
		odd.ID = ev.ID
	}

	return nil
}
