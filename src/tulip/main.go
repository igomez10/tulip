package main

import (
	"fmt"
	"tulip/client"
	"os"
)

func main() {
	// var authenticated bool = false
	const apiURL string = "https://www.surbtc.com/api/v2"
	var apiKey string = os.Getenv("budakey")
	var apiSecret string = os.Getenv("budasecret")
	var myClient client.Client
	myClient.ApiURL = apiURL
	if apiKey != "" && apiSecret != "" {

			myClient.ApiURL = apiURL
			myClient.ApiKey = apiKey
			myClient.ApiSecret = apiSecret
			myClient.Authenticated = true

		}
		fmt.Printf(myClient.GetMarkets())
		fmt.Printf(myClient.GetTicker("btc-cop"))
		fmt.Printf(myClient.GetOrderBook("btc-cop"))
		fmt.Printf(myClient.GetTrades("btc-cop"))

	}
