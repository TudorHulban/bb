package timeperiod

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"test/apperrors"
	"test/helpers"
	"time"

	"github.com/asaskevich/govalidator"
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

type TimePeriod struct {
	priceList

	name string

	noPriceChanges            atomic.Uint32
	percentDeltaIsPriceChange decimal.Decimal

	minimumNoPriceChanges    uint32
	minimumDurationTimeframe time.Duration
}

type ParamsNewTimePeriod struct {
	Name string `valid:"required"`

	MinimumPriceChanges       uint32        `valid:"required"`
	MinimumDurationTimeframe  time.Duration `valid:"oprional"`
	PercentDeltaIsPriceChange float64       `valid:"required"`
}

func NewTimePeriod(params *ParamsNewTimePeriod) (*TimePeriod, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: fmt.Sprintf("NewTimePeriod: %s", params.Name),
				Issue:  errVa,
			}
	}

	delta, errConv := decimal.NewFromFloat64(params.PercentDeltaIsPriceChange)
	if errConv != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewTimePeriod",
				Issue:  errConv,
			}
	}

	return &TimePeriod{
			priceList: priceList{
				prices: list.New(),
			},

			minimumNoPriceChanges:     params.MinimumPriceChanges,
			minimumDurationTimeframe:  params.MinimumDurationTimeframe,
			percentDeltaIsPriceChange: delta,

			name: params.Name,
		},
		nil
}

func (p *TimePeriod) isPriceChange(priceNew decimal.Decimal) error {
	p.priceList.mux.Lock()
	defer p.priceList.mux.Unlock()

	if p.priceList.prices.Front() == nil {
		return nil
	}

	return helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: p.priceList.prices.Front().Value.(Price).Value,
			PriceNew: priceNew,
			Delta:    p.percentDeltaIsPriceChange,
		},
	)
}

func (p *TimePeriod) AddPriceChange(price decimal.Decimal) {
	if p.isPriceChange(price) != nil {
		return
	}

	p.mux.Lock()

	if p.minimumDurationTimeframe != 0 {
		last := p.prices.Back()

		if last != nil && time.Since(last.Value.(Price).AtTime) > p.minimumDurationTimeframe {
			p.prices.Remove(last)
		}
	}

	p.prices.PushFront(
		Price{
			AtTime: time.Now(),
			Value:  price,
		},
	)

	p.noPriceChanges.Add(1)

	p.mux.Unlock()
}

func (p *TimePeriod) GetNoPriceChanges() uint32 {
	return p.noPriceChanges.Load()
}

func (p *TimePeriod) getPeriodValues() DecimalValues {
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

func (p *TimePeriod) GetPeriodAverage() decimal.Decimal {
	if errPeriodValid := p.Valid(); errPeriodValid != nil {
		fmt.Println(errPeriodValid)

		return decimal.Zero
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	var sum decimal.Decimal

	for e := p.prices.Front(); e != nil; e = e.Next() {
		var errSum error

		sum, errSum = sum.Add(e.Value.(Price).Value)
		if errSum != nil {
			fmt.Println(errSum)

			return decimal.Zero
		}
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

func (p *TimePeriod) Valid() error {
	p.priceList.mux.Lock()
	defer p.priceList.mux.Unlock()

	last := p.prices.Back()

	timeSinceLastValue := time.Since(last.Value.(Price).AtTime)

	if last != nil && timeSinceLastValue < p.minimumDurationTimeframe {
		return apperrors.ErrValidation{
			Issue: fmt.Errorf(
				"period %s too short (%d), minimum %d",
				p.name,
				timeSinceLastValue,
				p.minimumDurationTimeframe,
			),

			Caller: "timePeriod - Valid",
		}
	}

	numberPriceChanges := p.priceList.prices.Len()

	if numberPriceChanges < int(p.minimumNoPriceChanges) {
		return fmt.Errorf(
			"%s: current number of price changes (%d) smaller than minimum number of price changes (%d)",
			p.name,
			numberPriceChanges,
			p.minimumNoPriceChanges,
		)
	}

	return nil
}

func (p TimePeriod) String() string {
	validity := "is valid"

	if errValid := p.Valid(); errValid != nil {
		validity = errValid.Error()
	}

	return strings.Join(
		[]string{
			"",
			"TimePeriod:",
			fmt.Sprintf("Valid: %s", validity),
			fmt.Sprintf("Minimum number Price Changes: %d", p.minimumNoPriceChanges),
			fmt.Sprintf("Number Price Changes: %d", p.GetNoPriceChanges()),
			fmt.Sprintf("Period Average: %s", p.GetPeriodAverage().String()),
			fmt.Sprintf("Values: %s", p.getPeriodValues().String()),
			"",
		},
		"\n",
	)
}
