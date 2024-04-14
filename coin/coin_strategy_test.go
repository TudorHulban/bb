package coin

import (
	"test/ordering"
	"test/strategies"
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestStrategyDropSudden(t *testing.T) {
	dropPercent, errDrop := decimal.NewFromFloat64(10.)
	require.NoError(t, errDrop)

	c := NewCoin(
		&ParamsNewCoin{
			Ordering: ordering.NewOrderingLogOnly(),
		},

		WithStrategy(
			strategies.NewStrategyDropSudden(
				strategies.ParamsNewStrategyDropSudden{
					DropPercent: dropPercent,
				},
			),
		),
	)
	require.NotNil(t, c)

	require.NoError(t,
		c.AddPriceChangesFloat(
			[]float64{1., 1.2, 1., .5},
		),
	)
}
