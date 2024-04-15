package coin

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/govalues/decimal"
)

type priceList struct {
	prices *list.List

	mux sync.RWMutex
}

type timePeriod struct {
	priceList

	minimumNoPriceChanges   int
	minimumSecondsTimeframe uint
}

func NewTimePeriod(minimumPriceChanges int) timePeriod {
	return timePeriod{
		priceList: priceList{
			prices: list.New(),
		},

		minimumNoPriceChanges: minimumPriceChanges,
	}
}

func (p *timePeriod) Valid() bool {
	p.priceList.mux.Lock()
	defer p.priceList.mux.Unlock()

	// TODO: add timeframe check

	return p.priceList.prices.Len() > p.minimumNoPriceChanges
}

func (p *timePeriod) GetNoPriceChanges() int {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.prices.Len()
}

func (p *timePeriod) GetPeriodAverage() (decimal.Decimal, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.prices.Len() == 0 {
		return decimal.Zero,
			errors.New("empty")
	}

	last := p.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) < time.Duration(p.minimumSecondsTimeframe)*time.Second {
		return decimal.Zero,
			errors.New("medium period too short")
	}

	var sum decimal.Decimal

	for e := p.prices.Front(); e != nil; e = e.Next() {
		sum.Add(e.Value.(Price).Value)
	}

	length, errConverasion := decimal.NewFromInt64(
		int64(p.prices.Len()),
		0,
		0,
	)
	if errConverasion != nil {
		return decimal.Zero,
			errConverasion
	}

	return sum.Quo(length)
}
