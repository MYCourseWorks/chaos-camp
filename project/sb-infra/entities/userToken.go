package entities

import "github.com/dgrijalva/jwt-go"

// UserToken comment
type UserToken struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	jwt.StandardClaims
}
