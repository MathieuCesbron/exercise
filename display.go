package main

import (
	"fmt"
)

func PrintOrderBook(ob *OrderBook) {
	fmt.Println("Order Book:")
	fmt.Println("  Asks:")
	it := ob.Asks.Iterator()
	for it.Begin(); it.Next(); {
		queue := it.Value().(*OrderQueue)
		fmt.Printf("    Price: %v, Qty: %v\n", queue.price, queue.quantity)
		for e := queue.orders.Front(); e != nil; e = e.Next() {
			order := e.Value.(*Order)
			fmt.Printf("      OrderID: %s, Qty: %v, Time: %v\n", order.id, order.quantity, order.time)
		}
	}
	fmt.Println("  Bids:")
	it = ob.Bids.Iterator()
	for it.End(); it.Prev(); {
		queue := it.Value().(*OrderQueue)
		fmt.Printf("    Price: %v, Qty: %v\n", queue.price, queue.quantity)
		for e := queue.orders.Front(); e != nil; e = e.Next() {
			order := e.Value.(*Order)
			fmt.Printf("      OrderID: %s, Qty: %v, Time: %v\n", order.id, order.quantity, order.time)
		}
	}
}

func PrintTrades(trades []Trade) {
	fmt.Println("Trades:")
	for _, t := range trades {
		fmt.Printf("  Price: %v, Qty: %v, MakerID: %s, TakerID: %s\n", t.price, t.quantity, t.makerID, t.takerID)
	}
}
