package coin

import (
	"fmt"
	"test/ordering"
	"test/strategies"

	"github.com/govalues/decimal"
)

func (c *Coin) validatePriceChange(price decimal.Decimal) {
	for _, strategy := range c.strategies {
		if !strategy.IsReady() {
			fmt.Println("xxxxxxxxxxxxx", strategy.IsReady())

			if errPeriodMedium := c.periodMedium.Valid(); errPeriodMedium != nil {
				fmt.Println(errPeriodMedium)

				periodMediumAverage := c.periodMedium.GetPeriodAverage()
				if periodMediumAverage == decimal.Zero {
					continue
				}

				strategy.SetPrice(periodMediumAverage)

				continue
			}

			c.periodMedium.AddPriceChange(price)

			continue
		}

		action, errStrategy := strategy.AddPriceChange(
			&strategies.ParamsAddPriceChange{
				PriceNow: price,

				NoPriceChangesPeriodShort:  c.periodShort.GetNoPriceChanges(),
				NoPriceChangesPeriodMedium: c.periodMedium.GetNoPriceChanges(),
			},
		)
		if errStrategy != nil {
			fmt.Println(errStrategy)
		}
		if action != ordering.DoNothing {
			fmt.Printf(
				"%s at %s.\n",
				action,
				price,
			)
		}
	}
}
