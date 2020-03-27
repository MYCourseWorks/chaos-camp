package sportsrepo

import "github.com/MartinNikolovMarinov/sb-infra/entities"

// SportsRepo comment
type SportsRepo interface {
	All() ([]entities.Sport, error)
}
