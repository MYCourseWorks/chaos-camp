package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	createuser "github.com/homework/hw3/handlers/createUser"
	deleteuser "github.com/homework/hw3/handlers/deleteUser"
	getuser "github.com/homework/hw3/handlers/getUser"
	listusers "github.com/homework/hw3/handlers/listUsers"
	updateuser "github.com/homework/hw3/handlers/updateUser"
	"github.com/homework/hw3/infra"
	user "github.com/homework/hw3/models"
)

var db infra.UserDatabase

func init() {
	db = infra.NewUserDb()
	db[0] = user.NewAdmin(0, "admin", "secret", "admin.com", "Admin", "Admin")
}

func main() {
	router := mux.NewRouter()
	router.Handle("/users", listusers.Handle(db)).Methods(http.MethodGet)
	router.Handle("/users", createuser.Handle(db)).Methods(http.MethodPost)
	router.Handle("/user/{userID}", getuser.Handle(db)).Methods(http.MethodGet)
	router.Handle("/user/{userID}", updateuser.Handle(db)).Methods(http.MethodPut)
	router.Handle("/user/{userID}", deleteuser.Handle(db)).Methods(http.MethodDelete)

	port := ":8080"
	fmt.Println("\nListening on port " + port)
	http.Handle("/", router)
	http.ListenAndServe(port, router) // mux.Router now in play
}
