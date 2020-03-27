package entities

// Site TODO:
type Site struct {
	Metadata `json:"-"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
}
