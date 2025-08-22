package main

import (
	"container/list"
	"time"

	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type Side int

const (
	Buy Side = iota
	Sell
)

type Order struct {
	id       string
	side     Side
	price    decimal.Decimal
	quantity decimal.Decimal
	time     time.Time
}

type OrderQueue struct {
	orders   *list.List
	price    decimal.Decimal
	quantity decimal.Decimal
}

type Trade struct {
	price    decimal.Decimal
	quantity decimal.Decimal
	makerID  string
	takerID  string
}

func rbtComparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

func NewTree() *rbtx.RedBlackTreeExtended {
	return &rbtx.RedBlackTreeExtended{
		Tree: rbt.NewWith(rbtComparator),
	}
}
