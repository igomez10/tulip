package tulip

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
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
	u := path.Join(c.APIURL, "markets")
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
	u := path.Join(c.APIURL, "markets", ticker)
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
	u := path.Join(c.APIURL, "markets", marketID, "order_book")
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
	u := path.Join(c.APIURL, "markets", marketID, "trades")
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
	return c.GetBalance("")
}

// GetBalance get the wallet balance in a specific cryptocurrency or fiat currency
func (c *Client) GetBalance(currency string) (BalanceResponse, error) {
	u := path.Join(c.APIURL, "balances", currency)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return BalanceResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return BalanceResponse{}, err
	}

	var jsonBalance BalanceResponse
	err = json.Unmarshal([]byte(resp), &jsonBalance)
	return jsonBalance, err
}

// GetOrders gets your orders made in a specific market with a specific status
func (c *Client) GetOrders(marketID, state string, per, page int, minimumExchanged float64) (MyOrdersResponse, error) {
	const method string = "GET"

	baseQuery := "%s/markets/%s/orders?per=%d&page=%d&state=%s%minimumExchanged=%f"
	u := fmt.Sprintf(baseQuery, c.APIURL, marketID, per, page, state, minimumExchanged)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return MyOrdersResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return MyOrdersResponse{}, err
	}

	var jsonOrders MyOrdersResponse
	err = json.Unmarshal([]byte(resp), &jsonOrders)
	return jsonOrders, err
}

// PostOrder creates a new order (bid or ask) in a specific market
func (c *Client) PostOrder(marketID, orderType, priceType string, limit, amount float64) (OrderResponse, error) {
	var newOrder interface{}
	switch priceType {
	case "market":
		newOrder = MarketOrder{orderType, priceType, amount}
	case "limit":
		newOrder = LimitOrder{orderType, priceType, limit, amount}
	default:
		return OrderResponse{}, fmt.Errorf(`Only "limit" and "market" supported - invalid %q`, priceType)
	}

	myOrder, err := json.Marshal(newOrder)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("Invalid order values, cannot write as json \n%+v", err)
	}

	var payload []byte
	base64.StdEncoding.Encode(payload, myOrder)

	u := fmt.Sprintf("%s/markets/%s/orders", c.APIURL, marketID)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(payload))
	if err != nil {
		return OrderResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return OrderResponse{}, err
	}

	var o OrderResponse
	err = json.Unmarshal([]byte(resp), &o)
	return o, err
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

	var o OrderResponse
	err = json.Unmarshal([]byte(resp), &o)
	return o, err
}

// GetOrder returns the current state of the order
func (c *Client) GetOrder(orderID string) (OrderResponse, error) {
	u := fmt.Sprintf("%s/orders/%s", c.APIURL, orderID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return OrderResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return OrderResponse{}, err
	}

	var o OrderResponse
	err = json.Unmarshal([]byte(resp), &o)
	return o, err
}

// GetDepositHistory returns the historic deposits
func (c *Client) GetDepositHistory(currency string) (HistoricDespositsResponse, error) {
	u := fmt.Sprintf("%s/currencies/%s/deposits", c.APIURL, currency)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return HistoricDespositsResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return HistoricDespositsResponse{}, err
	}

	var h HistoricDespositsResponse
	err = json.Unmarshal([]byte(resp), &h)
	return h, err
}

// GetWithdrawHistory returns the historic withdrawls
func (c *Client) GetWithdrawHistory(currency string) (HistoricWithdrawResponse, error) {

	u := fmt.Sprintf("%s/currencies/%s/withdrawals", c.APIURL, currency)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return HistoricWithdrawResponse{}, err
	}

	resp, err := c.execute(req)
	if err != nil {
		return HistoricWithdrawResponse{}, err
	}

	var h HistoricWithdrawResponse
	err = json.Unmarshal([]byte(resp), &h)
	return h, err
}
