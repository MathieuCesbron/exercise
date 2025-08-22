package main

import (
	"container/list"
	"time"

	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	// key=price, value=OrderQueue.
	Asks *rbtx.RedBlackTreeExtended
	Bids *rbtx.RedBlackTreeExtended
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: NewTree(),
		Bids: NewTree(),
	}
}

func (ob *OrderBook) currentTree(side Side) *rbtx.RedBlackTreeExtended {
	if side == Buy {
		return ob.Bids
	}

	return ob.Asks
}

func (ob *OrderBook) oppositeTree(side Side) *rbtx.RedBlackTreeExtended {
	if side == Buy {
		return ob.Asks
	}

	return ob.Bids
}

type Level struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}

// BestBuy returns the highest bid and its total quantity.
func (ob *OrderBook) BestBuy() (Level, bool) {
	if v, ok := ob.Bids.GetMax(); ok {
		oq := v.(*OrderQueue)

		return Level{Price: oq.price, Quantity: oq.quantity}, true
	}

	return Level{}, false
}

// BestSell returns the lowest ask and its total quantity.
func (ob *OrderBook) BestSell() (Level, bool) {
	if v, ok := ob.Asks.GetMin(); ok {
		oq := v.(*OrderQueue)

		return Level{Price: oq.price, Quantity: oq.quantity}, true
	}

	return Level{}, false
}

func cross(p1, p2 decimal.Decimal, side Side) bool {
	if side == Buy {
		return p1.LessThanOrEqual(p2)
	}

	return p1.GreaterThanOrEqual(p2)
}

func (oq *OrderQueue) matchOrder(quantity decimal.Decimal, id string) ([]Trade, decimal.Decimal) {
	trades := []Trade{}

	for e := oq.orders.Front(); e != nil && !quantity.IsZero(); {
		current := e
		e = e.Next()

		order := current.Value.(*Order)
		filled := decimal.Min(quantity, order.quantity)
		trades = append(trades, Trade{
			price:    order.price,
			quantity: filled,
			makerID:  order.id,
			takerID:  id,
		})

		oq.quantity = oq.quantity.Sub(filled)
		quantity = quantity.Sub(filled)
		order.quantity = order.quantity.Sub(filled)

		// Remove the current order if no quantity left.
		if order.quantity.IsZero() {
			oq.orders.Remove(current)
		}
	}

	return trades, quantity
}

func (ob *OrderBook) addOrder(order *Order) {
	tree := ob.currentTree(order.side)

	if v, ok := tree.Get(order.price); ok {
		oq := v.(*OrderQueue)
		oq.orders.PushBack(order)
		oq.quantity = oq.quantity.Add(order.quantity)
	} else {
		orders := list.New()
		orders.PushFront(order)
		oq := &OrderQueue{
			price:    order.price,
			quantity: order.quantity,
			orders:   orders,
		}
		tree.Put(order.price, oq)
	}
}

// PlaceOrder matches an incoming order against the orderbook and
// adds any remaining quantity as a new order.
func (ob *OrderBook) PlaceOrder(side Side, price decimal.Decimal, quantity decimal.Decimal, id string) []Trade {
	trades := []Trade{}

	// Match order with orderbook.
	it := ob.oppositeTree(side).Iterator()
	for it.Begin(); it.Next(); {
		priceLevel := it.Key().(decimal.Decimal)
		if !cross(priceLevel, price, side) {
			break
		}

		var tradesLevel []Trade
		oq := it.Value().(*OrderQueue)
		tradesLevel, quantity = oq.matchOrder(quantity, id)
		trades = append(trades, tradesLevel...)

		if oq.quantity.IsZero() {
			it.Prev()
			ob.oppositeTree(side).Remove(priceLevel)
		}
	}

	// Add remaining quantity to the orderbook.
	if quantity.GreaterThan(decimal.Zero) {
		order := &Order{
			id:       id,
			side:     side,
			price:    price,
			quantity: quantity,
			time:     time.Now().UTC(),
		}
		ob.addOrder(order)
	}

	return trades
}
