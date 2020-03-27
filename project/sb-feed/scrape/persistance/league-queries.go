package persistance

import (
	"database/sql"
	"fmt"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrLeagueCreation comment
const ErrLeagueCreation = "Can't create a league from entity"

const findLeagueByIDQuery = `
select l.ID, l.Name, l.Country from League as l
where l.IsDeleted = 0 and l.ID = ?
`

// FindLeagueByID comment
func (dl *DataLayer) FindLeagueByID(id int64) (*entities.League, error) {
	row := dl.db.Connection().QueryRow(findLeagueByIDQuery, id)
	return extratLeagueFromRow(row)
}

const findLeagueQuery = `
select l.ID, l.Name, l.Country from League as l
where l.IsDeleted = 0 and l.Name = ? and l.SportID = ? and l.Country = ?
`

// FindLeague comment
func (dl *DataLayer) FindLeague(name string, sportID int64, country string) (*entities.League, error) {
	row := dl.db.Connection().QueryRow(findLeagueQuery, name, sportID, country)
	return extratLeagueFromRow(row)
}

func extratLeagueFromRow(row *sql.Row) (*entities.League, error) {
	var id int64
	var name string
	var country string
	err := row.Scan(&id, &name, &country)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	league := entities.NewLeague(id, name, nil, country)

	return league, nil
}

const createLeagueQuery = `
	insert into League (Name, SportID, Country) values (?, ?, ?)
`

// CreateLeague creates a league and sets the ID.
// NOTE: Assumes that the related sport exists!
func (dl *DataLayer) CreateLeague(league *entities.League) error {
	if league == nil || league.Name == "" || len(league.Name) > 254 {
		return fmt.Errorf(ErrLeagueCreation)
	}
	if league.ID > 0 {
		infra.Warn(infra.WarningColor, "League might already exist in DB")
	}

	result, err := dl.db.Connection().
		Exec(createLeagueQuery, league.Name, league.Sport.ID, league.Country)
	if err != nil {
		return err
	}

	league.ID, err = result.LastInsertId()
	return err
}

const deleteLeagueQuery = `
	update League
	set IsDeleted = 1
	where ID = ?
`

// DeleteLeague soft delete from database
func (dl *DataLayer) DeleteLeague(id int64) error {
	_, err := dl.db.Connection().Exec(deleteLeagueQuery, id)
	return err
}

const allLeagues = `
	select l.ID, l.Name, l.Country, s.ID as SportID, s.Name as SportName from League as l
	join Sport as s on s.ID = l.SportID and s.IsDeleted = 0
	where l.IsDeleted = 0
`

// AllLeagues commnet
func (dl *DataLayer) AllLeagues() ([]*entities.League, error) {
	ret := make([]*entities.League, 0)
	rows, err := dl.db.Connection().Query(allLeagues)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		var country string
		var sportID int64
		var sportName string
		err := rows.Scan(&id, &name, &country, &sportID, &sportName)
		if err != nil {
			return nil, err
		}

		s := entities.NewLeague(id, name, entities.NewSport(sportID, sportName), country)
		ret = append(ret, s)
	}

	return ret, nil
}
