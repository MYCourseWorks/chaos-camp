package betshandler

import (
	"net/http"
	"strconv"

	"github.com/MartinNikolovMarinov/sb-bets/models"
	betsrepo "github.com/MartinNikolovMarinov/sb-bets/repositories/bets-repo"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

// All comment
func All(b betsrepo.BetsRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var (
			err       error
			userID    int64
			bets      []models.Bet
			userIDStr string
		)

		userIDStr = chi.URLParam(r, "userID")
		if userIDStr != "" {
			userID, err = strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				return infra.NewHTTPError(err, 400, "userID paramter was invalid")
			}

			bets, err = b.AllByUserID(userID)
			if err != nil {
				return err
			}
		} else {
			bets, err = b.All()
			if err != nil {
				return err
			}
		}

		err = infra.WriteResponseJSON(w, bets, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// PlaceBetCMD comment
type PlaceBetCMD struct {
	UserID     int64  `json:"userID" validate:"required,gte=1"`
	LineID     int64  `json:"lineID" validate:"required,gte=1"`
	Value      string `json:"value" validate:"required"`
	BetOnIndex int64  `json:"betOnIndex" validate:"gte=0"`
}

// Place comment
func Place(b betsrepo.BetsRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var (
			c     PlaceBetCMD
			err   error
			betID int64
		)

		err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		betID, err = b.PlaceBet(c.UserID, c.LineID, c.BetOnIndex, c.Value)
		if err != nil {
			return err
		}

		err = infra.WriteResponseJSON(w, map[string]int64{"betID": betID}, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// CancelCMD comment
type CancelCMD struct {
	UserID int64 `json:"userID" validate:"required,gte=1"`
	BetID  int64 `json:"betID" validate:"required,gte=1"`
}

// Cancel comment
func Cancel(b betsrepo.BetsRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var (
			c   CancelCMD
			err error
		)

		err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = b.CancelBet(c.UserID, c.BetID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = infra.WriteResponseJSON(w, map[string]bool{"ok": true}, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// PayoutCMD comment
type PayoutCMD struct {
	BetID int64 `json:"betID" validate:"required,gte=1"`
}

// Payout comment
func Payout(b betsrepo.BetsRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var (
			c   PayoutCMD
			err error
		)

		err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = b.PayoutBet(c.BetID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = infra.WriteResponseJSON(w, map[string]bool{"ok": true}, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
