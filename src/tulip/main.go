package main

import (
	"fmt"
	"os"
	"tulip/client"
)

func main() {
	// var authenticated bool = false
	const APIURL string = "https://www.surbtc.com/api/v2"
	APIKey := os.Getenv("BUDAKEY")
	APISecret := os.Getenv("BUDASECRET")
	var myClient client.Client
	myClient.APIURL = APIURL
	if APIKey != "" && APISecret != "" {

		myClient.APIURL = APIURL
		myClient.APIKey = APIKey
		myClient.APISecret = APISecret
		myClient.Authenticated = true
		fmt.Println(myClient.GetOrders("btc-cop", 300, 1, "received", float64(0.00001)))
	}
	// fmt.Printf(myClient.GetMarkets())
	// fmt.Printf(myClient.GetTicker("btc-cop"))
	// fmt.Printf(myClient.GetOrderBook("btc-cop"))
	// fmt.Printf(myClient.GetTrades("btc-cop"))

}
