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

	MinimumPriceChangesShortPeriod  uint `valid:"required"`
	MinimumPriceChangesMediumPeriod uint `valid:"required"`

	MinimumSecondsTimeframeShort  uint `valid:"required"`
	MinimumSecondsTimeframeMedium uint `valid:"required"`
}

func NewCoin(params *ParamsNewCoin, options ...OptionCoin) (*Coin, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return nil,
			apperrors.ErrServiceValidation{
				Caller: "NewCoin",
				Issue:  errVa,
			}
	}

	deltaChange, _ := decimal.NewFromFloat64(configuration.DefaultPercentDeltaIsPriceChange)

	c := Coin{
		periodShort: NewTimePeriod(&ParamsNewTimePeriod{
			MinimumPriceChanges:     params.MinimumPriceChangesShortPeriod,
			MinimumSecondsTimeframe: params.MinimumSecondsTimeframeShort,
		}),
		periodMedium: NewTimePeriod(&ParamsNewTimePeriod{
			MinimumPriceChanges:     params.MinimumPriceChangesMediumPeriod,
			MinimumSecondsTimeframe: params.MinimumPriceChangesShortPeriod,
		}),

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

	c.periodShort.mux.Lock()

	last := c.periodShort.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) > time.Duration(c.periodShort.minimumSecondsTimeframe)*time.Second {
		c.periodShort.prices.Remove(last)
	}

	c.periodShort.prices.PushFront(
		Price{
			AtTime: time.Now(),
			Value:  price,
		},
	)

	c.periodShort.mux.Unlock()

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
