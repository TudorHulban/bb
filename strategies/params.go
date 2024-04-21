package strategies

import "github.com/govalues/decimal"

type ParamsAddPriceChangeBuy struct {
	AverageMediumPeriodPrice decimal.Decimal
	PriceNow                 decimal.Decimal

	NoPriceChangesPeriodShort  uint32
	NoPriceChangesPeriodMedium uint32
}

type ParamsAddPriceChangeSell struct {
	PriceBuy decimal.Decimal
	PriceNow decimal.Decimal
}
