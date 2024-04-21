package coin

import (
	"test/strategies"
)

type OptionCoin func(*Coin)

func WithStrategy(strategy strategies.IStrategyBuy) OptionCoin {
	return func(c *Coin) {
		c.strategiesBuy = append(c.strategiesBuy, strategy)
	}
}
