package strategies

import (
	"test/ordering"

	"github.com/govalues/decimal"
)

type IStrategy interface {
	SetPrice(price decimal.Decimal) error
	AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error)

	CanPlaceOrder() bool
	IncrementSimultaneousOrders()
	DecrementSimultaneousOrders()
}

var _ IStrategy = &StrategyDropSudden{}
