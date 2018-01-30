package tulip

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"types"
)

// Client is the struct that interacts with buda server and executes the requests
type client struct {
	apiURL        string
	apiKey        string
	apiSecret     string
	authenticated bool
}

// order is the struct that will be sent in the request payload as a json to create a new order
type order struct {
	OrderType string  `json:"type"`
	PriceType string  `json:"price_type"`
	Limit     float64 `json:"limit"`
	Amount    float64 `json:"amount"`
}

// CreateClient returns a new client
func CreateClient(apikey string, apisecret string) *client {
	var newClient client
	if apikey != "" && apisecret != "" {
		newClient.authenticated = true
	}
	newClient.apiURL = "https://www.surbtc.com/api/v2"
	newClient.apiKey = apikey
	newClient.apiSecret = apisecret

	return &newClient
}

// GetMarkets Returns info about all markets
func (c *client) GetMarkets() (types.MarketsResponse, error) {
	finalURL := c.apiURL + "/markets"
	resp := execute("GET", finalURL, "", "", "", "")
	var jsonMarketsResponse types.MarketsResponse
	err := json.Unmarshal([]byte(resp), &jsonMarketsResponse)
	if err != nil {
		return jsonMarketsResponse, err
	}
	return jsonMarketsResponse, nil
}

// GetTicker Returns info about a specific market
func (c *client) GetTicker(ticker string) (types.MarketResponse, error) {
	finalURL := c.apiURL + "/markets/" + ticker
	resp := execute("GET", finalURL, "", "", "", "")
	var jsonMarketResponse types.MarketResponse
	err := json.Unmarshal([]byte(resp), &jsonMarketResponse)
	if err != nil {
		return jsonMarketResponse, err
	}
	return jsonMarketResponse, nil
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *client) GetOrderBook(marketID string) (types.OrderBook, error) {
	finalURL := c.apiURL + "/markets/" + marketID + "/order_book"
	resp := execute("GET", finalURL, "", "", "", "")
	var orderBook types.OrderBook
	err := json.Unmarshal([]byte(resp), &orderBook)
	if err != nil {
		return orderBook, err
	}
	return orderBook, nil
}

//GetTrades returns a list of recent trades
func (c *client) GetTrades(marketID string) (types.TradesResponse, error) {
	finalURL := c.apiURL + "/markets/" + marketID + "/trades"
	resp := execute("GET", finalURL, "", "", "", "")
	var jsonTrades types.TradesResponse
	err := json.Unmarshal([]byte(resp), &jsonTrades)
	if err != nil {
		return jsonTrades, err
	}
	return jsonTrades, nil
}

func execute(method string, completeURL string, apikey string, signature string, Nonce string, reqPayload string) string {
	// responseData will contain the body of the response from the server, execute(...) will return this variable as a string
	responseData := "An error ocurred, check the request method, check your apikey ,signature and nonce"
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

		// res is the response from the server when the request is executed
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
func (c *client) GetBalances() (types.BalancesResponse, error) {
	const method string = "GET"
	const query string = "balances"
	var jsonBalances types.BalancesResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonBalances)
		if err != nil {
			return jsonBalances, err
		}
		return jsonBalances, nil
	}
	err := errors.New("Authentication Required GetBalances")
	return jsonBalances, err
}

// GetBalance get the wallet balance in a specific cryptocurrency or fiat currency
func (c *client) GetBalance(currency string) (types.BalanceResponse, error) {
	const method string = "GET"
	var query = "balances/" + currency
	var jsonBalance types.BalanceResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonBalance)
		if err != nil {
			return jsonBalance, err
		}
		return jsonBalance, nil
	}
	err := errors.New("Authentication Required GetBalance")
	return jsonBalance, err
}

// GetOrders gets your orders made in a specific market with a specific status
func (c *client) GetOrders(marketID string, per int, page int, state string, minimumExchanged float64) (types.MyOrdersResponse, error) {
	const method string = "GET"
	var query = "markets/" + marketID + "/orders?per=" + strconv.Itoa(per) + "&page=" + strconv.Itoa(page) + "&state=" + state + "&minimumExchanged=" + strconv.FormatFloat(minimumExchanged, 'g', 20, 32)
	var jsonOrders types.MyOrdersResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonOrders)
		if err != nil {
			return jsonOrders, err
		}
		return jsonOrders, nil
	}
	err := errors.New("Authentication Required GetOrders")
	return jsonOrders, err
}

