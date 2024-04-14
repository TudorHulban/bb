package strategies

import "test/ordering"

type IStrategy interface {
	AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error)
}

var _ IStrategy = &StrategyDropSudden{}
