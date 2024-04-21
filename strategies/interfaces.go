package strategies

import (
	"test/ordering"

	"github.com/govalues/decimal"
)

type IStrategyBuy interface {
	SetPrice(price decimal.Decimal) error
	AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error)

	CanPlaceOrder() bool
	IncrementSimultaneousOrders()
	DecrementSimultaneousOrders()
}

var _ IStrategyBuy = &StrategyDropSudden{}

type IStrategySell interface {
	SetPrice(price decimal.Decimal) error
	AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error)
}
