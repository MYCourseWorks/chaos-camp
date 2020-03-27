package main

import (
	"fmt"

	"github.com/MartinNikolovMarinov/sb-feed/rest"
)

const port = ":8081"

func main() {
	fmt.Println("Sports Betting User App starting on port:", port)
	a := rest.App{}
	a.Init()
	// a.Scrape()
	a.Run(port)
}
