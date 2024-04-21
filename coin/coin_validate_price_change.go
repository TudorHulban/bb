package coin

import (
	"fmt"
	"test/ordering"
	"test/strategies"

	"github.com/govalues/decimal"
)

func (c *Coin) validatePriceChange(price decimal.Decimal) {
	for _, strategy := range c.strategiesBuy {
		periodMediumAverage := c.periodMedium.GetPeriodAverage()
		if periodMediumAverage == decimal.Zero {
			continue
		}

		strategy.SetPrice(periodMediumAverage)

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

		if action == ordering.Buy {
			if c.currentQuantity.Load() == 0 {
				go func() {
					if !strategy.CanPlaceOrder() {
						return
					}

					strategy.IncrementSimultaneousOrders()
					defer strategy.DecrementSimultaneousOrders()

					if errBuy := c.ordering.Buy(
						ordering.ParamsOder{
							Code:     c.code,
							Quantity: c.defaultQuantityBuy,
						},
					); errBuy != nil {
						fmt.Println(errBuy)

						return
					}

					c.currentQuantity.Add(
						c.defaultQuantityBuy,
					)
				}()
			}
		}
	}
}
