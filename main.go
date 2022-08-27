package main

import (
	"fmt"

	"github.com/ramin0x53/pump_scanner/api"
)

func main() {
	tcoin := api.Topcoins()
	println(len(tcoin))
	e := api.GetAllKlines(tcoin, "1h", 10)
	fmt.Println(e["ETHUSDT"])
	fmt.Println(len(e))
}
