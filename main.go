package main

import (
	"github.com/shopspring/decimal"
)

func main() {
	ob := NewOrderBook()

	// PlaceOrder to place orders to the orderbook.
	ob.PlaceOrder(Buy, decimal.NewFromInt(10), decimal.NewFromFloat(1), "1")
	ob.PlaceOrder(Sell, decimal.NewFromInt(11), decimal.NewFromFloat(1), "2")

	trades := ob.PlaceOrder(Buy, decimal.NewFromInt(11), decimal.NewFromFloat(0.5), "3")

	// Display trades and orderbook state.
	PrintTrades(trades)
	PrintOrderBook(ob)
}
