package main

import (
	"fmt"

	"tulip"
)

func main() {
	buda := tulip.CreateClient("", "")
	somestruct := buda.GetMarkets()
	fmt.Println(somestruct.markets)
}
