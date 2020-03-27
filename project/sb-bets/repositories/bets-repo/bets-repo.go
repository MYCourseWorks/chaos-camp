package betsrepo

import "github.com/MartinNikolovMarinov/sb-bets/models"

// BetsRepo comment
type BetsRepo interface {
	PlaceBet(userID, lineID, betOnIndex int64, value string) (int64, error)
	CancelBet(userID, betID int64) error
	PayoutBet(betID int64) error
	AllByUserID(userID int64) ([]models.Bet, error)
	All() ([]models.Bet, error)
}
