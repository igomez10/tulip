package client

import (
	"io/ioutil"
	"net/http"
)

type Client struct {
	ApiURL string
  ApiKey string
  ApiSecret string
  Nonce int32
	Authenticated bool
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() string {
	finalURL := c.ApiURL + "/markets"
	return execute("GET", finalURL)
}

// GetTicker Returns info about a specific market
func (c *Client) GetTicker(ticker string) string {
	finalURL := c.ApiURL + "/markets/" + ticker
	return execute("GET" , finalURL)
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) string {
	finalURL := c.ApiURL + "/markets/" + marketID + "/order_book"
	return execute("GET", finalURL)
}

//GetTrades returns a list of recent trades
func (c *Client) GetTrades(marketID string) string {
	finalURL := c.ApiURL  + "/markets/" + marketID + "/trades"
	return execute("GET", finalURL)
}


func execute(method string , completeURL string) string {
	if method == "GET" {
	resp, err := http.Get(completeURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error getting data, submitted api url was: " + completeURL + "with method" + method
	}
	return string(body)
	}
	return "TODO METHOD POST"
}
