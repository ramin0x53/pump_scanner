package api

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/adshao/go-binance/v2"
)

type Klinef struct {
	OpenTime                 int64
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                int64
	QuoteAssetVolume         float64
	TradeNum                 int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

func klineWorker(symbols <-chan string, results chan<- []Klinef, coin chan<- string, tf string, limit int, wg *sync.WaitGroup) {
	for symbol := range symbols {
		results <- GetKlines(symbol, tf, limit)
		coin <- symbol
		wg.Done()
	}
}

func GetAllKlines(coins []string, tf string, limit int) map[string][]Klinef {
	allData := make(map[string][]Klinef)

	results := make(chan []Klinef, len(coins))
	coin := make(chan string, len(coins))
	jobs := make(chan string, len(coins))

	var wg sync.WaitGroup
	wg.Add(len(coins))
	for w := 1; w <= threadNum; w++ {
		go klineWorker(jobs, results, coin, tf, limit, &wg)
	}

	for _, i := range coins {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	for s := 0; s < len(coins); s++ {
		// r := <-results
		// c := <-coin
		// if c == "BTCUSDT" {
		// 	fmt.Println(r)
		// 	fmt.Println("-----------------------------------------------------------------------------")
		// }
		// allData[c] = r
		allData[<-coin] = <-results
	}
	close(coin)
	close(results)

	return allData
}

func GetKlines(symbol string, tf string, limit int) []Klinef {
	client := binance.NewClient("", "")
	klines, err := client.NewKlinesService().Symbol(symbol).Interval(tf).Limit(limit).Do(context.Background())
	if err != nil {
		log.Println(err)
	}
	var klinesf []Klinef
	for _, kline := range klines {
		klinesf = append(klinesf, Klinef{
			OpenTime:                 kline.OpenTime,
			Open:                     stringToFloat64(kline.Open),
			High:                     stringToFloat64(kline.High),
			Low:                      stringToFloat64(kline.Low),
			Close:                    stringToFloat64(kline.Close),
			Volume:                   stringToFloat64(kline.Volume),
			CloseTime:                kline.CloseTime,
			QuoteAssetVolume:         stringToFloat64(kline.QuoteAssetVolume),
			TradeNum:                 kline.TradeNum,
			TakerBuyBaseAssetVolume:  stringToFloat64(kline.TakerBuyBaseAssetVolume),
			TakerBuyQuoteAssetVolume: stringToFloat64(kline.TakerBuyQuoteAssetVolume),
		})
	}

	return klinesf
}

func Exist(symbol string) bool {
	client := binance.NewClient("", "")
	_, err := client.NewKlinesService().Symbol(symbol).Interval("1h").Limit(1).Do(context.Background())
	if err != nil {
		return false
	}
	return true
}

func stringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
