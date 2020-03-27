package persistance

import (
	"database/sql"
	"fmt"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrGameCreation comment
const ErrGameCreation = "Can't create a game from entity"

const findGameByIDQuery = `
select g.ID, g.Name from Game as g
where g.IsDeleted = 0 and g.ID = ?
`

// FindGameByID comment
func (dl *DataLayer) FindGameByID(id int64) (*entities.Game, error) {
	row := dl.db.Connection().QueryRow(findGameByIDQuery, id)
	return extratGameFromRow(row)
}

const findGameQuery = `
select g.ID, g.Name from Game as g
where g.IsDeleted = 0 and g.Name = ? and g.SportID = ? and g.LeagueID = ?
`

// FindGame comment
func (dl *DataLayer) FindGame(name string, sportID int64, leagueID int64) (*entities.Game, error) {
	row := dl.db.Connection().QueryRow(findGameQuery, name, sportID, leagueID)
	return extratGameFromRow(row)
}

func extratGameFromRow(row *sql.Row) (*entities.Game, error) {
	var id int64
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	game := entities.NewGame(id, name, nil, nil)

	return game, nil
}

const createGameQuery = `
	insert into Game (Name, SportID, LeagueID) values (?, ?, ?)
`
const insertGame2TeamQuery = `
	insert into Game2Team (GameID, TeamID) values (?, ?)
`

// CreateGame creates a game and sets the ID.
// NOTE: Assumes that the related sport, league and teams exists!
func (dl *DataLayer) CreateGame(game *entities.Game) error {
	if game == nil || game.Name == "" || len(game.Name) > 254 {
		return fmt.Errorf(ErrGameCreation)
	}
	if game.ID > 0 {
		infra.Warn(infra.WarningColor, "Game might already exist in DB")
	}

	var err error
	var result sql.Result
	result, err = dl.db.Connection().Exec(createGameQuery, game.Name, game.Sport.ID, game.League.ID)
	if err != nil {
		return err
	}

	game.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	for _, t := range game.Teams {
		result, err = dl.db.Connection().Exec(insertGame2TeamQuery, game.ID, t.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

const deleteGameQuery = `
	update Game
	set IsDeleted = 1
	where ID = ?
`

// DeleteGame soft delete from database
func (dl *DataLayer) DeleteGame(id int64) error {
	_, err := dl.db.Connection().Exec(deleteGameQuery, id)
	return err
}
