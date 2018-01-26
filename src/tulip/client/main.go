package client

import (
	"io/ioutil"
	"net/http"
)

// GetMarkets Returns info about all markets
func GetMarkets(url string) string {
	resp, err := http.Get(url + "/markets")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error getting market prices, submitted url was: " + url
	}
	return string(body)
}

// GetTicker Returns info about a specific market
func GetTicker(url string, ticker string) string {
	resp, err := http.Get(url + "/markets/" + ticker)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error getting market price, submitted api url was: " + url + " and submitted ticker was:" + ticker + ". A working example is 'btc-cop' "
	}
	return string(body)
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func GetOrderBook(url string, marketID string) string {
	resp, err := http.Get(url + "/markets/" + marketID + "/order_book")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error getting market price, submitted api url was: " + url + " and submitted ticker was:" + marketID + ". A working example is 'btc-cop' "
	}
	return string(body)
}

//GetTrades returns a list of recent trades
func GetTrades(url string, marketID string) string {
	resp, err := http.Get(url + "/markets/" + marketID + "/order_book")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error getting market price, submitted api url was: " + url + " and submitted ticker was:" + marketID + ". A working example is 'btc-cop' "
	}
	return string(body)
}
