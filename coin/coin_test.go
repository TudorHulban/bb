package coin

import (
	"test/ordering"
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestCoin(t *testing.T) {
	c := NewCoin(
		&ParamsNewCoin{
			Ordering: ordering.NewOrderingLogOnly(),
		},
	)

	c.AddPriceChange(decimal.One)

	lowPrice, errConversionLow := decimal.NewFromFloat64(1.00)
	require.NoError(t, errConversionLow)

	require.Error(t,
		c.isPriceChange(lowPrice),
	)

	highPrice, errConversionHigh := decimal.NewFromFloat64(1.001)
	require.NoError(t, errConversionHigh)

	require.NoError(t,
		c.isPriceChange(highPrice),
		highPrice.String(),
	)
}

func BenchmarkCoinPriceChange(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	c := NewCoin(
		&ParamsNewCoin{
			Ordering: ordering.NewOrderingLogOnly(),
		},
	)

	c.AddPriceChange(decimal.One)

	newPrice, _ := decimal.NewFromFloat64(1.00)

	for n := 0; n < b.N; n++ {
		c.isPriceChange(newPrice)
	}
}
