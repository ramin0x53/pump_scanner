package main

import (
	"fmt"
	"time"

	"github.com/ramin0x53/pump_scanner/api"
)

/*          KLINE_INTERVAL_1MINUTE = '1m'
            KLINE_INTERVAL_3MINUTE = '3m'
            KLINE_INTERVAL_5MINUTE = '5m'
            KLINE_INTERVAL_15MINUTE = '15m'
            KLINE_INTERVAL_30MINUTE = '30m'
            KLINE_INTERVAL_1HOUR = '1h'
            KLINE_INTERVAL_2HOUR = '2h'
            KLINE_INTERVAL_4HOUR = '4h'
            KLINE_INTERVAL_6HOUR = '6h'
            KLINE_INTERVAL_8HOUR = '8h'
            KLINE_INTERVAL_12HOUR = '12h'
            KLINE_INTERVAL_1DAY = '1d'
*/

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

var data = make(map[string][]api.Klinef)
var tcoin []string
var tf = "5m"
var lenght = 10
var baseN float64 = 2

func main() {
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
	fmt.Println("-------------------------")
}

func dump(symbol string) {
	fmt.Println(symbol + " dumped")
	fmt.Println("-------------------------")
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
