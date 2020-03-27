package infra

import (
	"log"
	"net/http"
)

// RootHandler commnet
type RootHandler func(w http.ResponseWriter, r *http.Request) error

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r) // Call handler function
	if err == nil {
		return
	}

	// NOTE: This is where our error handling logic starts :

	log.Printf("An error accured: %v", err) // Log the error.

	clientError, ok := err.(ClientError)
	if !ok {
		// If the error is not ClientError, assume that it is ServerError.
		w.WriteHeader(500)
		return
	}

	body, err := clientError.ResponseBody()
	if err != nil {
		log.Printf("An error accured: %v", err)
		w.WriteHeader(500)
		return
	}

	status, headers := clientError.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
	w.Write(body)
}
