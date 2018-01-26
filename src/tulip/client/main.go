package client

import (
	"io/ioutil"
	"net/http"
	"time"
	// "fmt"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strconv"
)

type Client struct {
	ApiURL        string
	ApiKey        string
	ApiSecret     string
	Nonce         int32
	Authenticated bool
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() string {
	finalURL := c.ApiURL + "/markets"
	return execute("GET", finalURL, false)
}

// GetTicker Returns info about a specific market
func (c *Client) GetTicker(ticker string) string {
	finalURL := c.ApiURL + "/markets/" + ticker
	return execute("GET", finalURL, false)
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) string {
	finalURL := c.ApiURL + "/markets/" + marketID + "/order_book"
	return execute("GET", finalURL, false)
}

//GetTrades returns a list of recent trades
func (c *Client) GetTrades(marketID string) string {
	finalURL := c.ApiURL + "/markets/" + marketID + "/trades"
	return execute("GET", finalURL, false)
}

func execute(method string, completeURL string, requiresAuth bool) string {
	if method == "GET" && requiresAuth == false {

		req, err := http.NewRequest("GET", completeURL, nil)
		if err != nil {
			log.Fatal(("Error creating new request, submitted url was: " +
				completeURL + "with method" + method +
				". Error: "), err)
		}

		httpClient := &http.Client{}

		res, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(("Error executing new request, submitted url was: " +
				completeURL + "with method" + method +
				". Error: "), err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(("Error reading response, submitted url was: " +
				completeURL + "with method" + method +
				". Error: "), err)
		}
		return string(body)

		// TODO ADD POST METHODS
	}
	// if method == "GET" && requiresAuth == true {
	// 	Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	// 	shortURL := comple
	// 	stringMessage := method + " " + strings.Replace(completeURL, old, new, n)+" " + Nonce
	//
	// 	// if c.Authenticated  == true{
	// 	//
	// 	// } else {
	// 	// 	log.Fatal("this method requires Authentication. BUDAKEY(api public key), BUDASECRET (api secret key) env variables were not found")
	// 	// 	}
	// }

	return "INCORRECT METHOD, ONLY GET IS AVAILABLE AT THE TIME"
}

// HERE STARTS THE PRIVATE CALLS

func (c *Client) signMessage(method string, query string, Nonce string) string {
	stringMessage := method + " /api/v2/" + query + " " + Nonce
	key := []byte(c.ApiSecret)
	h := hmac.New(sha512.New384, key)
	h.Write([]byte(stringMessage))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

func (c *Client) GetBalances() string {
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.ApiURL + "/balances"
		httpClient := &http.Client{}

		req, err := http.NewRequest("GET", finalURL, nil)
		if err != nil {
			log.Fatal("err1", err)
		}

		signature := c.signMessage("GET", "balances", Nonce)

		req.Header.Set("X-SBTC-SIGNATURE", signature)
		req.Header.Set("X-SBTC-APIKEY", c.ApiKey)
		req.Header.Set("X-SBTC-NONCE", Nonce)

		res, err := httpClient.Do(req)
		if err != nil {
			log.Fatal("err2", err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("err3", err)
		}
		return string(body)
	}

	return "NOT AUTHENTICATED"

}
