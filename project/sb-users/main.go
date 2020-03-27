package main

import (
	"fmt"

	"github.com/MartinNikolovMarinov/sb-users/rest"
)

const port = ":8080"

func main() {
	fmt.Println("Sports Betting User App starting on port:", port)
	a := rest.App{}
	a.Init()
	a.Run(port)
}
