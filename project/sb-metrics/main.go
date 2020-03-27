package main

import (
	"fmt"

	"github.com/MartinNikolovMarinov/sb-metrics/rest"
)

const port = ":8083"

func main() {
	fmt.Println("Sports Betting Metrics App starting on port:", port)
	a := rest.App{}
	a.Init()
	a.Run(port)
}
