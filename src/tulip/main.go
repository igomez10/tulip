package main

import (
	"fmt"
	"os"
	"tulip/client"
)

func main() {
	const APIURL string = "https://www.surbtc.com/api/v2"
	APIKey := os.Getenv("BUDAKEY")       // you can modify this variable and hardocde your own apikey
	APISecret := os.Getenv("BUDASECRET") // you can modify this variable and hardcode your own apisecret
	myClient := tulip.CreateClient(APIKey, APISecret)
	fmt.Println(myClient.GetMarkets())
	// myClient.APIURL = APIURL
	// if APIKey != "" && APISecret != "" {
	//
	// 	myClient.APIURL = APIURL
	// 	myClient.APIKey = APIKey
	// 	myClient.APISecret = APISecret
	// 	myClient.Authenticated = true
	//
	// }
	// EXAMPLES:
	// to create a new order that will never be fullfilled fmt.Println(myClient.PostOrder("btc-cop", "bid", "limit", float64(0.00001), float64(0.0001)))
	// fmt.Println(myClient.PostOrder("btc-cop", "bid", "limit", float64(0.00001), float64(0.0001)))
	// fmt.Println(myClient.GetOrders("btc-cop", 300, 1, "pending", float64(0)))
	// fmt.Println(myClient.GetWithdrawHistory("cop"))
	// fmt.Println(myClient.GetOrder("5205228"))
	// fmt.Println(myClient.CancelOrder("5205033"))
	// fmt.Printf(myClient.GetMarkets())
	// fmt.Printf(myClient.GetTicker("btc-cop"))
	// fmt.Printf(myClient.GetOrderBook("btc-cop"))
	// fmt.Printf(myClient.GetTrades("btc-cop"))

}
