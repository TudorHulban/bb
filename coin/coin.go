package coin

import (
	"container/list"
	"errors"
	"sync"
	"test/configuration"
	"test/helpers"
	"test/ordering"
	"time"

	"github.com/govalues/decimal"
)

type Price struct {
	AtTime time.Time
	Value  decimal.Decimal
}

type priceList struct {
	prices *list.List

	mux sync.RWMutex
}

type Coin struct {
	ordering ordering.IOrdering

	pricesShortPeriod  priceList
	pricesMediumPeriod priceList

	secondsShortPeriod          uint
	secondsMediumPeriodInterval uint
	secondsMinimumTimeframe     uint // between oldest medium and now

	percentDeltaChange decimal.Decimal
	currentQuantity    decimal.Decimal

	mux sync.RWMutex
}

type ParamsNewCoin struct {
	Ordering ordering.IOrdering
}

func NewCoin(params *ParamsNewCoin, options ...OptionCoin) *Coin {
	deltaChange, _ := decimal.NewFromFloat64(configuration.DefaultPercentDeltaChange)

	c := Coin{
		pricesShortPeriod: priceList{
			prices: list.New(),
		},
		pricesMediumPeriod: priceList{
			prices: list.New(),
		},

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
	c.pricesMediumPeriod.mux.Lock()
	defer c.pricesMediumPeriod.mux.Unlock()

	if c.pricesMediumPeriod.prices.Len() == 0 {
		return decimal.Zero,
			errors.New("empty")
	}

	last := c.pricesMediumPeriod.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) < time.Duration(c.secondsMediumPeriodInterval)*time.Second {
		return decimal.Zero,
			errors.New("medium period too short")
	}

	var sum decimal.Decimal

	for e := c.pricesMediumPeriod.prices.Front(); e != nil; e = e.Next() {
		sum.Add(e.Value.(Price).Value)
	}

	length, errConverasion := decimal.NewFromInt64(
		int64(c.pricesMediumPeriod.prices.Len()),
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
	c.pricesShortPeriod.mux.Lock()
	defer c.pricesShortPeriod.mux.Unlock()

	if c.pricesShortPeriod.prices.Front() == nil {
		return nil
	}

	return helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: c.pricesShortPeriod.prices.Front().Value.(Price).Value,
			PriceNew: priceNew,
			Delta:    c.percentDeltaChange,
		},
	)
}

func (c *Coin) AddPriceChange(price decimal.Decimal) {
	if c.isPriceChange(price) != nil {
		return
	}

	c.pricesShortPeriod.mux.Lock()

	last := c.pricesShortPeriod.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) > time.Duration(c.secondsShortPeriod)*time.Second {
		c.pricesShortPeriod.prices.Remove(last)
	}

	c.pricesShortPeriod.prices.PushFront(
		Price{
			AtTime: time.Now(),
			Value:  price,
		},
	)

	c.pricesShortPeriod.mux.Unlock()
}

func (c *Coin) PriceChanges() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.pricesShortPeriod.prices.Len()
}
