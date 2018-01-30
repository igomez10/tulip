package main

import (
	"fmt"
	"os"
	"tulip"
)

func main() {

	APIKey := os.Getenv("BUDAKEY")       // you can modify this variable and hardocde your own apikey
	APISecret := os.Getenv("BUDASECRET") // you can modify this variable and hardcode your own apisecret

	buda := tulip.CreateClient(APIKey, APISecret)
	results, err := (buda.GetTicker("btc-clp"))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(results.Market.ID)                 // "BTC-CLP"
		fmt.Println(results.Market.Name)               // "btc-clp"
		fmt.Println(results.Market.BaseCurrency)       // "BTC"
		fmt.Println(results.Market.QuoteCurrency)      // "CLP"
		fmt.Println(results.Market.MinimumOrderAmount) // ["0.0001 BTC"]
	}
	// Uncomment a line to try it

	// to create a new order that will never be fullfilled fmt.Println(myClient.PostOrder("btc-cop", "bid", "limit", float64(0.00001), float64(0.0001)))
	// fmt.Println(myClient.PostOrder("btc-cop", "bid", "limit", float64(0.00001), float64(0.0001)))
	// fmt.Println(myClient.GetOrders("btc-cop", 300, 1, "pending", float64(0)))
	// fmt.Println(myClient.GetWithdrawHistory("cop"))
	// fmt.Println(myClient.GetOrder("1234567"))
	// fmt.Println(myClient.CancelOrder("1234567"))
	// fmt.Printf(myClient.GetMarkets())
	// fmt.Printf(myClient.GetTicker("btc-cop"))
	// fmt.Printf(myClient.GetOrderBook("btc-cop"))
	// fmt.Printf(myClient.GetTrades("btc-cop"))

}
