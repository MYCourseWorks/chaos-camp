package main

import (
	"fmt"

	"github.com/MartinNikolovMarinov/sb-bets/rest"
)

const port = ":8082"

func main() {
	fmt.Println("Sports Betting Bet App starting on port:", port)
	a := rest.App{}
	a.Init()
	a.Run(port)
}
