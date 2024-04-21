package coin

import (
	"fmt"
	"test/configuration"
	"test/ordering"
	"test/strategies"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestStrategies(t *testing.T) {
	dropPercent, errDrop := decimal.NewFromFloat64(20.)
	require.NoError(t, errDrop)

	strategyDropSudden, errCr := strategies.NewStrategyDropSudden(
		&strategies.ParamsNewStrategyDropSudden{
			PercentDrop: dropPercent,
		},
	)
	require.NoError(t, errCr)

	raisePercent, errRaise := decimal.NewFromFloat64(10.)
	require.NoError(t, errRaise)

	strategySellSimple, errCr := strategies.NewStrategySellSimple(
		&strategies.ParamsNewStrategySellSimple{
			PercentRaise: raisePercent,
		},
	)
	require.NoError(t, errCr)

	c, errCr := NewCoin(
		&ParamsNewCoin{
			Code:     configuration.MON,
			Ordering: ordering.NewOrderingLogOnly(),

			MinimumPriceChangesMediumPeriod: 1,
			MinimumPriceChangesShortPeriod:  1,

			DefaultQuantityBuy: 1,
			MaximumQuantity:    100,

			PercentDeltaIsPriceChangeShort:  5,
			PercentDeltaIsPriceChangeMedium: 5,
		},

		WithStrategyBuy(
			strategyDropSudden,
		),
		WithStrategySell(
			strategySellSimple,
		),
	)
	require.NoError(t, errCr)
	require.NotNil(t, c)

	priceChanges1 := []float64{1., 1.2, 1., .5, .4}
	require.NoError(t,
		c.AddPriceChangesFloat(
			priceChanges1,
		),
	)

	priceChanges2 := []float64{.77}
	require.NoError(t,
		c.AddPriceChangesFloat(
			priceChanges2,
		),
	)

	require.NotZero(t,
		c.periodShort.GetNoPriceChanges(),
	)
	fmt.Println(c.periodShort)

	require.NotZero(t,
		c.periodMedium.GetNoPriceChanges(),
		"periodMedium price changes",
	)

	time.Sleep(2 * time.Second)

	fmt.Println("current quantity:", c.currentQuantity.Load())
}