// PostOrder creates a new order (bid or ask) in a specific market
func (c *client) PostOrder(marketID string, orderType string, priceType string, limit float64, amount float64) (types.OrderResponse, error) {
	const method string = "POST"
	var query = "markets/" + marketID + "/orders"
	var jsonPostOrder types.OrderResponse
	neworder := &order{OrderType: orderType, PriceType: priceType, Limit: limit, Amount: amount}
	myOrder, err := json.Marshal(neworder)
	if err != nil {
		fmt.Println("THE ORDER HAS WRONG VALUES, CHECK THE API DOCUMENTATION" + marketID + " , " + orderType + ", " + priceType + ", " + strconv.FormatFloat(limit, 'g', 20, 32) + ", " + strconv.FormatFloat(amount, 'g', 20, 32))
		return jsonPostOrder, err
	}
	encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(myOrder))
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		signature := signMessage(c.apiSecret, method, query, Nonce, encodedRequestPayload)
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, string(myOrder))
		err := json.Unmarshal([]byte(resp), &jsonPostOrder)
		if err != nil {
			return jsonPostOrder, err
		}
		return jsonPostOrder, nil
	}
	err2 := errors.New("AUTHENTICATION REQUIRED PostOrder")
	return jsonPostOrder, err2
}

// CancelOrder cancels a specified order
func (c *client) CancelOrder(orderID string) (types.OrderResponse, error) {
	const method string = "PUT"
	var query = "orders/" + orderID
	requestPayloadString := "{" + `"state":` + `"canceling"` + "}"
	encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(requestPayloadString))
	var jsonCancelOrder types.OrderResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		signature := signMessage(c.apiSecret, method, query, Nonce, encodedRequestPayload)
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, requestPayloadString)
		err := json.Unmarshal([]byte(resp), &jsonCancelOrder)
		if err != nil {
			return jsonCancelOrder, err
		}
		return jsonCancelOrder, nil
	}
	err2 := errors.New("AUTHENTICATION REQUIRED CancelOrder")
	return jsonCancelOrder, err2
}

// GetOrder returns the current state of the order
func (c *client) GetOrder(orderID string) (types.OrderResponse, error) {
	const method string = "GET"
	var query = "orders/" + orderID
	var jsonOrder types.OrderResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonOrder)
		if err != nil {
			return jsonOrder, err
		}
		return jsonOrder, nil
	}
	err2 := errors.New("AUTHENTICATION REQUIRED GetOrder")
	return jsonOrder, err2

}

// GetDepositHistory returns the historic deposits
func (c *client) GetDepositHistory(currency string) (types.HistoricDespositsResponse, error) {
	const method string = "GET"
	var query = "currencies/" + currency + "/deposits"
	var jsonHistoricDeposit types.HistoricDespositsResponse
	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonHistoricDeposit)
		if err != nil {
			return jsonHistoricDeposit, err
		}
		return jsonHistoricDeposit, nil
	}
	err2 := errors.New("AUTHENTICATION REQUIRED GetDepositHistory")
	return jsonHistoricDeposit, err2
}

// GetWithdrawHistory returns the historic withdrawls
func (c *client) GetWithdrawHistory(currency string) (types.HistoricWithdrawResponse, error) {
	const method string = "GET"
	var query = "currencies/" + currency + "/withdrawals"
	var jsonHistoricWithdraw types.HistoricWithdrawResponse

	if c.authenticated == true {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := c.apiURL + "/" + query
		fmt.Println("finalurl " + finalURL)
		signature := signMessage(c.apiSecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.apiKey, signature, Nonce, "")
		err := json.Unmarshal([]byte(resp), &jsonHistoricWithdraw)
		if err != nil {
			return jsonHistoricWithdraw, err
		}
		return jsonHistoricWithdraw, nil
	}
	err2 := errors.New("AUTHENTICATION REQUIRED GetWithdrawHistory")
	return jsonHistoricWithdraw, err2
}
