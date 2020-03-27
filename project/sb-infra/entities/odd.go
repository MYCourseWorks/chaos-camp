package entities

// Odd represents a single odd from a specific source
type Odd struct {
	Metadata `json:"-"`
	ID       int64     `json:"id"`
	Source   string    `json:"source"`
	Values   []float32 `json:"values"`
	Line     *Line     `json:"line"`
}

// NewOdd creates a new instance of Odd
func NewOdd(id int64, source string, values []float32, line *Line) *Odd {
	return &Odd{
		ID:     id,
		Source: source,
		Values: values,
		Line:   line,
	}
}
