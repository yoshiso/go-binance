package main

import (
	"fmt"

	"github.com/yoshiso/go-binance/binance"
)

func main() {
	cli := binance.New("", "")
	resp, err := cli.GetKlines(binance.KlineQuery{
		Symbol:   "ETHBTC",
		Interval: "m1",
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
