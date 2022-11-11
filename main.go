package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"

	"github.com/ramin0x53/pump_scanner/api"
)

var tframes = map[string]int64{
	"1m":  60,
	"3m":  180,
	"5m":  300,
	"15m": 900,
	"30m": 1800,
	"1h":  3600,
	"2h":  7200,
	"4h":  14400,
	"6h":  21600,
	"8h":  28800,
	"12h": 43200,
	"1d":  86400,
}

var divide = "-------------------------"
var data = make(map[string][]api.Klinef)
var tcoin []string
var tf string
var lenght int
var baseN float64

func options() {
	flag.StringVar(&tf, "tf", "5m", "Timeframe")
	flag.Float64Var(&baseN, "b", 2, "Base number for comparing volatility")
	flag.IntVar(&lenght, "l", 10, "Count of candles to compare")
	flag.IntVar(&api.ThreadNum, "t", 10, "Threads number")
	flag.Parse()
	fmt.Println("Timeframe: " + tf)
	fmt.Printf("Base number: %f\n", baseN)
	fmt.Printf("Lenght: %d\n", lenght)
	fmt.Printf("Threads number: %d\n", api.ThreadNum)
	fmt.Println(divide)
}

func main() {
	options()
	Run()
}

func Run() {
	tcoin = api.Topcoins()
	for {
		if len(data) == 0 {
			updateData()
		}
		timeNow := time.Now().Unix()
		if timeNow == (data["BTCUSDT"][len(data["BTCUSDT"])-1].OpenTime/1000)+(tframes[tf]*2)+1 {
			go updateData()
		}
		time.Sleep(1 * time.Second)
	}
}

func updateData() {
	if len(data) == 0 {
		data = api.GetAllKlines(tcoin, tf, lenght+1)
		for key := range data {
			data[key] = data[key][:len(data[key])-1]
		}
	} else {
		temp := api.GetAllKlines(tcoin, tf, 2)
		for key := range temp {
			temp[key] = temp[key][:len(temp[key])-1]
		}

		for key := range data {
			data[key] = append(data[key], temp[key]...)
		}

		for key := range data {
			data[key] = data[key][1:len(data[key])]
		}
	}
	checkOpentimes()
	findPumpDump()
}

func checkOpentimes() {
	if (data["BTCUSDT"][len(data["BTCUSDT"])-1].OpenTime / 1000) != (data["BTCUSDT"][len(data["BTCUSDT"])-2].OpenTime/1000)+tframes[tf] {
		fmt.Println("Uncoordinated opentimes !!!")
		fmt.Println("Restarting...")
		data = make(map[string][]api.Klinef)
	}
}

func highLowPerc(high float64, low float64, candleColor string) float64 {
	if candleColor == "green" {
		return ((high - low) * 100) / low
	} else if candleColor == "red" {
		return ((high - low) * 100) / high
	} else {
		panic("Error: Candle color not recognized")
	}
}

func openClosePerc(open float64, close float64) float64 {
	if open <= close {
		return ((close - open) * 100) / open
	} else {
		return (-(open - close) * 100) / open
	}
}

func averagePerc(symbol string) float64 {
	var r float64
	for i := 0; i < len(data[symbol]); i++ {
		if data[symbol][i].Open <= data[symbol][i].Close {
			r = r + highLowPerc(data[symbol][i].High, data[symbol][i].Low, "green")
		} else {
			r = r + highLowPerc(data[symbol][i].High, data[symbol][i].Low, "red")
		}
	}
	return r / float64(len(data[symbol])-1)
}

func calVol(symbol string) float64 {
	return openClosePerc(data[symbol][len(data[symbol])-1].Open, data[symbol][len(data[symbol])-1].Close) / averagePerc(symbol)
}

func pump(symbol string) {
	fmt.Println(symbol + " pumped")
	// alert(symbol, symbol+" pumped!")
	fmt.Println(divide)
}

func dump(symbol string) {
	fmt.Println(symbol + " dumped")
	// alert(symbol, symbol+" dumped!")
	fmt.Println(divide)
}

func findPumpDump() {
	for key := range data {
		volatility := calVol(key)
		if volatility >= baseN {
			pump(key)
		} else if volatility <= -baseN {
			dump(key)
		}
	}
}

func alert(symbol string, msg string) {
	err := beeep.Alert(symbol, msg, "")
	if err != nil {
		panic(err)
	}
}
