package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const expDurationInH = 2 // User expiration duration in hours

// Roles represents all valid user permission roles
type Roles int

const (
	// USER permission
	USER = 1 << iota
	// ADMIN permission
	ADMIN = 1 << iota
)

func (r Roles) String() string {
	switch r {
	case USER:
		return "USER"
	case USER & ADMIN:
		return "ADMIN"
	default:
		return "INVALID ROLE" // NOTE: should we panic here (probably not) ?
	}
}

// User comment
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Enabled   bool   `json:"-"`
	Expired   bool   `json:"-"`
	Roles     Roles  `json:"roles"`
}

func (u *User) String() string {
	return fmt.Sprintf(
		"ID: %d\nUsername: %s\nEmail: %s\nFirstName: %s\nLastName: %s\nEnabled: %t\nExpired: %t\n",
		u.ID,
		u.Username,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Enabled,
		u.Expired,
	)
}

// NewAdmin creates a new admin user
func NewAdmin(id int, username string, password string, email string,
	firstName string, lastName string) *User {

	user := &User{
		ID:        id,
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Enabled:   true,
		Expired:   false,
		Roles:     USER | ADMIN,
	}

	user.SetPassword(password)
	return user
}

// NewUser comment
func NewUser(id int, username string, password string, email string,
	firstName string, lastName string) *User {

	user := &User{
		ID:        id,
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Enabled:   true,
		Expired:   false,
		Roles:     USER,
	}

	user.SetPassword(password)
	return user
}

// SetPassword hashes the users password and saves it
func (u *User) SetPassword(password string) error {
	if password != "" {
		passAsBytes := []byte(password)
		hash, err := bcrypt.GenerateFromPassword(passAsBytes, bcrypt.MinCost)
		if err != nil {
			return err
		}

		u.Password = string(hash)
	} else {
		return errors.New("Provided password is invalid")
	}

	return nil
}
