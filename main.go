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
)

// APIURL is the url used to call Buda services
const APIURL = "https://www.buda.com/api/v2"

// Client is the struct that interacts with buda server and executes the requests
type Client struct {
	APIURL        string
	APIKey        string
	APISecret     string
	Authenticated bool
}

// NewClient returns a new Client
func NewClient(key string, apisecret string) *Client {
	var newClient Client
	if key != "" && apisecret != "" {
		newClient.Authenticated = true
	}

	newClient.APIURL = APIURL
	newClient.APIKey = key
	newClient.APISecret = apisecret

	return &newClient
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() (MarketsResponse, error) {
	finalURL := fmt.Sprintf("%s/markets", c.APIURL)

	resp := execute("GET", finalURL, "", "", "", "")

	var jsonMarketsResponse MarketsResponse
	err := json.Unmarshal([]byte(resp), &jsonMarketsResponse)
	if err != nil {
		return jsonMarketsResponse, err
	}

	return jsonMarketsResponse, nil
}

// GetTicker Returns the exchange rate for a given ticker
func (c *Client) GetTicker(ticker string) (MarketResponse, error) {
	finalURL := fmt.Sprintf("%s/markets/%s", c.APIURL, ticker)

	resp := execute("GET", finalURL, "", "", "", "")

	var jsonMarketResponse MarketResponse
	err := json.Unmarshal([]byte(resp), &jsonMarketResponse)
	if err != nil {
		return jsonMarketResponse, err
	}

	return jsonMarketResponse, nil
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) (OrderBook, error) {
	finalURL := fmt.Sprintf("%s/markets/%s/order_book", c.APIURL, marketID)
	resp := execute("GET", finalURL, "", "", "", "")

	var orderBook OrderBook
	err := json.Unmarshal([]byte(resp), &orderBook)
	if err != nil {
		return orderBook, err
	}

	return orderBook, nil
}

//GetTrades returns a list of recent trades in a given market
func (c *Client) GetTrades(marketID string) (TradesResponse, error) {
	finalURL := fmt.Sprintf("%s/markets/%s/trades", c.APIURL, marketID)
	resp := execute("GET", finalURL, "", "", "", "")
	var jsonTrades TradesResponse

	err := json.Unmarshal([]byte(resp), &jsonTrades)
	if err != nil {
		return jsonTrades, err
	}

	return jsonTrades, nil
}

func execute(method, completeURL, key, signature, Nonce, reqPayload string) string {
	// responseData will contain the body of the response
	// from the server, execute(...) will return this variable as a string
	responseData := "Error - check the request method, check your apikey ,signature and nonce"

	// httpClient	will make the http requests to the server
	httpClient := &http.Client{}

	// req is the request that will contain all the info
	req, err := http.NewRequest(method, completeURL, nil)
	if err != nil {
		log.Fatalf("Error creating new request: \n %+v ", err)
	}

	if method == "GET" && key == "" && signature == "" && Nonce == "" {
		// The GET requests that do not require authentication will end here

		// res is the response from the server when the request is executed
		res, err := httpClient.Do(req)
		if err != nil {
			// TODO return error instead of fatal
			log.Fatalf("Error executing new request \n %+v", err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response \n %+v", err)
		}

		responseData = string(body)

	} else if key != "" && signature != "" && Nonce != "" {
		// The POST AND GET requests that DO NEED AUTHENTICATION will end here
		req.Header.Set("X-SBTC-SIGNATURE", signature)
		req.Header.Set("X-SBTC-APIKEY", key)
		req.Header.Set("X-SBTC-NONCE", Nonce)
		req.Header.Set("Content-Type", "application/json")

		if method == "POST" || method == "PUT" {
			req.Body = ioutil.NopCloser(strings.NewReader(reqPayload))
		}

		res, err := httpClient.Do(req)
		if err != nil {
			log.Fatalf("Error executing new request \n %s", err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response \n %+v", err)
		}

		responseData = string(body)
	}

	return responseData
}

// HERE STARTS THE PRIVATE CALLS

func signMessage(APISecret, method, query, Nonce, body string) string {
	var stringMessage string

	if body != "" {
		stringMessage = fmt.Sprintf("%s/api/v2/%s/%s/%s", method, query, body, Nonce)
	} else {
		stringMessage = fmt.Sprintf("%s/api/v2/%s/%s", method, query, Nonce)
	}

	key := []byte(APISecret)
	h := hmac.New(sha512.New384, key)
	h.Write([]byte(stringMessage))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

// GetBalances get the wallet balances in all cryptocurrencies and fiat currencies
func (c *Client) GetBalances() (BalancesResponse, error) {
	var jsonBalances BalancesResponse
	method := "GET"
	query := "balances"

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)
		signature := signMessage(c.APISecret, method, query, Nonce, "")

		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

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
func (c *Client) GetBalance(currency string) (BalanceResponse, error) {
	var jsonBalance BalanceResponse
	method := "GET"
	query := fmt.Sprintf("balances/%s", currency)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)
		signature := signMessage(c.APISecret, method, query, Nonce, "")

		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

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
func (c *Client) GetOrders(marketID string, per, page int, state string, minimumExchanged float64) (MyOrdersResponse, error) {
	var jsonOrders MyOrdersResponse
	const method string = "GET"

	baseQuery := "markets/%s/orders?per=%d&page=%d&state=%s%minimumExchanged=%f"
	query := fmt.Sprintf(baseQuery, marketID, per, page, state, minimumExchanged)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)
		signature := signMessage(c.APISecret, method, query, Nonce, "")

		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

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
func (c *Client) PostOrder(marketID, orderType, priceType string, limit, amount float64) (OrderResponse, error) {
	var jsonPostOrder OrderResponse

	if c.Authenticated {
		method := "POST"
		query := fmt.Sprintf("markets/%s/orders", marketID)
		var newOrder interface{}

		if priceType == "market" {
			newOrder = MarketOrder{orderType, priceType, amount}
		} else if priceType == "limit" {
			newOrder = LimitOrder{orderType, priceType, limit, amount}
		}

		myOrder, err := json.Marshal(newOrder)
		if err != nil {
			fmt.Printf("Unexpected error marshaling order values, check API docs \n %+v", err)
			return jsonPostOrder, err
		}

		encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(myOrder))
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)
		signature := signMessage(c.APISecret, method, query, Nonce, encodedRequestPayload)

		resp := execute(method, finalURL, c.APIKey, signature, Nonce, string(myOrder))

		err = json.Unmarshal([]byte(resp), &jsonPostOrder)
		if err != nil {
			err := fmt.Errorf("Error unmarshaling response from API, check docs\n %+v", err)
			return jsonPostOrder, err
		}

		return jsonPostOrder, nil
	}

	err := errors.New("AUTHENTICATION REQUIRED PostOrder")
	return jsonPostOrder, err
}

// CancelOrder cancels a specified order
func (c *Client) CancelOrder(orderID string) (OrderResponse, error) {
	var jsonCancelOrder OrderResponse
	method := "PUT"
	query := fmt.Sprintf("orders/%s", orderID)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)
		requestPayloadString := `{ "state": "canceling" }`
		encodedRequestPayload := base64.StdEncoding.EncodeToString([]byte(requestPayloadString))

		signature := signMessage(c.APISecret, method, query, Nonce, encodedRequestPayload)
		resp := execute(method, finalURL, c.APIKey, signature, Nonce, requestPayloadString)

		err := json.Unmarshal([]byte(resp), &jsonCancelOrder)
		if err != nil {
			return jsonCancelOrder, err
		}

		return jsonCancelOrder, nil
	}

	err := errors.New("AUTHENTICATION REQUIRED CancelOrder")
	return jsonCancelOrder, err
}

