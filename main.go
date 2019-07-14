package tulip

import (
	"bytes"
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
	httpClient    http.Client
}

// NewClient returns a new Client
func NewClient(key string, APISecret string) *Client {
	var newClient Client
	if key != "" && APISecret != "" {
		newClient.Authenticated = true
	}

	newClient.APIURL = APIURL
	newClient.APIKey = key
	newClient.APISecret = APISecret
	newClient.httpClient = http.Client{Timeout: time.Second * 10}
	return &newClient
}

func getNonce() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// GetMarkets Returns info about all markets
func (c *Client) GetMarkets() (MarketsResponse, error) {
	u := fmt.Sprintf("%s/markets", c.APIURL)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return MarketsResponse{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return MarketsResponse{}, err
	}

	var jsonMarketsResponse MarketsResponse
	err = json.Unmarshal([]byte(res), &jsonMarketsResponse)
	return jsonMarketsResponse, err
}

// GetTicker Returns the exchange rate for a given ticker
func (c *Client) GetTicker(ticker string) (MarketResponse, error) {
	u := fmt.Sprintf("%s/markets/%s", c.APIURL, ticker)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return MarketResponse{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return MarketResponse{}, err
	}

	var jsonMarketResponse MarketResponse
	err = json.Unmarshal([]byte(res), &jsonMarketResponse)
	return jsonMarketResponse, nil
}

//GetOrderBook is used to get current state of the market.
// It shows the best offers (bid, ask) and the price from the
// last transaction, daily volume and the price in the last 24 hours
func (c *Client) GetOrderBook(marketID string) (OrderBook, error) {
	u := fmt.Sprintf("%s/markets/%s/order_book", c.APIURL, marketID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return OrderBook{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return OrderBook{}, err
	}

	var orderBook OrderBook
	err = json.Unmarshal([]byte(res), &orderBook)
	return orderBook, nil
}

//GetTrades returns a list of recent trades in a given market
func (c *Client) GetTrades(marketID string) (TradesResponse, error) {
	u := fmt.Sprintf("%s/markets/%s/trades", c.APIURL, marketID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return TradesResponse{}, err
	}

	res, err := c.execute(req)
	if err != nil {
		return TradesResponse{}, err
	}

	var jsonTrades TradesResponse
	err = json.Unmarshal([]byte(res), &jsonTrades)
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

func (c *Client) execute(r *http.Request) (string, error) {
	nonce := getNonce()
	r.Header.Set("X-SBTC-APIKEY", c.APIKey)
	r.Header.Set("X-SBTC-NONCE", nonce)
	r.Header.Set("Content-Type", "application/json")

	if c.Authenticated {
		signature, err := c.signMessage(r, nonce)
		if err != nil {
			return "", err
		}
		r.Header.Set("X-SBTC-SIGNATURE", signature)
	}

	res, err := c.httpClient.Do(r)
	if err != nil {
		return "", fmt.Errorf("error on http request \n %+v", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response \n %+v", err)
	}

	return string(body), nil
}

// HERE STARTS THE PRIVATE CALLS

func (c *Client) signMessage(r *http.Request, nonce string) (string, error) {
	// {GET|POST|PUT} {path} {base64_encoded_body} {nonce}
	p := strings.TrimPrefix(r.URL.String(), c.APIURL)
	b := fmt.Sprintf("%s %s %s %s", r.Method, p, r.Body, nonce)

	key := []byte(c.APISecret)
	h := hmac.New(sha512.New384, key)

	_, err := h.Write([]byte(b))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

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
func (c *Client) GetBalances() (BalanceResponse, error) {
	u := fmt.Sprintf("%s/balances", c.APIURL)
	req, err := http.NewRequest("GET", u, nil)
	res, err := c.execute(req)
	if err != nil {
		return BalanceResponse{}, err
	}

	var jsonBalances BalanceResponse
	err = json.Unmarshal([]byte(res), &jsonBalances)
	if err != nil {
		return jsonBalances, err
	}

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
	var encodedRequestPayload []byte
	base64.StdEncoding.Encode(encodedRequestPayload, []byte(`{ "state": "canceling" }`))
	u := fmt.Sprintf("%s/orders/%s", c.APIURL, orderID)
	req, err := http.NewRequest("PUT", u, bytes.NewBuffer(encodedRequestPayload))
	if err != nil {
		return OrderResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return OrderResponse{}, err
	}

	var jsonCancelOrder OrderResponse
	err = json.Unmarshal([]byte(resp), &jsonCancelOrder)
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
