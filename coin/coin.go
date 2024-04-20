package coin

import (
	"test/apperrors"
	"test/configuration"
	"test/helpers"
	"test/ordering"
	"test/strategies"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/govalues/decimal"
)

type Price struct {
	AtTime time.Time
	Value  decimal.Decimal
}

type Coin struct {
	ordering ordering.IOrdering

	periodShort  timePeriod
	periodMedium timePeriod

	percentDeltaIsPriceChange decimal.Decimal
	currentQuantity           decimal.Decimal

	strategies []strategies.IStrategy
}

type ParamsNewCoin struct {
	Ordering ordering.IOrdering

	MinimumPriceChangesShortPeriod  uint32 `valid:"required"`
	MinimumPriceChangesMediumPeriod uint32 `valid:"required"`

	MinimumDurationTimeframeShort  time.Duration
	MinimumDurationTimeframeMedium time.Duration
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

	deltaChange, _ := decimal.NewFromFloat64(configuration.DefaultPercentDeltaIsPriceChange)

	c := Coin{
		periodShort: NewTimePeriod(
			&ParamsNewTimePeriod{
				name:                     namePeriodShort,
				MinimumPriceChanges:      params.MinimumPriceChangesShortPeriod,
				minimumDurationTimeframe: params.MinimumDurationTimeframeShort,
			},
		),
		periodMedium: NewTimePeriod(
			&ParamsNewTimePeriod{
				name:                     namePeriodMedium,
				MinimumPriceChanges:      params.MinimumPriceChangesMediumPeriod,
				minimumDurationTimeframe: params.MinimumDurationTimeframeMedium,
			},
		),

		percentDeltaIsPriceChange: deltaChange,

		ordering: params.Ordering,
	}

	for _, option := range options {
		option(&c)
	}

	return &c, nil
}

func (c *Coin) isPriceChange(priceNew decimal.Decimal) error {
	c.periodShort.mux.Lock()
	defer c.periodShort.mux.Unlock()

	if c.periodShort.prices.Front() == nil {
		return nil
	}

	return helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: c.periodShort.prices.Front().Value.(Price).Value,
			PriceNew: priceNew,
			Delta:    c.percentDeltaIsPriceChange,
		},
	)
}

func (c *Coin) AddPriceChange(price decimal.Decimal) {
	if c.isPriceChange(price) != nil {
		return
	}

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
