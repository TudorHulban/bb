package strategies

import (
	"test/ordering"

	"github.com/govalues/decimal"
)

type IStrategy interface {
	IsReady() bool
	SetPrice(price decimal.Decimal)
	AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error)
}

var _ IStrategy = &StrategyDropSudden{}
