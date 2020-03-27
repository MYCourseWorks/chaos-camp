package betsrepo

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MartinNikolovMarinov/sb-bets/models"
	repo "github.com/MartinNikolovMarinov/sb-bets/repositories"
)

type betsRepoMysql struct {
	db *sql.DB
}

// NewMysqlRepo is a  constructor
func NewMysqlRepo(connectionString string) BetsRepo {
	repo := &betsRepoMysql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func (br *betsRepoMysql) PlaceBet(userID, lineID, betOnIndex int64, value string) (int64, error) {
	const placeBetQuery = `
		insert into Bet (LineID, UserID, Value, BetOnIndex) Values (?, ?, ?, ?)
	`

	var (
		result sql.Result
		err    error
		id     int64
	)

	result, err = br.db.Exec(placeBetQuery, lineID, userID, value, betOnIndex)
	if err != nil {
		return -1, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (br *betsRepoMysql) CancelBet(userID, betID int64) error {
	const cancelBetQuery = `
		update Bet
		set IsDeleted = 1
		where UserID = ? and ID = ?
	`

	result, err := br.db.Exec(cancelBetQuery, userID, betID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("Bet was not found")
	}

	return nil
}

func (br *betsRepoMysql) PayoutBet(betID int64) error {
	const payoutBetQuery = `
		update Bet
		set IsPayed = 1
		where ID = ?
	`

	result, err := br.db.Exec(payoutBetQuery, betID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("Bet was not found")
	}

	return nil
}

func (br *betsRepoMysql) AllByUserID(userID int64) ([]models.Bet, error) {
	const allByUserID = `
	select
		b.ID as BetID,
		b.UserID as UserID,
		u.Name as UserName,
		b.Value as BetValue,
		b.IsPayed as IsPayed,
		b.BetOnIndex as BetOnIndex,
		l.ID as LineID,
		l.LineType as LineType,
		l.Description as LineDescription,
		g.ID as GameID,
		g.Name as GameName
	from Bet as b
	join Line as l on l.ID = b.LineID and l.IsDeleted = 0
	join GameEvent as e on e.ID = l.EventID and e.IsDeleted = 0
	join Game as g on g.ID = e.RelatedGameID and g.IsDeleted = 0
	join User as u on u.ID = b.UserID and u.IsDeleted = 0
	where
		b.IsDeleted = 0 and
		b.UserID = ?
	`

	rows, err := br.db.Query(allByUserID, userID)
	if err != nil {
		return nil, err
	}

	return parseBets(rows)
}

func (br *betsRepoMysql) All() ([]models.Bet, error) {
	const allByUserID = `
	select
		b.ID as BetID,
		b.UserID as UserID,
		u.Name as UserName,
		b.Value as BetValue,
		b.IsPayed as IsPayed,
		b.BetOnIndex as BetOnIndex,
		l.ID as LineID,
		l.LineType as LineType,
		l.Description as LineDescription,
		g.ID as GameID,
		g.Name as GameName
	from Bet as b
	join Line as l on l.ID = b.LineID and l.IsDeleted = 0
	join GameEvent as e on e.ID = l.EventID and e.IsDeleted = 0
	join Game as g on g.ID = e.RelatedGameID and g.IsDeleted = 0
	join User as u on u.ID = b.UserID and u.IsDeleted = 0
	where
		b.IsDeleted = 0
	`

	rows, err := br.db.Query(allByUserID)
	if err != nil {
		return nil, err
	}

	return parseBets(rows)
}

func parseBets(rows *sql.Rows) ([]models.Bet, error) {
	var (
		betID           int64
		userID          int64
		userName        string
		betValueStr     string
		isPayedBit      []byte
		isPlayed        bool
		betOnIndex      int
		lineID          int64
		lineType        int
		lineDescription string
		gameID          int64
		gameName        string
		err             error
		ret             []models.Bet = make([]models.Bet, 0)
	)

	for rows.Next() {
		err = repo.ScanRow(
			rows,
			&betID, &userID, &userName, &betValueStr, &isPayedBit,
			&betOnIndex, &lineID, &lineType, &lineDescription,
			&gameID, &gameName,
		)
		if err != nil {
			return nil, err
		}

		isPlayed = isPayedBit[0] == byte(1)

		bet := models.Bet{
			ID:              betID,
			UserID:          userID,
			UserName:        userName,
			Value:           betValueStr,
			IsPayed:         isPlayed,
			BetOnIndex:      betOnIndex,
			LineID:          lineID,
			LineType:        lineType,
			LineDescription: lineDescription,
			GameID:          gameID,
			GameName:        gameName,
		}
		ret = append(ret, bet)
	}

	return ret, nil
}
