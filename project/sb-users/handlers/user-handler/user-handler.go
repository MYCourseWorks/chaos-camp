package userhandler

import (
	"net/http"
	"strconv"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-users/models"
	siterepo "github.com/MartinNikolovMarinov/sb-users/repositories/site-repo"
	userrepo "github.com/MartinNikolovMarinov/sb-users/repositories/user-repo"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

// LoginCMD comment
type LoginCMD struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=5,max=255"`
}

// Login comment
func Login(ur userrepo.UserRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c LoginCMD
		var err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return err
		}

		var u *models.User
		u, err = ur.Find(c.Name)
		if err != nil {
			return infra.NewHTTPError(err, 403, "Ivalid Username or password")
		}

		err = ur.MatchPassword(c.Password, u.Password)
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return infra.NewHTTPError(err, 403, "Ivalid Username or password")
		}
		if err != nil {
			return err
		}

		token, err := infra.GenerateToken(u.ID, u.Name)
		if err != nil {
			return err
		}

		var resp = map[string]interface{}{
			"token": token,
			"id":    u.ID,
			"roles": u.Roles,
		}
		err = infra.WriteResponseJSON(w, resp, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// RegisterCMD comment
type RegisterCMD struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=5,max=255"`
	SiteID   int64  `json:"siteID" validate:"required,numeric,gte=0"`
}

// Register comment
func Register(ur userrepo.UserRepo, sr siterepo.SiteRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c RegisterCMD
		var err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return err
		}

		site, err := sr.FindByID(c.SiteID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}
		if site == nil {
			return infra.NewHTTPError(err, 400, "Site does not exist")
		}

		userInDb, err := ur.Find(c.Name)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}
		if userInDb != nil {
			return infra.NewHTTPError(err, 400, "User already exists")
		}

		var u *models.User
		u, err = ur.Create(c.Name, c.Password, models.SiteUser, c.SiteID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = infra.WriteResponseJSON(w, u, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// CreateCMD comment
type CreateCMD struct {
	Name     string           `json:"name" validate:"required,min=3,max=255"`
	Password string           `json:"password" validate:"required,min=5,max=255"`
	Roles    models.RoleFlags `json:"roles" validate:"required"`
	SiteID   int64            `json:"siteID" validate:"required,numeric,gte=0"`
}

// Create comment
func Create(ur userrepo.UserRepo, sr siterepo.SiteRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c CreateCMD
		var err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return err
		}

		site, err := sr.FindByID(c.SiteID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}
		if site == nil {
			return infra.NewHTTPError(err, 400, "Site does not exist")
		}

		userInDb, err := ur.Find(c.Name)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}
		if userInDb != nil {
			return infra.NewHTTPError(err, 400, "User already exists")
		}

		var u *models.User
		u, err = ur.Create(c.Name, c.Password, c.Roles, c.SiteID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = infra.WriteResponseJSON(w, u, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// All comment
func All(ur userrepo.UserRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		var siteID int64
		siteID, err = strconv.ParseInt(chi.URLParam(r, "siteID"), 10, 64)
		if err != nil {
			return infra.NewHTTPError(err, 400, "siteID paramter was invalid")
		}

		var users []models.User
		users, err = ur.AllPerSite(siteID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		err = infra.WriteResponseJSON(w, users, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// Delete comment
func Delete(ur userrepo.UserRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var err error
		var userID int64
		userID, err = strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			return infra.NewHTTPError(err, 400, "siteID paramter was invalid")
		}

		err = ur.Delete(userID)
		if err != nil {
			return infra.NewHTTPError(err, 400, err.Error())
		}

		resp := map[string]string{"detail": "user deleted"}
		err = infra.WriteResponseJSON(w, resp, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}

// UpdateCMD comment
type UpdateCMD struct {
	UserID   int64  `json:"userID" validate:"required,min=3,max=255"`
	Username string `json:"username" validate:"omitempty,min=3,max=255"`
	Password string `json:"password" validate:"omitempty,min=5,max=255"`
}

// Update comment
func Update(ur userrepo.UserRepo, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		var c UpdateCMD
		var err = infra.DecodeAndValidate(&c, r.Body, v)
		if err != nil {
			return err
		}

		if c.Username != "" {
			err = ur.UpdateUsername(c.UserID, c.Username)
			if err != nil {
				return infra.NewHTTPError(err, 400, err.Error())
			}
		}

		if c.Password != "" {
			err = ur.UpdatePassword(c.UserID, c.Password)
			if err != nil {
				return infra.NewHTTPError(err, 400, err.Error())
			}
		}

		resp := map[string]bool{"ok": true}
		err = infra.WriteResponseJSON(w, resp, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
