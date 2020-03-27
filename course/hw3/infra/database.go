package infra

import (
	user "github.com/homework/hw3/models"
)

// UserDatabase commnet
type UserDatabase map[int]*user.User

// NewUserDb comment
func NewUserDb() UserDatabase {
	var db UserDatabase = make(map[int]*user.User)
	return db
}

// FindUserByID comment
func (d UserDatabase) FindUserByID(id int) *user.User {
	if u, ok := d[id]; ok {
		return u
	}

	return nil
}

var nextID int = 0

// NextID comment
func (d UserDatabase) NextID() int {
	nextID++
	return nextID
}
