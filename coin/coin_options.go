package coin

import (
	"test/strategies"

	"github.com/govalues/decimal"
)

type OptionCoin func(*Coin)

func WithPercentageDelta(delta decimal.Decimal) OptionCoin {
	return func(c *Coin) {
		c.percentDeltaChange = delta
	}
}

func WithStrategy(strategy strategies.IStrategy) OptionCoin {
	return func(c *Coin) {
		c.strategies = append(c.strategies, strategy)
	}
}
