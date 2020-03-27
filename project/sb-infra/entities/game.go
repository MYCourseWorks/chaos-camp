package entities

// Game TODO:
type Game struct {
	Metadata `json:"-"`
	ID       int64        `json:"id"`
	Name     string       `json:"name"`
	Sport    *Sport       `json:"sport"`
	League   *League      `json:"league"`
	Teams    []*Team      `json:"teams"`
	Events   []*GameEvent `json:"events"`
}

// NewGame creates a new instane of Game
func NewGame(id int64, name string, sport *Sport, league *League) *Game {
	return &Game{
		ID:     id,
		Name:   name,
		Sport:  sport,
		League: league,
		Teams:  make([]*Team, 0),
		Events: make([]*GameEvent, 0),
	}
}

// AddTeam adds a new team
func (g *Game) AddTeam(team *Team) {
	if g.Teams == nil {
		g.Teams = make([]*Team, 0)
	}
	g.Teams = append(g.Teams, team)
}

// AddEvent adds a new event
func (g *Game) AddEvent(event *GameEvent) {
	if g.Events == nil {
		g.Events = make([]*GameEvent, 0)
	}
	g.Events = append(g.Events, event)
}

// Equals compares two gamse
func (g *Game) Equals(other *Game) bool {
	if g.ID > 0 && other.ID > 0 {
		return g.ID == other.ID
	}

	ret := (g.Name == other.Name) &&
		(g.Sport.Name == other.Sport.Name) &&
		(g.League.Name == other.League.Name)

	return ret
}
