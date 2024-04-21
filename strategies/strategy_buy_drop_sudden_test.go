package strategies

import (
	"test/ordering"
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestStrategyDropSudden(t *testing.T) {
	p1, errP1 := decimal.NewFromFloat64(2.)
	require.NoError(t, errP1)

	p2, errP2 := decimal.NewFromFloat64(50.)
	require.NoError(t, errP2)

	s := StrategyDropSudden{
		PriceBeforeDrop: p1,
		PercentDrop:     p2,
	}

	p3, errP3 := decimal.NewFromFloat64(1.)
	require.NoError(t, errP3)

	action, errAdd := s.AddPriceChange(
		&ParamsAddPriceChangeBuy{
			AverageMediumPeriodPrice: p1,
			PriceNow:                 p3,
		},
	)
	require.NoError(t, errAdd)
	require.NotNil(t, action)
	require.EqualValues(t,
		ordering.Buy, action,
	)
}
