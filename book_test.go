package main

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestSideTree(t *testing.T) {
	ob := NewOrderBook()

	if res := ob.currentTree(Buy); res != ob.Bids {
		t.Errorf("want Bids, got %+v", res)
	}

	if res := ob.currentTree(Sell); res != ob.Asks {
		t.Errorf("want Asks, got %+v", res)
	}
}

func TestOppositeTree(t *testing.T) {
	ob := NewOrderBook()

	if res := ob.oppositeTree(Buy); res != ob.Asks {
		t.Errorf("expected Asks, got %+v", res)
	}

	if res := ob.oppositeTree(Sell); res != ob.Bids {
		t.Errorf("expected Bids, got %+v", res)
	}
}

func TestBestBuy(t *testing.T) {
	tests := []struct {
		name      string
		orders    []*Order
		wantLevel Level
		wantOK    bool
	}{
		{
			name:      "empty order book",
			orders:    []*Order{},
			wantLevel: Level{},
			wantOK:    false,
		},
		{
			name: "one order",
			orders: []*Order{
				{
					side:     Buy,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(1),
				},
			},
			wantLevel: Level{
				Price:    decimal.NewFromInt(10),
				Quantity: decimal.NewFromInt(1),
			},
			wantOK: true,
		},
		{
			name: "multiple orders",
			orders: []*Order{
				{
					side:     Buy,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(1),
				},
				{
					side:     Buy,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(2),
				},
				{
					side:     Buy,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(3),
				},
			},
			wantLevel: Level{
				Price:    decimal.NewFromInt(10),
				Quantity: decimal.NewFromInt(6),
			},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := NewOrderBook()

			for _, order := range tt.orders {
				ob.addOrder(order)
			}

			gotLevel, gotOK := ob.BestBuy()

			priceCheck := !gotLevel.Price.Equal(tt.wantLevel.Price)
			quantityCheck := !gotLevel.Quantity.Equal(tt.wantLevel.Quantity)
			okCheck := gotOK != tt.wantOK

			if priceCheck || quantityCheck || okCheck {
				t.Errorf("want level=%+v, ok=%v; got level=%+v, ok=%v",
					tt.wantLevel, tt.wantOK, gotLevel, gotOK)
			}
		})
	}
}

func TestBestSell(t *testing.T) {
	tests := []struct {
		name      string
		orders    []*Order
		wantLevel Level
		wantOK    bool
	}{
		{
			name:      "empty order book",
			orders:    []*Order{},
			wantLevel: Level{},
			wantOK:    false,
		},
		{
			name: "one order",
			orders: []*Order{
				{
					side:     Sell,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(1),
				},
			},
			wantLevel: Level{
				Price:    decimal.NewFromInt(10),
				Quantity: decimal.NewFromInt(1),
			},
			wantOK: true,
		},
		{
			name: "multiple orders",
			orders: []*Order{
				{
					side:     Sell,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(1),
				},
				{
					side:     Sell,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(2),
				},
				{
					side:     Sell,
					price:    decimal.NewFromInt(10),
					quantity: decimal.NewFromInt(3),
				},
			},
			wantLevel: Level{
				Price:    decimal.NewFromInt(10),
				Quantity: decimal.NewFromInt(6),
			},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := NewOrderBook()

			for _, order := range tt.orders {
				ob.addOrder(order)
			}

			gotLevel, gotOK := ob.BestSell()

			priceCheck := !gotLevel.Price.Equal(tt.wantLevel.Price)
			quantityCheck := !gotLevel.Quantity.Equal(tt.wantLevel.Quantity)
			okCheck := gotOK != tt.wantOK

			if priceCheck || quantityCheck || okCheck {
				t.Errorf("BestSell() = level=%+v, ok=%v; want level=%+v, ok=%v",
					gotLevel, gotOK, tt.wantLevel, tt.wantOK)
			}
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	ob := NewOrderBook()

	order := &Order{
		id:       "1",
		side:     Buy,
		price:    decimal.NewFromInt(10),
		quantity: decimal.NewFromInt(1),
	}

	// First Buy order placed.
	ob.PlaceOrder(order.side, order.price, order.quantity, order.id)

	if ob.Bids.Size() != 1 {
		t.Errorf("Want 1 bid, got %d", ob.Bids.Size())
	}
	if ob.Asks.Size() != 0 {
		t.Errorf("Want 0 asks, got %d", ob.Asks.Size())
	}

	// Second Buy order placed at same price.
	order2 := &Order{
		id:       "2",
		side:     Buy,
		price:    decimal.NewFromInt(10),
		quantity: decimal.NewFromInt(2),
	}

	ob.PlaceOrder(order2.side, order2.price, order2.quantity, order2.id)

	if ob.Bids.Size() != 1 {
		t.Errorf("Want 1 bid, got %d", ob.Bids.Size())
	}
	if ob.Asks.Size() != 0 {
		t.Errorf("Want 0 asks, got %d", ob.Asks.Size())
	}

	// Third Buy order placed at different price.
	order3 := &Order{
		id:       "3",
		side:     Buy,
		price:    decimal.NewFromInt(11),
		quantity: decimal.NewFromInt(3),
	}

	ob.PlaceOrder(order3.side, order3.price, order3.quantity, order3.id)

	if ob.Bids.Size() != 2 {
		t.Errorf("Want 2 bids, got %d", ob.Bids.Size())
	}
	if ob.Asks.Size() != 0 {
		t.Errorf("Want 0 asks, got %d", ob.Asks.Size())
	}

	// Place Sell order that partially matches.
	sellOrder := &Order{
		id:       "4",
		side:     Sell,
		price:    decimal.NewFromInt(10),
		quantity: decimal.NewFromInt(4),
	}

	trades := ob.PlaceOrder(sellOrder.side, sellOrder.price, sellOrder.quantity, sellOrder.id)

	if ob.Bids.Size() != 1 {
		t.Errorf("Want 1 bid, got %d", ob.Bids.Size())
	}
	if ob.Asks.Size() != 0 {
		t.Errorf("Want 0 asks, got %d", ob.Asks.Size())
	}

	if len(trades) != 3 {
		t.Errorf("Want 3 trades, got %d", len(trades))
	}

	want := []Trade{
		{makerID: "1", takerID: "4", price: decimal.NewFromInt(10), quantity: decimal.NewFromInt(1)},
		{makerID: "2", takerID: "4", price: decimal.NewFromInt(10), quantity: decimal.NewFromInt(2)},
		{makerID: "3", takerID: "4", price: decimal.NewFromInt(11), quantity: decimal.NewFromInt(1)},
	}

	if !reflect.DeepEqual(trades, want) {
		t.Errorf("trades mismatch, got %+v", trades)
	}
}
