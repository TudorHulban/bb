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
	dropPercent, errDrop := decimal.NewFromFloat64(10.)
	require.NoError(t, errDrop)

	c, errCr := NewCoin(
		&ParamsNewCoin{
			Ordering: ordering.NewOrderingLogOnly(),

			MinimumPriceChangesMediumPeriod: 1,
			MinimumPriceChangesShortPeriod:  1,
			MinimumSecondsTimeframeShort:    1,
			MinimumSecondsTimeframeMedium:   1,
		},

		WithStrategy(
			strategies.NewStrategyDropSudden(
				strategies.ParamsNewStrategyDropSudden{
					DropPercent: dropPercent,
				},
			),
		),
	)
	require.NoError(t, errCr)
	require.NotNil(t, c)

	priceChanges1 := []float64{1., 1.2, 1., .5}
	require.NoError(t,
		c.AddPriceChangesFloat(
			priceChanges1,
		),
	)

	time.Sleep(1 * time.Second)

	priceChanges2 := []float64{.77}
	require.NoError(t,
		c.AddPriceChangesFloat(
			priceChanges2,
		),
	)

	require.EqualValues(t,
		c.periodShort.GetNoPriceChanges(),
		len(priceChanges1)+len(priceChanges2)-1,
	)
	fmt.Println(c.periodShort)

	require.NotZero(t,
		c.periodMedium.GetNoPriceChanges(),
		"periodMedium price changes",
	)
}
