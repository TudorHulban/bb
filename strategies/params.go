package strategies

import "github.com/govalues/decimal"

type ParamsAddPriceChange struct {
	AverageMediumPeriodPrice decimal.Decimal
	PriceNow                 decimal.Decimal

	NoPriceChangesPeriodShort  uint32
	NoPriceChangesPeriodMedium uint32
}
