package coin

import (
	"sync/atomic"
	"test/apperrors"
	"test/configuration"
	"test/ordering"
	"test/strategies"
	"test/timeperiod"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/govalues/decimal"
)

type Coin struct {
	ordering       ordering.IOrdering
	strategiesBuy  []strategies.IStrategyBuy
	strategiesSell []strategies.IStrategySell

	code configuration.CoinCODE

	periodShort  *timeperiod.TimePeriod
	periodMedium *timeperiod.TimePeriod

	defaultQuantityBuy uint32
	currentQuantity    atomic.Uint32
}

type ParamsNewCoin struct {
	Ordering ordering.IOrdering
	Code     configuration.CoinCODE

	MinimumPriceChangesShortPeriod  uint32 `valid:"required"`
	MinimumPriceChangesMediumPeriod uint32 `valid:"required"`

	DefaultQuantityBuy uint32 `valid:"required"`

	MinimumDurationTimeframeShort  time.Duration
	MinimumDurationTimeframeMedium time.Duration

	percentDeltaIsPriceChangeShort float64
}

// NewCoin - if minimum duration is zero, no deletion would occur for old element.
func NewCoin(params *ParamsNewCoin, options ...OptionCoin) (*Coin, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errVa,
			}
	}

	periodShort, errCrShort := timeperiod.NewTimePeriod(
		&timeperiod.ParamsNewTimePeriod{
			Name: timeperiod.NamePeriodShort,

			MinimumPriceChanges:       params.MinimumPriceChangesShortPeriod,
			MinimumDurationTimeframe:  params.MinimumDurationTimeframeShort,
			PercentDeltaIsPriceChange: params.percentDeltaIsPriceChangeShort,
		},
	)
	if errCrShort != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errCrShort,
			}
	}

	periodMedium, errCrMedium := timeperiod.NewTimePeriod(
		&timeperiod.ParamsNewTimePeriod{
			Name: timeperiod.NamePeriodMedium,

			MinimumPriceChanges:      params.MinimumPriceChangesMediumPeriod,
			MinimumDurationTimeframe: params.MinimumDurationTimeframeMedium,
		},
	)
	if errCrMedium != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errCrMedium,
			}
	}

	c := Coin{
		code:               params.Code,
		defaultQuantityBuy: params.DefaultQuantityBuy,

		periodShort:  periodShort,
		periodMedium: periodMedium,

		ordering: params.Ordering,
	}

	for _, option := range options {
		option(&c)
	}

	return &c, nil
}

func (c *Coin) AddPriceChange(price decimal.Decimal) {
	c.periodShort.AddPriceChange(price)
	c.periodMedium.AddPriceChange(price)

	c.validatePriceChange(price)
}

func (c *Coin) AddPriceChanges(prices []decimal.Decimal) {
	for _, price := range prices {
		c.AddPriceChange(price)
	}
}

func (c *Coin) AddPriceChangesFloat(prices []float64) error {
	for _, price := range prices {
		priceDecimal, errDecimal := decimal.NewFromFloat64(price)
		if errDecimal != nil {
			return errDecimal
		}

		c.AddPriceChange(priceDecimal)
	}

	return nil
}
