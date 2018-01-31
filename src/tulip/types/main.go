package types

import "time"

// Parsed with https://mholt.github.io/json-to-go/ , huge thanks to Matt Holt for this tool

func Describe(i interface{}) interface{} {
	return i
}

type LimitOrder struct {
	OrderType string  `json:"type"`
	PriceType string  `json:"price_type"`
	Limit     float64 `json:"limit"`
	Amount    float64 `json:"amount"`
}

type MarketOrder struct {
	OrderType string  `json:"type"`
	PriceType string  `json:"price_type"`
	Amount    float64 `json:"amount"`
}

type Market struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	BaseCurrency       string   `json:"base_currency"`
	QuoteCurrency      string   `json:"quote_currency"`
	MinimumOrderAmount []string `json:"minimum_order_amount"`
}

type MarketResponse struct {
	Market Market
}

type MarketsResponse struct {
	Markets []Market
}

type OrderBook struct {
	OrderBook struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	} `json:"order_book"`
}

type TradesResponse struct {
	Trades struct {
		MarketID      string          `json:"market_id"`
		Timestamp     interface{}     `json:"timestamp"`
		LastTimestamp string          `json:"last_timestamp"`
		Entries       [][]interface{} `json:"entries"`
	} `json:"trades"`
}

type BalancesResponse struct {
	Balances []Balance
}

type BalanceResponse struct {
	Balance Balance
}

type Balance struct {
	ID                    string   `json:"id"`
	Amount                []string `json:"amount"`
	AvailableAmount       []string `json:"available_amount"`
	FrozenAmount          []string `json:"frozen_amount"`
	PendingWithdrawAmount []string `json:"pending_withdraw_amount"`
	AccountID             int      `json:"account_id"`
}

type Order struct {
	ID             int       `json:"id"`
	MarketID       string    `json:"market_id"`
	AccountID      int       `json:"account_id"`
	Type           string    `json:"type"`
	State          string    `json:"state"`
	CreatedAt      time.Time `json:"created_at"`
	FeeCurrency    string    `json:"fee_currency"`
	PriceType      string    `json:"price_type"`
	Limit          []string  `json:"limit"`
	Amount         []string  `json:"amount"`
	OriginalAmount []string  `json:"original_amount"`
	TradedAmount   []string  `json:"traded_amount"`
	TotalExchanged []string  `json:"total_exchanged"`
	PaidFee        []string  `json:"paid_fee"`
}

type MyOrdersResponse struct {
	Orders []Order
	Meta   struct {
		TotalPages  int `json:"total_pages"`
		TotalCount  int `json:"total_count"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
}

type OrderResponse struct {
	Order Order
}

type HistoricDespositsResponse struct {
	Deposits []struct {
		ID          int       `json:"id"`
		State       string    `json:"state"`
		Currency    string    `json:"currency"`
		CreatedAt   time.Time `json:"created_at"`
		DepositData struct {
			Type      string      `json:"type"`
			CreatedAt time.Time   `json:"created_at"`
			UpdatedAt time.Time   `json:"updated_at"`
			UploadURL interface{} `json:"upload_url"`
		} `json:"deposit_data"`
		Amount []string `json:"amount"`
		Fee    []string `json:"fee"`
	} `json:"deposits"`
	Meta struct {
		TotalPages  int `json:"total_pages"`
		TotalCount  int `json:"total_count"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
}

type HistoricWithdrawResponse struct {
	Withdrawals []struct {
		ID             int       `json:"id"`
		State          string    `json:"state"`
		Currency       string    `json:"currency"`
		CreatedAt      time.Time `json:"created_at"`
		WithdrawalData struct {
			Type         string      `json:"type"`
			ID           int         `json:"id"`
			CreatedAt    time.Time   `json:"created_at"`
			UpdatedAt    time.Time   `json:"updated_at"`
			TransactedAt interface{} `json:"transacted_at"`
			StatementRef interface{} `json:"statement_ref"`
			FiatAccount  struct {
				ID                       int         `json:"id"`
				AccountNumber            string      `json:"account_number"`
				AccountType              string      `json:"account_type"`
				BankID                   int         `json:"bank_id"`
				CreatedAt                time.Time   `json:"created_at"`
				Currency                 string      `json:"currency"`
				DocumentNumber           string      `json:"document_number"`
				Email                    string      `json:"email"`
				FullName                 string      `json:"full_name"`
				NationalNumberIdentifier interface{} `json:"national_number_identifier"`
				Phone                    string      `json:"phone"`
				UpdatedAt                time.Time   `json:"updated_at"`
				BankName                 string      `json:"bank_name"`
				PeCciNumber              interface{} `json:"pe_cci_number"`
			} `json:"fiat_account"`
			SourceAccount interface{} `json:"source_account"`
		} `json:"withdrawal_data"`
		Amount []string `json:"amount"`
		Fee    []string `json:"fee"`
	} `json:"withdrawals"`
	Meta struct {
		TotalPages  int `json:"total_pages"`
		TotalCount  int `json:"total_count"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
}
