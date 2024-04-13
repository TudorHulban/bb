package ordering

import "log"

type OrderingLogOnly struct{}

func NewOrderingLogOnly() *OrderingLogOnly {
	return &OrderingLogOnly{}
}

func (o OrderingLogOnly) Buy() error {
	log.Println("buy")

	return nil
}

func (o OrderingLogOnly) Sell() error {
	log.Println("sell")

	return nil
}
