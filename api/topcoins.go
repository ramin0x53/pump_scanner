package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Coin struct {
	Id              string  `json:"id"`
	Symbol          string  `json:"symbol"`
	Market_cap      float64 `json:"market_cap"`
	Market_cap_rank int     `json:"market_cap_rank"`
	Total_volume    float64 `json:"total_volume"`
}

var threadNum = 10

func Topcoins() []string {
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=volume_desc&per_page=250&page=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	var coins []Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		log.Println(err)
	}
	var topcoins []string
	for _, coin := range coins {
		if !strings.Contains(coin.Symbol, "usd") {
			topcoins = append(topcoins, strings.ToUpper(coin.Symbol)+"USDT")
		}
	}

	results := make(chan string, len(topcoins))
	jobs := make(chan string, len(topcoins))

	for w := 1; w <= threadNum; w++ {
		go worker(jobs, results)
	}

	for _, i := range topcoins {
		jobs <- i
	}
	close(jobs)

	var binanceCoins []string
	for s := 0; s < len(topcoins); s++ {
		re := <-results
		if re != "false" {
			binanceCoins = append(binanceCoins, re)
		}
	}
	close(results)

	return binanceCoins
}

func worker(txt <-chan string, results chan<- string) {
	for j := range txt {
		if Exist(j) {
			results <- j
			continue
		} else {
			results <- "false"
		}
	}
}
