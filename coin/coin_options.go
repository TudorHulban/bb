package coin

import (
	"test/strategies"
)

type OptionCoin func(*Coin)

func WithStrategyBuy(strategy strategies.IStrategyBuy) OptionCoin {
	return func(c *Coin) {
		c.strategiesBuy = append(c.strategiesBuy, strategy)
	}
}

func WithStrategySell(strategy strategies.IStrategySell) OptionCoin {
	return func(c *Coin) {
		c.strategiesSell = append(c.strategiesSell, strategy)
	}
}
