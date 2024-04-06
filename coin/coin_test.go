package coin

import (
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestCoin(t *testing.T) {
	c := NewCoin()

	c.AddPriceChange(decimal.One)

	lowPrice, errConversionLow := decimal.NewFromFloat64(1.00)
	require.NoError(t, errConversionLow)

	require.Error(t,
		c.IsPriceChange(lowPrice),
	)

	highPrice, errConversionHigh := decimal.NewFromFloat64(1.001)
	require.NoError(t, errConversionHigh)

	require.NoError(t,
		c.IsPriceChange(highPrice),
		highPrice.String(),
	)
}

func BenchmarkCoinPriceChange(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	c := NewCoin()

	c.AddPriceChange(decimal.One)

	newPrice, _ := decimal.NewFromFloat64(1.00)

	for n := 0; n < b.N; n++ {
		c.IsPriceChange(newPrice)
	}
}
