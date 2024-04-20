package strategies

import (
	"errors"
	"fmt"
	"strings"
	"test/apperrors"
	"test/helpers"
	"test/ordering"

	"github.com/govalues/decimal"
)

const nameStrategy = "Drop Sudden"

type StrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal
}

type ParamsNewStrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal
}

func NewStrategyDropSudden(params ParamsNewStrategyDropSudden) (*StrategyDropSudden, error) {
	if !params.DropPercent.IsPos() {
		return nil,
			apperrors.ErrorInvalidInput{
				InputName: "DropPercent",
			}
	}

	return &StrategyDropSudden{
			PriceBeforeDrop: params.PriceBeforeDrop,
			DropPercent:     params.DropPercent,
		},
		nil
}

func (s *StrategyDropSudden) SetPrice(price decimal.Decimal) error {
	if !price.IsPos() {
		return apperrors.ErrorInvalidInput{
			InputName: "SetPrice",
		}
	}

	s.PriceBeforeDrop = price

	return nil
}

func (s *StrategyDropSudden) AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error) {
	difference, errSubtract := s.PriceBeforeDrop.Sub(params.PriceNow)
	if errSubtract != nil {
		return ordering.DoNothing, errSubtract
	}

	if difference == decimal.Zero {
		return ordering.DoNothing,
			errors.New("no price change")
	}

	if errPercent := helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: s.PriceBeforeDrop,
			PriceNew: params.PriceNow,
			Delta:    s.DropPercent,
		},
	); errPercent != nil {
		return ordering.DoNothing,
			errors.New("no sudden drop")
	}

	return ordering.Buy, nil
}

func (s StrategyDropSudden) String() string {
	return strings.Join(
		[]string{
			"",
			fmt.Sprintf("Strategy %s", nameStrategy),
			fmt.Sprintf("drop percent: %s%%.", s.DropPercent.String()),
			fmt.Sprintf("current price before drop: %s.", s.PriceBeforeDrop.String()),
			"",
		},
		"\n",
	)
}
