package helpers

import (
	"errors"

	"github.com/govalues/decimal"
)

func PriceChangeByPercent(priceOld, priceNew, delta decimal.Decimal) error {
	difference, errSubtract := priceNew.SubAbs(priceOld)
	if errSubtract != nil {
		return errSubtract
	}

	division, errDivision := difference.Quo(priceOld)
	if errDivision != nil {
		return errDivision
	}

	multiply100, errMultiply00 := division.FMA(decimal.Hundred, decimal.Zero)
	if errMultiply00 != nil {
		return errMultiply00
	}

	subtract, errSubtract := multiply100.Sub(delta)
	if errSubtract != nil {
		return errSubtract
	}

	if subtract.Sign() == -1 {
		return errors.New("no price change")
	}

	return nil
}
