package entities

import "time"

// EventType TODO:
type EventType int

const (
	// None no event type set
	None EventType = iota
	// T1x2 1x2 game event
	T1x2 EventType = iota
)

// GameEvent TODO:
type GameEvent struct {
	Metadata    `json:"-"`
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	RelatedGame *Game     `json:"relatedgame"`
	EventType   EventType `json:"eventtype"`
	Lines       []*Line   `json:"lines"`
}

// NewGameEvent creates a new instane of GameEvent
func NewGameEvent(id int64, date time.Time, relatedGame *Game, evType EventType) *GameEvent {
	return &GameEvent{
		ID:          id,
		Date:        date,
		RelatedGame: relatedGame,
		EventType:   evType,
		Lines:       make([]*Line, 0),
	}
}

// AddLine adds a new line
func (ge *GameEvent) AddLine(line *Line) {
	if ge.Lines == nil {
		ge.Lines = make([]*Line, 0)
	}
	ge.Lines = append(ge.Lines, line)
}
