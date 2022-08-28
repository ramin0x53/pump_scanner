package main

import (
	"fmt"

	"github.com/ramin0x53/pump_scanner/api"
)

func main() {
	tcoin := api.Topcoins()
	e := api.GetAllKlines(tcoin, "1h", 1)
	fmt.Println(len(e))
}
