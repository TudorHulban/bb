package coin

import (
	"errors"
	"fmt"
	"test/configuration"
	"test/helpers"
	"test/ordering"
	"test/strategies"
	"time"

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

	secondsShortPeriod          uint
	secondsMediumPeriodInterval uint

	percentDeltaChange decimal.Decimal
	currentQuantity    decimal.Decimal

	strategies []strategies.IStrategy
}

type ParamsNewCoin struct {
	Ordering ordering.IOrdering

	MinimumPriceChangesShortPeriod  int
	MinimumPriceChangesMediumPeriod int
}

func NewCoin(params *ParamsNewCoin, options ...OptionCoin) *Coin {
	deltaChange, _ := decimal.NewFromFloat64(configuration.DefaultPercentDeltaChange)

	c := Coin{
		periodShort:  NewTimePeriod(params.MinimumPriceChangesShortPeriod),
		periodMedium: NewTimePeriod(params.MinimumPriceChangesMediumPeriod),

		secondsShortPeriod:          configuration.DefaultSecondsShortPeriod,
		secondsMediumPeriodInterval: configuration.DefaultSecondsMediumPeriod,
		percentDeltaChange:          deltaChange,

		ordering: params.Ordering,
	}

	for _, option := range options {
		option(&c)
	}

	return &c
}

func (c *Coin) AverageMediumPeriod() (decimal.Decimal, error) {
	c.periodMedium.mux.Lock()
	defer c.periodMedium.mux.Unlock()

	if c.periodMedium.prices.Len() == 0 {
		return decimal.Zero,
			errors.New("empty")
	}

	last := c.periodMedium.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) < time.Duration(c.secondsMediumPeriodInterval)*time.Second {
		return decimal.Zero,
			errors.New("medium period too short")
	}

	var sum decimal.Decimal

	for e := c.periodMedium.prices.Front(); e != nil; e = e.Next() {
		sum.Add(e.Value.(Price).Value)
	}

	length, errConverasion := decimal.NewFromInt64(
		int64(c.periodMedium.prices.Len()),
		0,
		0,
	)
	if errConverasion != nil {
		return decimal.Zero,
			errConverasion
	}

	return sum.Quo(length)
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
			Delta:    c.percentDeltaChange,
		},
	)
}

func (c *Coin) AddPriceChange(price decimal.Decimal) {
	if c.isPriceChange(price) != nil {
		return
	}

	c.periodShort.mux.Lock()

	last := c.periodShort.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) > time.Duration(c.secondsShortPeriod)*time.Second {
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

func (c *Coin) getShortPeriodNoPriceChanges() int {
	c.periodShort.mux.Lock()
	defer c.periodShort.mux.Unlock()

	return c.periodShort.prices.Len()
}

func (c *Coin) getMediumPeriodNoPriceChanges() int {
	c.periodMedium.mux.Lock()
	defer c.periodMedium.mux.Unlock()

	return c.periodShort.prices.Len()
}

func (c *Coin) validatePriceChange(price decimal.Decimal) {
	for _, strategy := range c.strategies {
		if !strategy.IsReady() {

		}

		action, errStrategy := strategy.AddPriceChange(
			&strategies.ParamsAddPriceChange{
				PriceNow: price,

				NoPriceChangesPeriodShort:  c.getShortPeriodNoPriceChanges(),
				NoPriceChangesPeriodMedium: c.getMediumPeriodNoPriceChanges(),
			},
		)
		if errStrategy != nil {
			fmt.Println(errStrategy)
		}
		if action != ordering.DoNothing {
			fmt.Printf(
				"%s at %s.\n",
				action,
				price,
			)
		}
	}
}
