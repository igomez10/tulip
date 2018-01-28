
# 🌷Tulip🌷



An open source API for surBTC (Buda) written in Go. An alternative for current apis written in Java or Python.

You can check the current surBTC (Buda) API documentation at http://api.surbtc.com

DISCLAIMER

BEWARE THIS IS OPEN SOURCE IN EARLY STAGES, IM DOING MY BEST, FEEL FREE TO OPEN AN ISSUE / PULL REQUEST IF YOU FIND BUGS OR IF YOU WANT TO INCLUDE NEW FEATURES.
If your pull request doesnt break anything I will happily merge it.

This Client works with v2 of the surbtc API. It might need to change in the coming weeks with the change to Buda as their new name.
I assume all endpoints will need to change to something similar to buda.com/api/v2/something . 
In the meantime, I'm it  with the latest API


# TODO Endpoints

  - Do Deposits and Withdrawls

# Get it working

```sh
$ go get github.com/igomez10/tulip
```

Storing your keys in env variables ensures you will never accidentally upload your credentials to Github/Lab.
[Optional] You can change the env variables that tulip uses by modifying the code in tulip/main.go 
``` GO
	import (
	tulip "github.com/igomez10/tulip"
	"fmt"
	)
	
	buda := tulip.CreateClient(<youPublicAPIKEY> , <yourSecretPrivateKey>)
	results := buda.GetMarkets()
	fmt.Println(results)
	
```

