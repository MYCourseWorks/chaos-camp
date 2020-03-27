package main

import (
	"fmt"

	"github.com/homework/hw3/client/util"
)

// GET /user/{userID}
func test1GetUser() {
	fmt.Println(util.SendGet("http://localhost:8080/user/0"))

	// NO SUCH USER :
	fmt.Println(util.SendGet("http://localhost:8080/user/500"))

	// BAD Request:
	fmt.Println(util.SendGet("http://localhost:8080/user/not_a_number"))

	// Another BAD Request:
	fmt.Println(util.SendGet("http://localhost:8080/user/12.31"))

	// NOT FOUND 404:
	fmt.Println(util.SendGet("http://localhost:8080/user/1/should/be/invalid"))
}

// POST /user
func test2CreateUser() {
	var resp, location string

	// Valid creates a user
	resp, location = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "test",
		"email":      "test.com",
		"first-name": "Test",
		"last-name":  "asd",
		"password":   "123",
	})
	fmt.Println(resp)
	fmt.Println(util.SendGet(fmt.Sprintf("http://localhost:8080/%s", location)))

	// STATUS 400 NO EMAIL
	resp, _ = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "test",
		"first-name": "Test",
		"last-name":  "asd",
		"password":   "123",
	})
	fmt.Println(resp)

	// STATUS 400 NO USERNAME, NO EMAIL, NO PASSWORD
	resp, _ = util.SendPost("http://localhost:8080/users", map[string]string{
		"first-name": "Test",
		"last-name":  "asd",
		"password":   "",
	})
	fmt.Println(resp)

	// STATUS 400 nothing is set
	resp, _ = util.SendPost("http://localhost:8080/users", map[string]string{})
	fmt.Println(resp)

	// STATUS 400 request body is not valid
	resp, _ = util.SendPost("http://localhost:8080/users", nil)
	fmt.Println(resp)

	util.SendDelete(fmt.Sprintf("http://localhost:8080/%s", location))
	util.SendDelete(fmt.Sprintf("http://localhost:8080/%s", location))
}

// POST, PUT, DELETE /user
func test3CreateUpdateDelete() {
	var resp, location string

	// Create a valid user
	resp, location = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "test",
		"email":      "test.com",
		"first-name": "Test",
		"last-name":  "asd",
		"password":   "123",
	})
	fmt.Println(resp)
	locationURL := fmt.Sprintf("http://localhost:8080/%s", location)
	fmt.Println(util.SendGet(locationURL))

	// Update user 1
	resp = util.SendPut(locationURL, map[string]string{
		"username": "chaged",
		"email":    "cahanged.com",
	})
	fmt.Println(resp)

	// Delete that user
	resp = util.SendDelete(locationURL)
	fmt.Println(resp)
	// Try to find that user after delete :
	fmt.Println(util.SendGet(locationURL))
}

// GET /users
func test4ListingAndSorting() {
	var resp, l1, l2, l3 string

	// CREATE 3 users
	resp, l1 = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "a",
		"email":      "a.com",
		"first-name": "A",
		"last-name":  "A",
		"password":   "123",
	})
	fmt.Println(resp)

	resp, l2 = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "c",
		"email":      "c.com",
		"first-name": "C",
		"last-name":  "C",
		"password":   "123",
	})
	fmt.Println(resp)

	resp, l3 = util.SendPost("http://localhost:8080/users", map[string]string{
		"username":   "B",
		"email":      "b.com",
		"first-name": "B",
		"last-name":  "b",
		"password":   "123",
	})
	fmt.Println(resp)

	fmt.Println(util.SendGet("http://localhost:8080/users"))
	// fmt.Println(util.SendGet("http://localhost:8080/users?sortBy=username"))
	// fmt.Println(util.SendGet("http://localhost:8080/users?sortBy=role"))

	util.SendDelete(fmt.Sprintf("http://localhost:8080/%s", l1))
	util.SendDelete(fmt.Sprintf("http://localhost:8080/%s", l2))
	// PUT expired to true does the same thing as delete
	util.SendPut(fmt.Sprintf("http://localhost:8080/%s", l3), map[string]bool{
		"expired": true,
	})
}

func main() {
	test2CreateUser()
}
