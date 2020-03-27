package infra

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator"
)

// DecodeAndValidate comment
func DecodeAndValidate(addr interface{}, reader io.Reader, v *validator.Validate) error {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&addr)
	if err != nil {
		return NewHTTPError(err, 400, "Request body is not valid")
	}

	// VALIDATION :
	err = v.Struct(addr)
	if err != nil {
		return NewHTTPError(err, 400, err.Error())
	}

	return nil
}

// WriteResponseJSON comment
func WriteResponseJSON(w http.ResponseWriter, obj interface{}, status int) error {
	var JSON, err = json.Marshal(obj)
	if err != nil {
		return NewHTTPError(err, 500, "JSON marshall error")
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(JSON)
	return nil
}
