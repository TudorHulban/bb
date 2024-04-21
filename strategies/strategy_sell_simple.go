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

const nameStrategySellSimple = "Sell Simple"

type StrategySellSimple struct {
	PriceBeforeRaise decimal.Decimal
	PercentRaise     decimal.Decimal
}

type ParamsNewStrategySellSimple struct {
	PriceBeforeRaise decimal.Decimal
	PercentRaise     decimal.Decimal
}

func NewStrategySellSimple(params *ParamsNewStrategySellSimple) (*StrategySellSimple, error) {
	if !params.PercentRaise.IsPos() {
		return nil,
			apperrors.ErrorInvalidInput{
				InputName: "PercentRaise",
			}
	}

	return &StrategySellSimple{
			PriceBeforeRaise: params.PriceBeforeRaise,
			PercentRaise:     params.PercentRaise,
		},
		nil
}

func (s *StrategySellSimple) SetPrice(price decimal.Decimal) error {
	if !price.IsPos() {
		return apperrors.ErrorInvalidInput{
			InputName: "SetPrice",
		}
	}

	s.PriceBeforeRaise = price

	return nil
}

func (s *StrategySellSimple) AddPriceChange(params *ParamsAddPriceChangeSell) (ordering.Action, error) {
	difference, errSubtract := s.PriceBeforeRaise.Sub(params.PriceNow)
	if errSubtract != nil {
		return ordering.DoNothing, errSubtract
	}

	if difference.IsPos() {
		return ordering.DoNothing,
			errors.New("price went down")
	}

	if errPercent := helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: s.PriceBeforeRaise,
			PriceNew: params.PriceNow,
			Delta:    s.PercentRaise,
		},
	); errPercent != nil {
		return ordering.DoNothing,
			errors.New("not enough raise")
	}

	return ordering.Sell, nil
}

func (s StrategySellSimple) String() string {
	return strings.Join(
		[]string{
			"",
			fmt.Sprintf("Strategy '%s'", nameStrategySellSimple),
			fmt.Sprintf("raise percent: %s%%.", s.PercentRaise.String()),
			fmt.Sprintf("current price before raise: %s.", s.PriceBeforeRaise.String()),
			"",
		},
		"\n",
	)
}