// GetOrder returns the current state of the order
func (c *Client) GetOrder(orderID string) (OrderResponse, error) {
	var jsonOrder OrderResponse
	method := "GET"
	query := fmt.Sprintf("orders/%s", orderID)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)

		signature := signMessage(c.APISecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

		err := json.Unmarshal([]byte(resp), &jsonOrder)
		if err != nil {
			return jsonOrder, err
		}

		return jsonOrder, nil
	}

	err := errors.New("AUTHENTICATION REQUIRED GetOrder")
	return jsonOrder, err
}

// GetDepositHistory returns the historic deposits
func (c *Client) GetDepositHistory(currency string) (HistoricDespositsResponse, error) {
	var jsonHistoricDeposit HistoricDespositsResponse
	method := "GET"
	query := fmt.Sprintf("currencies/%s/deposits", currency)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)

		signature := signMessage(c.APISecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

		err := json.Unmarshal([]byte(resp), &jsonHistoricDeposit)
		if err != nil {
			return jsonHistoricDeposit, err
		}

		return jsonHistoricDeposit, nil
	}

	err := errors.New("AUTHENTICATION REQUIRED GetDepositHistory")
	return jsonHistoricDeposit, err
}

// GetWithdrawHistory returns the historic withdrawls
func (c *Client) GetWithdrawHistory(currency string) (HistoricWithdrawResponse, error) {
	var jsonHistoricWithdraw HistoricWithdrawResponse
	method := "GET"
	query := fmt.Sprintf("currencies/%s/withdrawals", currency)

	if c.Authenticated {
		Nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		finalURL := fmt.Sprintf("%s/%s", c.APIURL, query)

		signature := signMessage(c.APISecret, method, query, Nonce, "")
		resp := execute(method, finalURL, c.APIKey, signature, Nonce, "")

		err := json.Unmarshal([]byte(resp), &jsonHistoricWithdraw)
		if err != nil {
			return jsonHistoricWithdraw, err
		}

		return jsonHistoricWithdraw, nil
	}

	err := errors.New("AUTHENTICATION REQUIRED GetWithdrawHistory")
	return jsonHistoricWithdraw, err
}
