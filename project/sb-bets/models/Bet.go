package models

// Bet comment
type Bet struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"userID"`
	UserName        string `json:"userName"`
	Value           string `json:"value"`
	IsPayed         bool   `json:"isPayed"`
	BetOnIndex      int    `json:"betOnIndex"`
	LineID          int64  `json:"lineID"`
	LineType        int    `json:"lineType"`
	LineDescription string `json:"lineDescription"`
	GameID          int64  `json:"gameID"`
	GameName        string `json:"gameName"`
}
