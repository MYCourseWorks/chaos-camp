package entities

// Sport represents a sport
type Sport struct {
	Metadata `json:"-"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
}

// NewSport initializes a sport
func NewSport(id int64, name string) *Sport {
	return &Sport{ID: id, Name: name}
}
