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
			periodMediumAverage := c.periodMedium.GetPeriodAverage()
			if periodMediumAverage == decimal.Zero {
				continue
			}

			strategy.SetPrice(periodMediumAverage)
		}

		fmt.Println(strategy)

		action, errStrategy := strategy.AddPriceChange(
			&strategies.ParamsAddPriceChange{
				PriceNow: price,

				NoPriceChangesPeriodShort:  c.periodShort.GetNoPriceChanges(),
				NoPriceChangesPeriodMedium: c.periodMedium.GetNoPriceChanges(),
			},
		)
		if errStrategy != nil {
			fmt.Printf(
				"validatePriceChange: %s for %s.\n",
				errStrategy.Error(),
				price.String(),
			)
		}
		if action != ordering.DoNothing {
			fmt.Printf(
				"validatePriceChange: %s at %s.\n",
				action,
				price,
			)
		}
	}
}
