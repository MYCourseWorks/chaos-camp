package leaguesrepo

import "github.com/MartinNikolovMarinov/sb-infra/entities"

// LeaguesRepo comment
type LeaguesRepo interface {
	All() ([]entities.League, error)
}
