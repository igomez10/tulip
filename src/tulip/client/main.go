package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	// "fmt"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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

type order struct {
	OrderType string  `json:"type"`
	PriceType string  `json:"price_type"`
	Limit     float64 `json:"limit"`
	Amount    float64 `json:"amount"`
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() string {
	finalURL := c.APIURL + "/markets"
	return execute("GET", finalURL, "", "", "", "")
}

// GetTicker Returns info about a specific market
func (c *Client) GetTicker(ticker string) string {
	finalURL := c.APIURL + "/markets/" + ticker
	return execute("GET", finalURL, "", "", "", "")
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) string {
	finalURL := c.APIURL + "/markets/" + marketID + "/order_book"
	return execute("GET", finalURL, "", "", "", "")
}

//GetTrades returns a list of recent trades
func (c *Client) GetTrades(marketID string) string {
	finalURL := c.APIURL + "/markets/" + marketID + "/trades"
	return execute("GET", finalURL, "", "", "", "")
}

func execute(method string, completeURL string, apikey string, signature string, Nonce string, reqPayload string) string {
	// responseData will contain the body of the response from the server, execute(...) will return this variable as a string
	responseData := "An error ocurred, ocurred, check the request method, only GET is allowed at the moment. Also check your apikey ,signature and nonce"
	// httpClient	will make the http requests to the server
	httpClient := &http.Client{}
	// req is the request that will hold all the info
	req, err := http.NewRequest(method, completeURL, nil)
	if err != nil {
		log.Fatal(("Error creating new request, submitted url was: " +
			completeURL + "with method" + method +
			". Error: "), err)
	}

	if method == "GET" && apikey == "" && signature == "" && Nonce == "" {
		// The GET requests that do not require authentication will end here

		// res is the respons from the server when the request is executed
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

	} else if apikey != "" && signature != "" && Nonce != "" {
		// The POST AND GET requests that DO NEED AUTHENTICATION will end here
		req.Header.Set("X-SBTC-SIGNATURE", signature)
		req.Header.Set("X-SBTC-APIKEY", apikey)
		req.Header.Set("X-SBTC-NONCE", Nonce)
		req.Header.Set("Content-Type", "application/json")

		if method == "POST" || method == "PUT" {
			req.Body = ioutil.NopCloser(strings.NewReader(reqPayload))
		}

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

func signMessage(APISecret string, method string, query string, Nonce string, body string) string {
	var stringMessage string

	if body != "" {
		stringMessage = method + " /api/v2/" + query + " " + body + " " + Nonce
	} else {
		stringMessage = method + " /api/v2/" + query + " " + Nonce
	}

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
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
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
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
	}
	return "AUTHENTICATION REQUIRED GetBalance"
}

// GetOrders get the wallet balance in a specific cryptocurrency or fiat currency
func (c *Client) GetOrders(marketID string, per int, page int, state string, minimumExchanged float64) string {
	const method string = "GET"
	var query = "markets/" + marketID + "/orders?per=" + strconv.Itoa(per) + "&page=" + strconv.Itoa(page) + "&state=" + state + "&minimumExchanged=" + strconv.FormatFloat(minimumExchanged, 'g', 20, 32)
	// fmt.Println("query " + query)
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		// fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
	}
	return "AUTHENTICATION REQUIRED GetOrders"
}

// PostOrder creates a new order (bid or ask) in a specific market
func (c *Client) PostOrder(marketID string, orderType string, priceType string, limit float64, amount float64) string {
	const method string = "POST"
	var query = "markets/" + marketID + "/orders"

	neworder := &order{OrderType: orderType, PriceType: priceType, Limit: limit, Amount: amount}
	myOrder, err := json.Marshal(neworder)
	fmt.Println(string(myOrder))
	if err != nil {
		return "THE ORDER HAS WRONG VALUES, CHECK THE API DOCUMENTATION" + marketID + " , " + orderType + ", " + priceType + ", " + strconv.FormatFloat(limit, 'g', 20, 32) + ", " + strconv.FormatFloat(amount, 'g', 20, 32)
	}
	encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(myOrder))
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce, encodedRequestPayload)
		return execute(method, finalURL, c.APIKey, signature, Nonce, string(myOrder))
	}
	return "AUTHENTICATION REQUIRED PostOrder"
}

// CancelOrder cancels a specified order
func (c *Client) CancelOrder(orderID string) string {
	const method string = "PUT"
	var query = "orders/" + orderID

	requestPayloadString := "{" + `"state":` + `"canceling"` + "}"

	encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(requestPayloadString))
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		signature := signMessage(c.APISecret, method, query, Nonce, encodedRequestPayload)
		return execute(method, finalURL, c.APIKey, signature, Nonce, requestPayloadString)
	}
	return "AUTHENTICATION REQUIRED CancelOrder"
}

// GetOrder returns the current state of the order
func (c *Client) GetOrder(orderID string) string {
	const method string = "GET"
	var query = "orders/" + orderID
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
	}
	return "AUTHENTICATION REQUIRED GetOrder"
}

// GetDepositHistory returns the historic deposits/withdrawls
func (c *Client) GetDepositHistory(currency string) string {
	const method string = "GET"
	var query = "currencies/" + currency + "/deposits"
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
	}
	return "AUTHENTICATION REQUIRED GetDepositHistory"
}

// GetWithdrawHistory returns the historic deposits/withdrawls
func (c *Client) GetWithdrawHistory(currency string) string {
	const method string = "GET"
	var query = "currencies/" + currency + "/withdrawals"
	if c.Authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.APIURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.APISecret, method, query, Nonce, "")
		return execute(method, finalURL, c.APIKey, signature, Nonce, "")
	}
	return "AUTHENTICATION REQUIRED GetWithdrawHistory"
}
