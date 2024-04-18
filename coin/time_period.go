package coin

import (
	"container/list"
	"fmt"
	"strings"
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

	minimumNoPriceChanges   uint
	minimumSecondsTimeframe uint
}

type ParamsNewTimePeriod struct {
	MinimumPriceChanges     uint
	MinimumSecondsTimeframe uint
}

func NewTimePeriod(params *ParamsNewTimePeriod) timePeriod {
	return timePeriod{
		priceList: priceList{
			prices: list.New(),
		},

		minimumNoPriceChanges:   params.MinimumPriceChanges,
		minimumSecondsTimeframe: params.MinimumSecondsTimeframe,
	}
}

func (p *timePeriod) Valid() bool {
	p.priceList.mux.Lock()
	defer p.priceList.mux.Unlock()

	last := p.prices.Back()

	if last != nil && time.Since(last.Value.(Price).AtTime) < time.Duration(p.minimumSecondsTimeframe)*time.Second {
		return false // "period too short"
	}

	return p.priceList.prices.Len() > int(p.minimumNoPriceChanges)
}

func (p timePeriod) String() string {
	return strings.Join(
		[]string{
			"TimePeriod:",
			fmt.Sprintf("Valid: %t", p.Valid()),
			fmt.Sprintf("Number Price Changes: %d", p.GetNoPriceChanges()),
			fmt.Sprintf("Period Average: %s", p.GetPeriodAverage().String()),
			fmt.Sprintf("Values: %s", p.getPeriodValues().String()),
		},
		"\n",
	)
}

func (p *timePeriod) GetNoPriceChanges() int {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.prices.Len()
}

func (p *timePeriod) getPeriodValues() DecimalValues {
	p.mux.Lock()
	defer p.mux.Unlock()

	var result []decimal.Decimal

	for e := p.prices.Front(); e != nil; e = e.Next() {
		result = append(result,
			e.Value.(Price).Value,
		)
	}

	return result
}

func (p *timePeriod) GetPeriodAverage() decimal.Decimal {
	if !p.Valid() {
		return decimal.Zero
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	var sum decimal.Decimal

	for e := p.prices.Front(); e != nil; e = e.Next() {
		sum.Add(e.Value.(Price).Value)
	}

	length, errConversion := decimal.NewFromInt64(
		int64(p.prices.Len()),
		0,
		0,
	)
	if errConversion != nil {
		fmt.Println(errConversion)

		return decimal.Zero
	}

	average, errAverage := sum.Quo(length)
	if errAverage != nil {
		fmt.Println(errAverage)

		return decimal.Zero
	}

	return average
}
