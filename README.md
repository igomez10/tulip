# 🌷 Tulip 🌷



An open source API client for Buda written in Go.

You can check the current Buda API documentation at http://api.buda.com

NOTE: You need a developer account in order to make full use of this client.
You can make unauthenticated calls (http://api.buda.com/#llamadas-p-blicas)


### Available Methods ✅
- #### Public 🔓
  - Get markets info `GetMarkets()`
  - Get ticker info `GetTicker()`
  - Get all available orders for a market `GetOrderBook()`
  - Get info about recent trades `GetTrades()`
- #### Private (Requires APIkey) 🔑
  - Get your current balance `GetBalances()`
  - Get your orders `GetOrders()`
  - Crate a new order (bid or ask) `PostOrder()`
  - Cancel a current order `CancelOrder()`
  - Get info about a specific order `GetOrder()`
  - Get historic info about deposits `GetDepositHistory()`
  - Get historic info about withdrawals `GetWithdrawHistory()`


#### TODO 👨🏻‍💻

  - Do Deposits and Withdrawals
  - Do string interpolation for queries instead of appending strings together
  - ?? Make a suggestion

## Get it working 🦄

In your terminal:
```sh
$ go get github.com/igomez10/tulip
```


In your .go file:
```GO
package main

import (
	"fmt"
	"os"

	"github.com/igomez10/tulip"
)

func main() {

	APIKey := os.Getenv("BUDAKEY")       // you can modify this variable and hardocde your own apikey
	APISecret := os.Getenv("BUDASECRET") // you can modify this variable and hardcode your own apisecret

	buda := tulip.CreateClient(APIKey, APISecret)
	results, err := buda.GetTicker("btc-clp")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(results.Market.ID)                 // "BTC-CLP"
		fmt.Println(results.Market.Name)               // "btc-clp"
		fmt.Println(results.Market.BaseCurrency)       // "BTC"
		fmt.Println(results.Market.QuoteCurrency)      // "CLP"
		fmt.Println(results.Market.MinimumOrderAmount) // ["0.0001 BTC"]
	}
```



# Documentation 📄

## Public Calls

### Get Ticker
GetTicker Returns info about a specific market

```GO
buda.GetTicker(<marketID>)
```

###### marketID can be "btc-clp", "btc-cop" etc

```json
{
   "market":{
      "id":"BTC-CLP",
      "name":"btc-clp",
      "base_currency":"BTC",
      "quote_currency":"CLP",
      "minimum_order_amount":[
         "0.0001",
         "BTC"
      ]
   }
}
```


### Get Markets
GetMarkets Returns info about all markets
```GO
buda.GetMarkets()
 ```

 ```json
 {
    "markets":[
       {
          "id":"BTC-CLP",
          "name":"btc-clp",
          "base_currency":"BTC",
          "quote_currency":"CLP",
          "minimum_order_amount":[
             "0.0001",
             "BTC"
          ]
       },
       {
          "id":"BCH-PEN",
          "name":"bch-pen",
          "base_currency":"BCH",
          "quote_currency":"PEN",
          "minimum_order_amount":[
             "0.0001",
             "BCH"
          ]
       }
    ]
 }
 ```

### Get Order Book
GetOrderBook is used to get current state of the market. It shows the best offers (bid, ask) and the price from the last transaction, daily volume and the price in the last 24 hours
```GO
buda.GetOrderBook("btc-clp")
```

```json
{
 "order_book":{
    "asks":[
       [
          "7299998.0",
          "0.0374176"
         ],
       [
          "7299999.0",
          "0.07844341"
        ]
      ],
    "bids":[
       [
          "7217231.0",
          "0.02699449"
       ],
       [
          "7183000.0",
          "0.28"
       ],
       [
          "7181619.0",
          "0.0688"
       ]
    ]
 }
}

```
### Get Trades
GetTrades returns a list of recent trades in a specific market

```Go
buda.GetTrades("btc-clp")
```
```json
{  
  "trades":{  
       "market_id":"BTC-CLP",
       "timestamp":null,
       "last_timestamp":"1517177157715",
       "entries":[
       [  
          "1517188555937",
          "0.00300551",
          "7217231.0",
          "sell",
          316863
        ],
        [  
          "1517186567638",
          "0.00206959",
          "7300000.0",
          "buy",
          316844
        ]
      ]
    }
}
```

## Private Calls
#### NEEDS AN APIKEY/APISECRET

### Get Balances/Balances
  - GetBalances gets the wallet balances in all cryptocurrencies and fiat currencies
  - GetBalance(currency) gets the wallet balance in a specific cryptocurrency or fiat currency


```GO
buda.GetBalances()
```
```json
{
   "balances":[
      {
         "id":"BTC",
         "amount":[
            "0.0000001",
            "BTC"
         ],
         "available_amount":[
            "0.0000001",
            "BTC"
         ],
         "frozen_amount":[
            "0.0",
            "BTC"
         ],
         "pending_withdraw_amount":[
            "0.0",
            "BTC"
         ],
         "account_id":1234567
      },
      {
         "id":"COP",
         "amount":[
            "1000.33",
            "COP"
         ],
         "available_amount":[
            "1000.33",
            "COP"
         ],
         "frozen_amount":[
            "0.0",
            "COP"
         ],
         "pending_withdraw_amount":[
            "0.0",
            "COP"
         ],
         "account_id":1234567
      }
   ]
}
```


### Get info about all YOUR Orders
##### Warning: GetOrderBook(A)  !=  GetOrders(A,B,C,D,E)

- GetOrders gets your orders made in a specific market with a specific status or/and minimum amount


```GO
buda.GetOrders(marketID, ordersPerPage, page, state, minimumExchanged)
```

##### Example:
``` GO
buda.GetOrders("btc-cop", 300, 1, "pending", float64(0))
```

```json
{  
   "orders":[  
      {  
         "id":123456,
         "market_id":"BTC-COP",
         "account_id":1234567,
         "type":"Bid",
         "state":"pending",
         "created_at":"2018-01-29T02:28:44.658Z",
         "fee_currency":"BTC",
         "price_type":"limit",
         "limit":[  
            "0.0",
            "COP"
         ],
         "amount":[  
            "0.0001",
            "BTC"
         ],
         "original_amount":[  
            "0.0001",
            "BTC"
         ],
         "traded_amount":[  
            "0.0",
            "BTC"
         ],
         "total_exchanged":[  
            "0.0",
            "COP"
         ],
         "paid_fee":[  
            "0.0",
            "BTC"
         ]
      }
   ],
   "meta":{  
      "total_pages":1,
      "total_count":1,
      "current_page":1
   }
}
```

### Create a new order
#### Warning: Be sure to understand this method: http://api.buda.com/#nueva-orden

#### PostOrder creates a new order (bid or ask) in a specific market

 IMPORTANT: You can either crete a Limit of Market order specifying it on the _priceType_ parameter as "limit" or "market". Therefore, the parameter _limit_ is only read in the first case but it is necessary to make the call. In order to post a Market order, submit the limit as float(0) or float().

```GO
buda.PostOrder(marketID string, orderType string, priceType string, limit float64, amount float64)
```
##### At the moment, only "limit" is available for priceType
##### Example:
``` GO
buda.PostOrder("btc-cop", "bid", "limit", float64(0.00001), float64(0.0001)))
```

```json
{
   "order":{
      "id":1234567,
      "market_id":"BTC-COP",
      "account_id":1234567,
      "type":"Bid",
      "state":"received",
      "created_at":"2018-01-29T02:28:44.658Z",
      "fee_currency":"BTC",
      "price_type":"limit",
      "limit":[
         "0.0",
         "COP"
      ],
      "amount":[
         "0.0001",
         "BTC"
      ],
      "original_amount":[
         "0.0001",
         "BTC"
      ],
      "traded_amount":[
         "0.0",
         "BTC"
      ],
      "total_exchanged":[
         "0.0",
         "COP"
      ],
      "paid_fee":[
         "0.0",
         "BTC"
      ]
   }
}
```


### Cancel an order
#### CancelOrder cancels a specified order

```GO
buda.CancelOrder(orderID string)
```
#### CancelOrder cancels a specified order
```json
{  
   "order":{  
      "id":1234567,
      "market_id":"BTC-COP",
      "account_id":1234567,
      "type":"Bid",
      "state":"canceled",
      "created_at":"2018-01-29T02:28:44.658Z",
      "fee_currency":"BTC",
      "price_type":"limit",
      "limit":[  
         "0.0",
         "COP"
      ],
      "amount":[  
         "0.0001",
         "BTC"
      ],
      "original_amount":[  
         "0.0001",
         "BTC"
      ],
      "traded_amount":[  
         "0.0",
         "BTC"
      ],
      "total_exchanged":[  
         "0.0",
         "COP"
      ],
      "paid_fee":[  
         "0.0",
         "BTC"
      ]
   }
}
```


### Get info about a specific order
####GetOrder returns the current state of a specific order

```GO
buda.GetOrder(orderID string)
```

```json
{  
   "order":{  
      "id":1234567,
      "market_id":"BTC-COP",
      "account_id":1234567,
      "type":"Bid",
      "state":"received",
      "created_at":"2018-01-29T02:38:37.178Z",
      "fee_currency":"BTC",
      "price_type":"limit",
      "limit":[  
         "0.0",
         "COP"
      ],
      "amount":[  
         "0.0001",
         "BTC"
      ],
      "original_amount":[  
         "0.0001",
         "BTC"
      ],
      "traded_amount":[  
         "0.0",
         "BTC"
      ],
      "total_exchanged":[  
         "0.0",
         "COP"
      ],
      "paid_fee":[  
         "0.0",
         "BTC"
      ]
   }
}
```


### Get historic deposits in a specific fiat currency

```GO
buda.GetDepositHistory(currency string)
```

```json
{  
   "deposits":[  
      {  
         "id":1234567,
         "state":"confirmed",
         "currency":"COP",
         "created_at":"2018-01-26T03:01:34.791Z",
         "deposit_data":{  
            "type":"fiat_deposit_data",
            "created_at":"2018-01-26T03:01:34.783Z",
            "updated_at":"2018-01-26T03:01:34.783Z",
            "upload_url":null
         },
         "amount":[  
            "10000.0",
            "COP"
         ],
         "fee":[  
            "431.0",
            "COP"
         ]
      },
      {  
         "id":1234567,
         "state":"confirmed",
         "currency":"COP",
         "created_at":"2018-01-22T17:06:14.473Z",
         "deposit_data":{  
            "type":"fiat_deposit_data",
            "created_at":"2018-01-22T17:06:14.465Z",
            "updated_at":"2018-01-22T17:06:14.465Z",
            "upload_url":null
         },
         "amount":[  
            "1000.0",
            "COP"
         ],
         "fee":[  
            "431.0",
            "COP"
         ]
      }
   ],
   "meta":{  
      "total_pages":1,
      "total_count":2,
      "current_page":1
   }
}
```



### Get historic withdrawals in a specific fiat currency

```GO
buda.GetWithdrawHistory(currency string)
```

```json
{  
   "withdrawals":[  
      {  
         "id":1234567,
         "state":"pending_op_execution",
         "currency":"COP",
         "created_at":"2018-01-29T02:55:04.056Z",
         "withdrawal_data":{  
            "type":"fiat/withdrawal_data",
            "id":1234567,
            "created_at":"2018-01-29T02:55:04.048Z",
            "updated_at":"2018-01-29T02:55:04.048Z",
            "transacted_at":null,
            "statement_ref":null,
            "fiat_account":{  
               "id":1234567,
               "account_number":"000000000",
               "account_type":"Account description",
               "bank_id":99,
               "created_at":"2018-01-22T16:57:32.949Z",
               "currency":"COP",
               "document_number":"000000000",
               "email":"email@email.com",
               "full_name":"Your name",
               "national_number_identifier":null,
               "phone":"+000000000000",
               "updated_at":"2018-01-22T16:57:32.949Z",
               "bank_name":"BankName",
               "pe_cci_number":null
            },
            "source_account":null
         },
         "amount":[  
            "9000.0",
            "COP"
         ],
         "fee":[  
            "0.0",
            "COP"
         ]
      }
   ],
   "meta":{  
      "total_pages":1,
      "total_count":1,
      "current_page":1
   }
}
```


##### CONTRIBUTION ⚠️

FEEL FREE TO OPEN AN ISSUE / PULL REQUEST IF YOU FIND BUGS OR IF YOU WANT TO INCLUDE NEW FEATURES.
If your pull request doesn't break any old code I will happily merge it.

This client works with v2 of the Buda API.
