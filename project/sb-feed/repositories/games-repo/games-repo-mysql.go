package gamesrepo

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	repo "github.com/MartinNikolovMarinov/sb-feed/repositories"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/shopspring/decimal"
)

type gamesRepoMysql struct {
	db         *sql.DB
	dateFormat string
}

// NewMysqlRepo mysql repo
func NewMysqlRepo(connectionString string) GamesRepo {
	repo := &gamesRepoMysql{dateFormat: "2006-01-02 15:04:05"}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

// AllPerLeague comment
func (gr *gamesRepoMysql) AllPerLeague(leagueID int64) ([]DaoEvent, error) {
	const gamesPerLeagueQuery = `
		select
			e.ID as EventID,
			e.Date as EventDate,
			e.EventType as EventType,
			e.IsFrozen as EventIsFrozen,
			e.RelatedGameID as RelatedGameID,
			g.Name as RelatedGameName,
			l.ID as LeagueID,
			l.Name as LeagueName,
			l.Country as LeagueCountry
		from GameEvent as e
		join Game as g on g.ID = e.RelatedGameID and g.IsDeleted = 0
		join League as l on l.ID = g.LeagueID and l.IsDeleted = 0
		where
			e.IsDeleted = 0 and
			LeagueID = ?
	`

	rows, err := gr.db.Query(gamesPerLeagueQuery, leagueID)
	if err != nil {
		return nil, err
	}

	return gr.parseEvents(rows)
}

// All commnet
func (gr *gamesRepoMysql) All() ([]DaoEvent, error) {
	const allGamesQuery = `
		select
			e.ID as EventID,
			e.Date as EventDate,
			e.EventType as EventType,
			e.IsFrozen as EventIsFrozen,
			e.RelatedGameID as RelatedGameID,
			g.Name as RelatedGameName,
			l.ID as LeagueID,
			l.Name as LeagueName,
			l.Country as LeagueCountry
		from GameEvent as e
		join Game as g on g.ID = e.RelatedGameID and g.IsDeleted = 0
		join League as l on l.ID = g.LeagueID and l.IsDeleted = 0
		where
			e.IsDeleted = 0
	`

	rows, err := gr.db.Query(allGamesQuery)
	if err != nil {
		return nil, err
	}

	return gr.parseEvents(rows)
}

func (gr *gamesRepoMysql) parseEvents(rows *sql.Rows) ([]DaoEvent, error) {
	var (
		err              error
		eventID          int64
		leagueIDFromDB   int64
		relatedGameID    int64
		eventDateStr     string
		relatedGameName  string
		leagueName       string
		leagueCountry    string
		eventType        int
		eventIsFrozenBit []byte
		daoEvents        []DaoEvent = make([]DaoEvent, 0)
	)

	for rows.Next() {
		err = repo.ScanRow(
			rows,
			&eventID, &eventDateStr, &eventType, &eventIsFrozenBit,
			&relatedGameID, &relatedGameName,
			&leagueIDFromDB, &leagueName, &leagueCountry,
		)
		if err != nil {
			return nil, err
		}

		var eventDate time.Time
		eventDate, err = time.Parse(gr.dateFormat, eventDateStr)
		if err != nil {
			return nil, err
		}

		isFrozen := eventIsFrozenBit[0] == 1

		ev := DaoEvent{
			ID:        eventID,
			EventDate: eventDate,
			EventType: entities.EventType(eventType),
			RelatedGame: DaoRelatedGame{
				ID:   relatedGameID,
				Name: relatedGameName,
			},
			League: DaoLeague{
				ID:      leagueIDFromDB,
				Name:    leagueName,
				Country: leagueCountry,
			},
			IsFrozen: isFrozen,
		}

		ev.Lines, err = gr.getLinesForEvent(ev.ID)
		if err != nil {
			return nil, err
		}

		daoEvents = append(daoEvents, ev)
	}

	return daoEvents, nil
}

func (gr *gamesRepoMysql) getLinesForEvent(eventID int64) ([]DaoLine, error) {
	const getLinesForEventQuery = `
		select
			l.ID as LineID,
			l.LineType as LineType,
			l.Description as LineDescription,
			o.ID as OddID,
			o.Source as OddSource,
			o.Values as OddValues
		from Line as l
		join GameEvent as e on e.ID = l.EventID and e.IsDeleted = 0
		join Odd as o on o.LineID = l.ID and o.IsDeleted = 0
		where
			l.IsDeleted = 0 and e.ID = ?
	`

	linesMap := make(map[int64]DaoLine)            // lineID -> DaoLine
	oddsMap := make(map[int64][][]decimal.Decimal) // lineID -> odds
	rows, err := gr.db.Query(getLinesForEventQuery, eventID)
	if err != nil {
		return nil, err
	}

	var (
		lineID, oddID                         int64
		lineType                              int
		lineDescription, oddSource, oddValues string
	)

	for rows.Next() {
		err = repo.ScanRow(
			rows,
			&lineID, &lineType, &lineDescription,
			&oddID, &oddSource, &oddValues,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := linesMap[lineID]; !ok {
			linesMap[lineID] = DaoLine{
				ID:         lineID,
				LineType:   entities.LineType(lineType),
				Descripion: lineDescription,
				Odds:       make([]string, 0),
			}
		}

		if _, ok := oddsMap[lineID]; !ok {
			oddsMap[lineID] = make([][]decimal.Decimal, 0)
		}

		split := strings.Split(oddValues, ":")
		odds := make([]decimal.Decimal, 0)
		for i := 0; i < len(split); i++ {
			var val decimal.Decimal
			val, err := decimal.NewFromString(split[i])
			if err != nil {
				return nil, err
			}
			odds = append(odds, val)
		}

		oddsMap[lineID] = append(oddsMap[lineID], odds)
	}

	ret := make([]DaoLine, len(linesMap))
	i := 0
	for _, dl := range linesMap {
		if odds, ok := oddsMap[dl.ID]; ok && len(odds) > 0 {
			var aggr = make([]decimal.Decimal, len(odds[0]))
			for _, oddLine := range odds {
				for j, o := range oddLine {
					aggr[j] = aggr[j].Add(o)
				}
			}

			var aggrLen = int32(len(aggr))
			var oddsLen = int32(len(odds))
			for j := 0; int32(j) < aggrLen; j++ {
				el := aggr[j]
				el = el.Div(decimal.NewFromInt32(oddsLen))
				dl.Odds = append(dl.Odds, el.String())
			}

			ret[i] = dl
			i++
		}
	}

	return ret, nil
}

// Freeze comment
func (gr *gamesRepoMysql) Freeze(eventID int64) error {
	const freezeEventQuery = `
		update GameEvent
		set IsFrozen = CASE WHEN IsFrozen = 0 THEN 1 ELSE 0 END
		where ID = ?
	`

	result, err := gr.db.Exec(freezeEventQuery, eventID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil || n != 1 {
		return errors.New("Couldn't find event")
	}

	return nil
}
