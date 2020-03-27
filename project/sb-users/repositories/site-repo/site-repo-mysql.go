package siterepo

import (
	"database/sql"
	"log"

	"github.com/MartinNikolovMarinov/sb-users/models"
)

type siteRepoMysql struct {
	db *sql.DB
}

// NewMysqlSiteRepo is a SiteRepo constructor
func NewMysqlSiteRepo(connectionString string) SiteRepo {
	repo := &siteRepoMysql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (repo *siteRepoMysql) FindByID(id int64) (*models.Site, error) {
	const findSiteByIDQuery = `
		select ID, Name, Country from Site
		where IsDeleted = 0 and ID = ?
	`
	var name string
	var country string
	var err = repo.db.QueryRow(findSiteByIDQuery, id).
		Scan(&id, &name, &country)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var site = &models.Site{ID: id, Name: name, Country: country}
	return site, nil
}

func (repo *siteRepoMysql) SiteExists(name, country string) (bool, error) {
	const findSiteQuery = `
		select Count(*) from Site
		where IsDeleted = 0 and Name = ? and Country = ?
	`

	var count int
	var err = repo.db.QueryRow(findSiteQuery, name, country).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (repo *siteRepoMysql) All(name string) ([]models.Site, error) {
	const allSitesQuery = `
		select ID, Name, Country from Site
		where IsDeleted = 0
	`
	sites := make([]models.Site, 0)
	rows, err := repo.db.Query(allSitesQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name string
		var country string
		err := rows.Scan(&id, &name, &country)
		if err != nil {
			return nil, err
		}

		var s = models.Site{ID: id, Name: name, Country: country}
		sites = append(sites, s)
	}

	return sites, nil
}

func (repo *siteRepoMysql) Create(name string, country string) (int64, error) {
	const createSiteQuery = `
		insert into Site (Name, Country) values (?, ?)
	`

	result, err := repo.db.Exec(createSiteQuery, name, country)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
