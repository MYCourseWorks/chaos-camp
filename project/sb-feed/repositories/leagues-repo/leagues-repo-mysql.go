package leaguesrepo

import (
	"database/sql"
	"log"

	repo "github.com/MartinNikolovMarinov/sb-feed/repositories"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
)

type leaguesRepoMysql struct {
	db *sql.DB
}

// NewMysqlRepo mysql repo
func NewMysqlRepo(connectionString string) LeaguesRepo {
	repo := &leaguesRepoMysql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (sr *leaguesRepoMysql) All() ([]entities.League, error) {
	const allLeagues = `
		select l.ID, l.Name, l.Country, s.ID as SportID, s.Name as SportName from League as l
		join Sport as s on s.ID = l.SportID and s.IsDeleted = 0
		where l.IsDeleted = 0
	`

	ret := make([]entities.League, 0)
	rows, err := sr.db.Query(allLeagues)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, sportID int64
		var name, country, sportName string
		err = repo.ScanRow(rows, &id, &name, &country, &sportID, &sportName)
		if err != nil {
			return nil, err
		}

		s := entities.NewSport(sportID, sportName)
		league := entities.NewLeague(id, name, s, country)
		ret = append(ret, *league)
	}

	return ret, nil
}
