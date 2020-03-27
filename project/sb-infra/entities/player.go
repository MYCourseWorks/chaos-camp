package entities

// Player represents a sports competitor
type Player struct {
	Metadata `json:"-"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
}
