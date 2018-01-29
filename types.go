package tulip

type Market struct {
	id                   string
	name                 string
	base_currency        string
	minimum_order_amount string
}

type Order_book struct {
	asks [][]string
	bids [][]string
}

type Trade struct {
	market_id      string
	timestamp      string
	last_timestamp string
	entries        [][]string
}

type Response struct {
	Error  []string
	Result interface{}
}
