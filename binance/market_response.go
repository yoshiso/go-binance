/*

   market_response.go
       Stores response structs/handlers for API functions in market.go

*/

package binance

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// Result from: GET /api/v1/depth
type OrderBook struct {
	LastUpdateId int64   `json:"lastUpdateId"`
	Bids         []Order `json:"bids"`
	Asks         []Order `json:"asks"`
}

type Order struct {
	Price    decimal.Decimal `json:",string"`
	Quantity decimal.Decimal `json:",string"`
}

// Custom Unmarshal function to handle response data format
func (o *Order) UnmarshalJSON(b []byte) error {
	var s [2]string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	o.Price, err = decimal.NewFromString(s[0])
	if err != nil {
		return err
	}

	o.Quantity, err = decimal.NewFromString(s[1])
	if err != nil {
		return err
	}

	return nil
}

// Result from: GET /api/v1/ticker/24hr
type ChangeStats struct {
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	BidPrice           float64 `json:"bidPrice,string"`
	AskPrice           float64 `json:"askPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstId            int64   `json:"firstId"`
	LastId             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

// Result from: GET /api/v1/aggTrade
type AggTrade struct {
	TradeId      int64   `json:"a"`
	Price        float64 `json:"p,string"`
	Quantity     float64 `json:"q,string"`
	FirstTradeId int64   `json:"f"`
	LastTradeId  int64   `json:"l"`
	Timestamp    int64   `json:"T"`
	Maker        bool    `json:"m"`
	BestMatch    bool    `json:"M"`
}

// Result from: GET /api/v1/allPrices
type TickerPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// Result from: GET /api/v1/allBookTickers
type BookTicker struct {
	Symbol      string  `json:"symbol"`
	BidPrice    float64 `json:"bidPrice,string"`
	BidQuantity float64 `json:"bidQty,string"`
	AskPrice    float64 `json:"askPrice,string"`
	AskQuantity float64 `json:"askQty,string"`
}

// Result from: GET /api/v1/klines
type Kline struct {
	OpenTime         int64           `json:"t"`
	CloseTime        int64           `json:"T"`
	Open             decimal.Decimal `json:"o"`
	High             decimal.Decimal `json:"h"`
	Low              decimal.Decimal `json:"l"`
	Close            decimal.Decimal `json:"c"`
	Volume           decimal.Decimal `json:"v"`
	Trades           int64           `json:"n"`
	Closed           bool            `json:"x"`
	QuoteVolume      decimal.Decimal `json:"q"`
	TakerBaseVolume  decimal.Decimal `json:"V"`
	TakerQuoteVolume decimal.Decimal `json:"Q"`
}

// Result from: GET /api/v3/exchangeInfo

type ExchangeInfo struct {
	ExchangeFilters []string     `json:"ExchangeFilters"`
	RateLimits      []RateLimit  `json:"rateLimits"`
	ServerTime      int64        `json:"serverTime"`
	Symbols         []SymbolInfo `json:"symbols"`
	TimeZone        string       `json:"timezone"`
}

type SymbolInfo struct {
	Symbol             string         `json:"symbol"`
	BaseAsset          string         `json:"baseAsset"`
	QuotePrecision     int64          `json:"quotePrecision"`
	BaseAssetPrecision int64          `json:"baseAssetPrecision"`
	Status             string         `json:"status"`
	OrderTypes         []string       `json:"orderTypes"`
	Filters            []SymbolFilter `json:"filters"`
	QuoteAsset         string         `json:"quoteAsset"`
	IceBergAllowed     bool           `json:"icebergAllowed"`
}

type SymbolFilter struct {
	Type        string  `json:"filterType"`
	MinPrice    float64 `json:"minPrice,string"`
	MaxPrice    float64 `json:"maxPrice,string"`
	TickSize    float64 `json:"tickSize,string"`
	StepSize    float64 `json:"stepSize,string"`
	MinQty      float64 `json:"minQty,string"`
	MaxQty      float64 `json:"maxQty,string"`
	MinNotional float64 `json:"minNotional,string"`
}

type RateLimit struct {
	Limit         int64  `json:"limit"`
	Interval      string `json:"interval"`
	RateLimitType string `json:"rateLimitType"`
}
