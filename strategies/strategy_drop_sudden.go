package strategies

import (
	"test/ordering"

	"github.com/govalues/decimal"
)

type StrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal
}

func (s *StrategyDropSudden) AddPriceChange(params *ParamsAddPriceChange) (*ordering.Order, error) {
	if s.PriceBeforeDrop == decimal.Zero {
		s.PriceBeforeDrop = params.PriceNow
	}

	return nil, nil
}
