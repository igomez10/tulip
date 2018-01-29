# Tulip



An open source API for surBTC (Buda) written in Go. An alternative for current apis written in Java or Python.

You can check the current surBTC (Buda) API documentation at http://api.surbtc.com

NOTE: You need a developer account in order to make full use of this client.
Altough ya can make unauthenticated calls (http://api.surbtc.com/#llamadas-p-blicas)

##### DISCLAIMER

BEWARE THIS IS OPEN SOURCE IN EARLY STAGES, IM DOING MY BEST, FEEL FREE TO OPEN AN ISSUE / PULL REQUEST IF YOU FIND BUGS OR IF YOU WANT TO INCLUDE NEW FEATURES.
If your pull request doesnt break anything I will happily merge it.

This Client works with v2 of the surbtc API. It might need to change in the coming weeks with the change to Buda as their new name.
I assume all endpoints will need to change to something similar to buda.com/api/v2/something .
In the meantime, I'm testing it  with the latest API


#### TODO Endpoints

  - Do Deposits and Withdrawls

## Get it working

From you terminal:
```sh
$ go get github.com/igomez10/tulip
```


From your file:
``` GO
somefile.go

import (
    tulip "github.com/igomez10/tulip"
    "fmt"
)

buda := tulip.CreateClient(<youPublicAPIKEY> , <yourSecretPrivateKey>)
results := buda.GetMarkets()
fmt.Println(results)

```



# Documentation

## Public Calls

### Get Ticker

```Go
buda.GetTicker("btc-clp")
```

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
```Go
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
```Go
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
        ],       
      ]
    }
}
```
