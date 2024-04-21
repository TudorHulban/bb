package helpers

import (
	"sync"

	"github.com/govalues/decimal"
)

type Price struct {
	price decimal.Decimal

	mux sync.RWMutex
}

func NewPriceFrom(price float64) (*Price, error) {
	p, errCr := decimal.NewFromFloat64(price)
	if errCr != nil {
		return nil, errCr
	}

	return &Price{
			price: p,
		},
		nil
}

func (p *Price) GetPrice() decimal.Decimal {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.price
}

func (p *Price) SetPrice(newPrice decimal.Decimal) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.price = newPrice
}

func (p *Price) HoldAverage(withPrice decimal.Decimal) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.price.IsZero() {
		p.price = withPrice

		return nil
	}

	sum, errAdd := p.price.Add(withPrice)
	if errAdd != nil {
		return errAdd
	}

	average, errAverage := sum.Quo(decimal.Two)
	if errAverage != nil {
		return errAverage
	}

	p.price = average

	return nil
}
