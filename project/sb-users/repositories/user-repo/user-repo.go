package userrepo

import "github.com/MartinNikolovMarinov/sb-users/models"

// UserRepo comment
type UserRepo interface {
	AllPerSite(siteID int64) ([]models.User, error)
	Find(name string) (*models.User, error)
	Create(name, password string, roles models.RoleFlags, siteID int64) (*models.User, error)
	UpdateUsername(userID int64, userName string) error
	UpdatePassword(userID int64, password string) error
	MatchPassword(password, hashedPass string) error
	Delete(userID int64) error
}
