package main

import (
	"fmt"
	"tulip/client"
)

func main() {

	// getMarkets()
	// getTicker("btc-cop")
	getOrderBook("btc-cop")
}

// the apiUrl might change in the coming
// days with the introduction of Buda as the new company name
var apiURL = "https://www.surbtc.com/api/v2"

// getMarket is usted to get info about all markets
func getMarkets() {
	res := client.GetMarkets(apiURL)
	fmt.Print(res)
}

// getMarket is usted to get info about one (or every) market listed
func getTicker(ticker string) {
	res := client.GetTicker(apiURL, ticker)
	fmt.Print(res)
}

// getOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func getOrderBook(marketID string) {
	res := client.GetOrderBook(apiURL, marketID)
	fmt.Print(res)
}

// getTrades is used to get a list of most recent trades
func getTrades(marketID string) {
	res := client.GetOrderBook(apiURL, marketID)
	fmt.Print(res)
}
