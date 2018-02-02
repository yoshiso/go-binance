/*
   depth.go
       Connects to the Binance WebSocket and maintains
       local depth cache.

*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yoshiso/go-binance/binance"

	"github.com/gorilla/websocket"
)

const (
	MaxDepth = 100 // Size of order book
	MaxQueue = 100 // Size of message queue
)

// Message received from websocket
type State struct {
	EventType string          `json:"e"`
	EventTime int64           `json:"E"`
	Symbol    string          `json:"s"`
	UpdateId  int64           `json:"u"`
	BidDelta  []binance.Order `json:"b"`
	AskDelta  []binance.Order `json:"a"`
}

// Orderbook structure
type OrderBook struct {
	Bids     map[decimal.Decimal]decimal.Decimal // Map of all bids, key->price, value->quantity
	BidMutex sync.Mutex                          // Threadsafe

	Asks     map[decimal.Decimal]decimal.Decimal // Map of all asks, key->price, value->quantity
	AskMutex sync.Mutex                          // Threadsafe

	Updates chan State // Channel of all state updates
}

// Process all incoming bids
func (o *OrderBook) ProcessBids(bids []binance.Order) {
	for _, bid := range bids {
		o.BidMutex.Lock()
		if bid.Quantity.Equal(decimal.Zero) {
			delete(o.Bids, bid.Price)
		} else {
			o.Bids[bid.Price] = bid.Quantity
		}
		o.BidMutex.Unlock()
	}
}

// Process all incoming asks
func (o *OrderBook) ProcessAsks(asks []binance.Order) {
	for _, ask := range asks {
		o.AskMutex.Lock()
		if ask.Quantity.Equal(decimal.Zero) {
			delete(o.Asks, ask.Price)
		} else {
			o.Asks[ask.Price] = ask.Quantity
		}

		o.AskMutex.Unlock()
	}
}

// Hands off incoming messages to processing functions
func (o *OrderBook) Maintainer() {
	for {
		select {
		case job := <-o.Updates:
			if len(job.BidDelta) > 0 {
				go o.ProcessBids(job.BidDelta)
			}

			if len(job.AskDelta) > 0 {
				go o.ProcessAsks(job.AskDelta)
			}
		}
	}
}

func getKeys(mapItem map[decimal.Decimal]decimal.Decimal) []decimal.Decimal {
	keys := make([]decimal.Decimal, len(mapItem))
	i := 0

	for k := range mapItem {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].GreaterThan(keys[j])
	})

	return keys
}

func (o *OrderBook) Viewer() {
	t := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-t.C:
			o.AskMutex.Lock()
			o.BidMutex.Lock()

			akeys := getKeys(o.Asks)
			bkeys := getKeys(o.Bids)
			fmt.Println("--------")

			for _, key := range akeys {
				fmt.Println("ask", key, o.Asks[key])
			}
			for _, key := range bkeys {
				fmt.Println("bid", key, o.Bids[key])
			}

			o.AskMutex.Unlock()
			o.BidMutex.Unlock()
		}
	}
}

func main() {

	symbol := "ethbtc"
	address := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth", symbol)

	// Connect to websocket
	var wsDialer websocket.Dialer
	wsConn, _, err := wsDialer.Dial(address, nil)
	if err != nil {
		panic(err)
	}
	defer wsConn.Close()
	log.Println("Dialed:", address)

	// Set up Order Book
	ob := OrderBook{}
	ob.Bids = make(map[decimal.Decimal]decimal.Decimal, MaxDepth)
	ob.Asks = make(map[decimal.Decimal]decimal.Decimal, MaxDepth)
	ob.Updates = make(chan State, 500)

	// Get initial state of orderbook from rest api
	client := binance.New("", "")
	query := binance.OrderBookQuery{
		Symbol: strings.ToUpper(symbol),
	}
	orderBook, err := client.GetOrderBook(query)
	if err != nil {
		panic(err)
	}
	ob.ProcessBids(orderBook.Bids)
	ob.ProcessAsks(orderBook.Asks)

	// Start maintaining order book
	go ob.Maintainer()

	go ob.Viewer()

	// Read & Process Messages from wss stream
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("[ERROR] ReadMessage:", err)
		}

		msg := State{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("[ERROR] Parsing:", err)
			continue
		}

		ob.Updates <- msg
	}

}
