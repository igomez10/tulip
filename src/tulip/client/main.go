package client

import (
	"fmt"
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

// Client is the struct that interacts with buda server and executes the requests
type Client struct {
	APIURL        string
	APIKey        string
	APISecret     string
	Nonce         int32
	Authenticated bool
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() string {
	finalURL := c.APIURL + "/markets"
	return execute("GET", finalURL, "", "", "")
}

// GetTicker Returns info about a specific market
func (c *Client) GetTicker(ticker string) string {
	finalURL := c.APIURL + "/markets/" + ticker
	return execute("GET", finalURL, "", "", "")
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) string {
	finalURL := c.APIURL + "/markets/" + marketID + "/order_book"
	return execute("GET", finalURL, "", "", "")
}

//GetTrades returns a list of recent trades
func (c *Client) GetTrades(marketID string) string {
	finalURL := c.APIURL + "/markets/" + marketID + "/trades"
	return execute("GET", finalURL, "", "", "")
}

func execute(method string, completeURL string, apikey string, signature string, Nonce string) string {
	responseData := "An error ocurred, ocurred, check the request method, only GET is allowed at the moment. Also check your apikey ,signature and nonce"
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", completeURL, nil)
	if err != nil {
		log.Fatal(("Error creating new request, submitted url was: " +
			completeURL + "with method" + method +
			". Error: "), err)
	}

	if method == "GET" && apikey == "" && signature == "" && Nonce == "" {

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
		responseData = string(body)

		// TODO ADD POST METHODS
	} else if method == "GET" && apikey != "" && signature != "" && Nonce != "" {
		req.Header.Set("X-SBTC-SIGNATURE", signature)
		req.Header.Set("X-SBTC-APIKEY", apikey)
		req.Header.Set("X-SBTC-NONCE", Nonce)
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
		responseData = string(body)
	}
	return responseData
}

// HERE STARTS THE PRIVATE CALLS

func signMessage(APISecret string, method string, query string, Nonce string) string {
	stringMessage := method + " /api/v2/" + query + " " + Nonce
	fmt.Println(stringMessage)
	key := []byte(APISecret)
	h := hmac.New(sha512.New384, key)
	h.Write([]byte(stringMessage))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

// GetBalances get the wallet balances in all cryptocurrencies and fiat currencies
func (c *Client) GetBalances() string {
	const method string = "GET"
	const query string = "balances"
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		signature := signMessage(c.APISecret, method, query, Nonce)
		return execute("GET", finalURL, c.APIKey, signature, Nonce)
	}
	return "AUTHENTICATION REQUIRED GetBalances"
}

// GetBalance get the wallet balance in a specific cryptocurrency or fiat currency
func (c *Client) GetBalance(currency string) string {
	const method string = "GET"
	var query = "balances/" + currency
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		signature := signMessage(c.APISecret, method, query, Nonce)
		return execute("GET", finalURL, c.APIKey, signature, Nonce)
	}
	return "AUTHENTICATION REQUIRED GetBalance"
}

// GetOrders get the wallet balance in a specific cryptocurrency or fiat currency
func (c *Client) GetOrders(marketID string, per int, page int, state string, minimumExchanged float64) string {
	const method string = "GET"
	var query = "markets/" + marketID + "/orders?per=" + strconv.Itoa(per) + "&page=" + strconv.Itoa(page) + "&state=" + state + "&minimumExchanged=" + strconv.FormatFloat(minimumExchanged, 'g', 20, 32)
	fmt.Println("query " + query)
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce)
		return execute("GET", finalURL, c.APIKey, signature, Nonce)
	}
	return "AUTHENTICATION REQUIRED GetOrders"
}
