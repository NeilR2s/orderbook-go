package main

import (
	"fmt"
	"slices"
	"sort"
)

type OrderType int

const (
	GoodTillCancel = iota
	FillandKill
)

type Side int

const (
	Unknown = iota
	Buy
	Sell
)

type Price int32
type Quantity uint32
type OrderId uint64

type LevelInfo struct {
	Price    Price
	Quantity Quantity
}

type LevelInfos []LevelInfo

type Order struct {
	orderType         OrderType
	orderId           OrderId
	side              Side
	price             Price
	initialQuantity   Quantity
	remainingQuantity Quantity
}

type OrderBook struct {
	bids []*Order
	asks []*Order
}

func NewOrder(orderType OrderType, orderId OrderId, side Side, price Price, quantity Quantity) *Order {
	return &Order{
		orderType:         orderType,
		orderId:           orderId,
		side:              side,
		price:             price,
		initialQuantity:   quantity,
		remainingQuantity: quantity,
	}
}

func (o *Order) OrderId() OrderId {
	return o.orderId
}

func (o *Order) Side() Side {
	return o.side
}

func (o *Order) Price() Price {
	return o.price
}

func (o *Order) OrderType() OrderType {
	return o.orderType
}

func (o *Order) InitialQuantity() Quantity {
	return o.initialQuantity
}

func (o *Order) RemainingQuantity() Quantity {
	return o.remainingQuantity
}

func (o *Order) FilledQuantity() Quantity {
	return o.InitialQuantity() - o.RemainingQuantity()
}

func (s Side) String() string {
	orderSide := ""
	switch s {
	case Buy:
		orderSide = "Buy"
	case Sell:
		orderSide = "Sell"
	default:
		return "Unknown"
	}
	return orderSide
}

func (o *Order) String() string {

	return fmt.Sprintf("Order[%d] %v %d @ %d", o.orderId, o.side, o.remainingQuantity, o.price)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		bids: []*Order{},
		asks: []*Order{},
	}
}

func (ob *OrderBook) AddOrder(o *Order) {
	if o.side == Buy {
		ob.bids = append(ob.bids, o)
		sort.Slice(ob.bids, func(i, j int) bool {
			return ob.bids[i].price > ob.bids[j].price
		})
	}
	if o.side == Sell {
		ob.asks = append(ob.asks, o)
		sort.Slice(ob.asks, func(i, j int) bool {
			return ob.asks[i].price < ob.asks[j].price
		})
	}
	ob.MatchOrder()
}

func (ob *OrderBook) MatchOrder() {
	for {
		if len(ob.bids) == 0 || len(ob.asks) == 0 {
			break
		} else if ob.bids[0].Price() < ob.asks[0].Price() {
			break
		} else {
			topBid := ob.bids[0]
			topAsk := ob.asks[0]
			filledQuantity := min(topBid.remainingQuantity, topAsk.remainingQuantity)
			topBid.remainingQuantity -= filledQuantity
			topAsk.remainingQuantity -= filledQuantity

			if topBid.remainingQuantity == 0 {
				ob.bids = slices.Delete(ob.bids, 0, 1)
			}
			if topAsk.remainingQuantity == 0 {
				ob.asks = slices.Delete(ob.asks, 0, 1)
			}
		}
	}
}

func main() {
	testBid := NewOrder(GoodTillCancel, 1231230, Buy, 100, 10)
	testAsk := NewOrder(GoodTillCancel, 998139, Sell, 99, 10)

	testOrderBook := NewOrderBook()
	testOrderBook.AddOrder(testBid)
	testOrderBook.AddOrder(testAsk)

	fmt.Printf("%v", testBid.String())
	fmt.Printf("%v", testAsk.String())
}
