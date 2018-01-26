package main

import (
	"fmt"
	"os"
	"tulip/client"
)

func main() {
	// var authenticated bool = false
	const apiURL string = "https://www.surbtc.com/api/v2"
	var apiKey string = os.Getenv("BUDAKEY")
	var apiSecret string = os.Getenv("BUDASECRET")
	var myClient client.Client
	myClient.ApiURL = apiURL
	if apiKey != "" && apiSecret != "" {

		myClient.ApiURL = apiURL
		myClient.ApiKey = apiKey
		myClient.ApiSecret = apiSecret
		myClient.Authenticated = true
		fmt.Println(myClient.GetBalances())
	}
	// fmt.Printf(myClient.GetMarkets())
	// fmt.Printf(myClient.GetTicker("btc-cop"))
	// fmt.Printf(myClient.GetOrderBook("btc-cop"))
	// fmt.Printf(myClient.GetTrades("btc-cop"))

}
