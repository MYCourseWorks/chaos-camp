package siterepo

import "github.com/MartinNikolovMarinov/sb-users/models"

// SiteRepo comment
type SiteRepo interface {
	FindByID(id int64) (*models.Site, error)
	SiteExists(name, country string) (bool, error)
	All(name string) ([]models.Site, error)
	Create(name string, country string) (int64, error)
}
