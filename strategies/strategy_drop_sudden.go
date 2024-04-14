package strategies

import (
	"errors"
	"test/helpers"
	"test/ordering"

	"github.com/govalues/decimal"
)

type StrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal

	isReady bool
}

type ParamsNewStrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal
}

func NewStrategyDropSudden(params ParamsNewStrategyDropSudden) *StrategyDropSudden {
	// TODO: add input validation

	s := StrategyDropSudden{
		PriceBeforeDrop: params.PriceBeforeDrop,
		DropPercent:     params.DropPercent,
	}

	if params.PriceBeforeDrop.IsPos() {
		s.isReady = true
	}

	return &s
}

func (s *StrategyDropSudden) IsReady() bool {
	return s.isReady
}

func (s *StrategyDropSudden) SetPrice(price decimal.Decimal) {
	s.PriceBeforeDrop = price
}

func (s *StrategyDropSudden) AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error) {
	difference, errSubtract := s.PriceBeforeDrop.Sub(params.PriceNow)
	if errSubtract != nil {
		return ordering.DoNothing, errSubtract
	}

	if difference == decimal.Zero {
		return ordering.DoNothing, errors.New("no price change")
	}

	if errPercent := helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: s.PriceBeforeDrop,
			PriceNew: params.PriceNow,
			Delta:    s.DropPercent,
		},
	); errPercent != nil {
		return ordering.DoNothing, errors.New("no sudden drop")
	}

	return ordering.Buy, nil
}
