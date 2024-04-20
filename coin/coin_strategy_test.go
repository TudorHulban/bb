package coin

import (
	"fmt"
	"test/ordering"
	"test/strategies"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestStrategyDropSudden(t *testing.T) {
	dropPercent, errDrop := decimal.NewFromFloat64(20.)
	require.NoError(t, errDrop)

	strategyDropSudden, errCr := strategies.NewStrategyDropSudden(
		strategies.ParamsNewStrategyDropSudden{
			DropPercent: dropPercent,
		},
	)
	require.NoError(t, errCr)

	c, errCr := NewCoin(
		&ParamsNewCoin{
			Ordering: ordering.NewOrderingLogOnly(),

			MinimumPriceChangesMediumPeriod: 1,
			MinimumPriceChangesShortPeriod:  1,

			DefaultQuantityBuy: 1,
		},

		WithStrategy(
			strategyDropSudden,
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

	fmt.Println(c.currentQuantity.Load())
}
