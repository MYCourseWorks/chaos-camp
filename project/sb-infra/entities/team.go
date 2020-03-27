package entities

// Team represents a team
type Team struct {
	Metadata `json:"-"`
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	League   *League   `json:"league"`
	Players  []*Player `json:"players"`
}

// NewTeam creats a new team instnace
func NewTeam(id int64, name string, league *League) *Team {
	return &Team{
		ID:      id,
		Name:    name,
		League:  league,
		Players: make([]*Player, 0),
	}
}

// AddPlayer adds a new player to the team
func (t *Team) AddPlayer(p *Player) {
	if t.Players == nil {
		t.Players = make([]*Player, 0)
	}
	t.Players = append(t.Players, p)
}
