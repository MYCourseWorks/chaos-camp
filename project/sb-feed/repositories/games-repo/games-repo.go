package gamesrepo

import (
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
)

// DaoRelatedGame comment
type DaoRelatedGame struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DaoLeague comment
type DaoLeague struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// DaoLine commnet
type DaoLine struct {
	ID         int64             `json:"id"`
	Odds       []string          `json:"odds"` // aggregated odds from all sources as a string
	LineType   entities.LineType `json:"lineType"`
	Descripion string            `json:"descripion"`
}

// DaoEvent commet
type DaoEvent struct {
	ID          int64              `json:"id"`
	EventDate   time.Time          `json:"eventDate"`
	EventType   entities.EventType `json:"eventType"`
	RelatedGame DaoRelatedGame     `json:"relatedGame"`
	League      DaoLeague          `json:"league"`
	Lines       []DaoLine          `json:"lines"`
	IsFrozen    bool               `json:"isFrozen"`
}

// GamesRepo comment
type GamesRepo interface {
	All() ([]DaoEvent, error)
	AllPerLeague(leagueID int64) ([]DaoEvent, error)
	Freeze(eventID int64) error
}
