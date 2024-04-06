package coin

import "github.com/govalues/decimal"

type OptionCoin func(*Coin)

func WithPercentageDelta(delta decimal.Decimal) OptionCoin {
	return func(c *Coin) {
		c.percentDeltaChange = delta
	}
}
