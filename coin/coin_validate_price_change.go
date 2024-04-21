package coin

import (
	"fmt"
	"log"
	"sync/atomic"
	"test/ordering"
	"test/strategies"

	"github.com/govalues/decimal"
)

func (c *Coin) validatePriceChangeBuy(priceNow decimal.Decimal) {
	if c.canBuy() {
		for _, strategy := range c.strategiesBuy {
			periodMediumAverage := c.periodMedium.GetPeriodAverage()
			if periodMediumAverage == decimal.Zero {
				continue
			}

			strategy.SetPrice(periodMediumAverage)

			fmt.Println(strategy)

			action, errStrategy := strategy.AddPriceChange(
				&strategies.ParamsAddPriceChangeBuy{
					PriceNow: priceNow,

					NoPriceChangesPeriodShort:  c.periodShort.GetNoPriceChanges(),
					NoPriceChangesPeriodMedium: c.periodMedium.GetNoPriceChanges(),
				},
			)
			if errStrategy != nil {
				fmt.Printf(
					"validatePriceChange: %s for %s.\n",
					errStrategy.Error(),
					priceNow.String(),
				)
			}
			if action != ordering.DoNothing {
				fmt.Printf(
					"validatePriceChange: %s at %s.\n",
					action,
					priceNow,
				)
			}

			if action == ordering.Buy {
				if c.canBuy() {
					go func() {
						if !strategy.CanPlaceOrder() {
							return
						}

						strategy.IncrementSimultaneousOrders()
						defer strategy.DecrementSimultaneousOrders()

						if errBuy := c.ordering.Buy(
							ordering.ParamsOder{
								Code:     c.code,
								Quantity: c.quantityBuy,
							},
						); errBuy != nil {
							fmt.Println(errBuy)

							return
						}

						c.currentQuantity.Add(
							c.quantityBuy,
						)
						c.averageBuyPrice.HoldAverage(priceNow)

						log.Println("current quantity:", c.currentQuantity.Load())
					}()
				}
			}
		}
	}
}

func (c *Coin) validatePriceChangeSell(price decimal.Decimal) {
	log.Println("last buy price:", c.averageBuyPrice.GetPrice().String())

	if !c.canSell() || c.averageBuyPrice.GetPrice().IsZero() {
		return
	}

	for _, strategy := range c.strategiesSell {
		strategy.SetPrice(c.averageBuyPrice.GetPrice())

		fmt.Println(strategy)

		action, errStrategy := strategy.AddPriceChange(
			&strategies.ParamsAddPriceChangeSell{
				PriceBuy: c.averageBuyPrice.GetPrice(),
				PriceNow: price,
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

		if action == ordering.Sell {
			if c.canSell() && !c.averageBuyPrice.GetPrice().IsZero() {
				go func() {
					if errSell := c.ordering.Sell(
						ordering.ParamsOder{
							Code:     c.code,
							Quantity: c.currentQuantity.Load(),
						},
					); errSell != nil {
						fmt.Println(errSell)

						return
					}

					c.currentQuantity = atomic.Uint32{}
				}()
			}
		}
	}
}
