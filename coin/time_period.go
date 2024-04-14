package coin

import (
	"container/list"
	"sync"
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
