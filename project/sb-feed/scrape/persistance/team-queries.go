package persistance

import (
	"database/sql"
	"fmt"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrTeamCreation comment
const ErrTeamCreation = "Can't create a team from entity"

const findTeamByIDQuery = `
select t.ID, t.Name from Team as t
where t.IsDeleted = 0 and t.ID = ?
`

// FindTeamByID comment
func (dl *DataLayer) FindTeamByID(id int64) (*entities.Team, error) {
	row := dl.db.Connection().QueryRow(findTeamByIDQuery, id)
	return extractTeamFromRow(row)
}

const findTeamQuery = `
select t.ID, t.Name from Team as t
where t.IsDeleted = 0 and t.Name = ? and t.LeagueID = ?
`

// FindTeam comment
func (dl *DataLayer) FindTeam(name string, leagueID int64) (*entities.Team, error) {
	row := dl.db.Connection().QueryRow(findTeamQuery, name, leagueID)
	return extractTeamFromRow(row)
}

func extractTeamFromRow(row *sql.Row) (*entities.Team, error) {
	var id int64
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	team := entities.NewTeam(id, name, nil)

	return team, nil
}

const createTeamQuery = `
	insert into Team (Name, LeagueID) values (?, ?)
`

// CreateTeam creates a team and sets the ID.
// NOTE: Assumes that the related league exists!
func (dl *DataLayer) CreateTeam(team *entities.Team) error {
	if team == nil || team.Name == "" || len(team.Name) > 254 {
		return fmt.Errorf(ErrTeamCreation)
	}
	if team.ID > 0 {
		infra.Warn(infra.WarningColor, "Team might already exist in DB")
	}

	result, err := dl.db.Connection().Exec(createTeamQuery, team.Name, team.League.ID)
	if err != nil {
		return err
	}

	team.ID, err = result.LastInsertId()
	return err
}

const deleateTeamQuery = `
	update Team
	set IsDeleted = 1
	where ID = ?
`

// DeleteTeam soft delete from database
func (dl *DataLayer) DeleteTeam(id int64) error {
	_, err := dl.db.Connection().Exec(deleateTeamQuery, id)
	return err
}
