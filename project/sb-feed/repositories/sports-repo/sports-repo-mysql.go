package sportsrepo

import (
	"database/sql"
	"log"

	repo "github.com/MartinNikolovMarinov/sb-feed/repositories"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
)

type sportRepoMysql struct {
	db *sql.DB
}

// NewMysqlRepo is a SiteRepo constructor
func NewMysqlRepo(connectionString string) SportsRepo {
	repo := &sportRepoMysql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (sr *sportRepoMysql) All() ([]entities.Sport, error) {
	const allSportsQuery = `
		select s.ID, s.Name from Sport as s
		where s.IsDeleted = 0
	`

	ret := make([]entities.Sport, 0)
	rows, err := sr.db.Query(allSportsQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		err = repo.ScanRow(rows, &id, &name)
		if err != nil {
			return nil, err
		}

		sport := entities.NewSport(id, name)
		ret = append(ret, *sport)
	}

	return ret, nil
}
