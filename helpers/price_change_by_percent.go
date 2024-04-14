package helpers

import (
	"errors"
	"fmt"
	"test/apperrors"

	"github.com/govalues/decimal"
)

type ParamsPriceChangeByPercent struct {
	PriceOld decimal.Decimal
	PriceNew decimal.Decimal
	Delta    decimal.Decimal
}

func (p ParamsPriceChangeByPercent) String() string {
	return fmt.Sprintf(
		"PriceOld: %s,\nPriceNew: %s,\nDelta: %s",
		p.PriceOld.String(),
		p.PriceNew.String(),
		p.Delta.String(),
	)
}

// 100 * (newPrice - currentPrice) / currentPrice < delta
func PriceChangeByPercent(params *ParamsPriceChangeByPercent) error {
	if params.PriceOld == decimal.Zero && params.PriceNew == decimal.Zero {
		return apperrors.ErrorInvalidInputs{
			InputsName: []string{
				"PriceOld",
				"PriceNew",
			},
		}
	}

	if params.PriceOld == decimal.Zero {
		return nil
	}

	difference, errSubtract := params.PriceNew.SubAbs(params.PriceOld)
	if errSubtract != nil {
		return errSubtract
	}

	division, errDivision := difference.Quo(params.PriceOld)
	if errDivision != nil {
		return errDivision
	}

	multiply100, errMultiply00 := division.FMA(decimal.Hundred, decimal.Zero)
	if errMultiply00 != nil {
		return errMultiply00
	}

	subtract, errSubtract := multiply100.Sub(params.Delta)
	if errSubtract != nil {
		return errSubtract
	}

	if subtract.Sign() == -1 {
		return errors.New("no price change")
	}

	return nil
}
