package infra

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/dgrijalva/jwt-go"
)

const secret = "secret"

// GenerateToken comment
func GenerateToken(userID int64, userName string) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	claims := &entities.UserToken{
		UserID: string(userID),
		Name:   userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JwtVerifyMiddleware comment
func JwtVerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerSchema = "Bearer "
		var header = r.Header.Get("Authorization")

		var token string
		if header == "" || len(header) <= len(bearerSchema) {
			errResp := NewHTTPError(
				errors.New("Failed to authenticate"),
				http.StatusForbidden,
				"Failed to authenticate")
			WriteResponseJSON(w, errResp, http.StatusUnauthorized)
			return
		}

		token = header[len(bearerSchema):]
		token = strings.TrimSpace(token)

		if token == "" {
			// Token is missing, returns with error code 403 Unauthorized
			errResp := NewHTTPError(
				errors.New("Missing auth token"),
				http.StatusForbidden,
				"Missing auth token")
			WriteResponseJSON(w, errResp, http.StatusUnauthorized)
			return
		}

		claims := &entities.UserToken{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			errResp := NewHTTPError(
				errors.New("Missing auth token"),
				http.StatusForbidden,
				err.Error())
			WriteResponseJSON(w, errResp, http.StatusUnauthorized)
			return
		}

		// FIXME: How do we get roles from this ?
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
