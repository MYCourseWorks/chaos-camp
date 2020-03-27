package entities

// OddFormat format types
type OddFormat int

const (
	// BritishOdds odds 	-> 1/4, 1/3, 6/5
	BritishOdds OddFormat = iota
	// EuropeanOdds odds 	-> 1.25, 1.33, 2.20
	EuropeanOdds OddFormat = iota
	// AmericanOdds odds 	-> -400, -300, 120
	AmericanOdds OddFormat = iota
	// HongKongOdds odds 	-> 0.25, 0.33, 1.20
	HongKongOdds OddFormat = iota
)

// LineType FIXME: deside what line type stands for concretely.
type LineType int

const (
	// MoneyLine is a when you bet to pick a team to win straight-up.
	// It is the line type every event should have.
	MoneyLine LineType = iota
)

// Line represents somthing you can bet on
type Line struct {
	Metadata    `json:"-"`
	ID          int64      `json:"id"`
	Odds        []*Odd     `json:"odds"`
	OddFormat   OddFormat  `json:"oddformat"`
	LineType    LineType   `json:"linetype"`
	Event       *GameEvent `json:"event"`
	Description string     `json:"description"`
}

// NewLine creates and instance of Line
func NewLine(id int64,
	lineType LineType,
	description string,
	event *GameEvent) *Line {

	return &Line{
		ID:          id,
		OddFormat:   BritishOdds,
		LineType:    lineType,
		Odds:        make([]*Odd, 0),
		Description: description,
		Event:       event,
	}
}

// AddOdd adds odds to the line
func (l *Line) AddOdd(odd *Odd) {
	if l.Odds == nil {
		l.Odds = make([]*Odd, 0)
	}
	l.Odds = append(l.Odds, odd)
}
