package entities

// League represents a league
type League struct {
	Metadata `json:"-"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Sport    *Sport `json:"sport"`
	Country  string `json:"country"`
}

// NewLeague TODO:
func NewLeague(id int64, name string, sport *Sport, country string) *League {
	return &League{ID: id, Name: name, Sport: sport, Country: country}
}
