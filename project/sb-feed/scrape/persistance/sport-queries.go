package persistance

import (
	"database/sql"
	"fmt"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrSportCreation comment
const ErrSportCreation = "Can't create a sport from entity"

const findSportByIDQuery = `
select s.ID, s.Name from Sport as s
where s.IsDeleted = 0 and s.ID = ?
`

// FindSportByID comment
func (dl *DataLayer) FindSportByID(id int64) (*entities.Sport, error) {
	row := dl.db.Connection().QueryRow(findSportByIDQuery, id)
	return extratSportFromRow(row)
}

const findSportQuery = `
select s.ID, s.Name from Sport as s
where s.IsDeleted = 0 and s.Name = ?
`

// FindSport comment
func (dl *DataLayer) FindSport(name string) (*entities.Sport, error) {
	row := dl.db.Connection().QueryRow(findSportQuery, name)
	return extratSportFromRow(row)
}

func extratSportFromRow(row *sql.Row) (*entities.Sport, error) {
	var id int64
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	sport := entities.NewSport(id, name)

	return sport, nil
}

const createSportQuery = `
	insert into Sport (Name) values (?)
`

// CreateSport creates a sport and sets the ID
func (dl *DataLayer) CreateSport(sport *entities.Sport) error {
	if sport == nil || sport.Name == "" || len(sport.Name) > 254 {
		return fmt.Errorf(ErrSportCreation)
	}
	if sport.ID > 0 {
		infra.Warn(infra.WarningColor, "Sport might already exist in DB")
	}

	result, err := dl.db.Connection().Exec(createSportQuery, sport.Name)
	if err != nil {
		return err
	}

	sport.ID, err = result.LastInsertId()
	return err
}

const deleteSportQuery = `
	update Sport
	set IsDeleted = 1
	where ID = ?
`

// DeleteSport soft delete from database
func (dl *DataLayer) DeleteSport(id int64) error {
	_, err := dl.db.Connection().Exec(deleteSportQuery, id)
	return err
}

const allSportsQuery = `
	select s.ID, s.Name from Sport as s
	where s.IsDeleted = 0
`

// AllSports comment
func (dl *DataLayer) AllSports() ([]*entities.Sport, error) {
	ret := make([]*entities.Sport, 0)
	rows, err := dl.db.Connection().Query(allSportsQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		s := entities.NewSport(id, name)
		ret = append(ret, s)
	}

	return ret, nil
}
