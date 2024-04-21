package timeperiod

import (
	"fmt"
	"test/configuration"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestIsPriceChange(t *testing.T) {
	p, errCr := NewTimePeriod(
		&ParamsNewTimePeriod{
			Name:                      NamePeriodShort,
			MinimumPriceChanges:       1,
			PercentDeltaIsPriceChange: configuration.DefaultPercentDeltaIsPriceChangeShort,
		},
	)
	require.NoError(t, errCr)
	require.NotZero(t, p)

	p.AddPriceChange(decimal.One)

	lowPrice, errConversionLow := decimal.NewFromFloat64(1.00)
	require.NoError(t, errConversionLow)

	require.Error(t,
		p.isPriceChange(lowPrice),
	)

	highPrice, errConversionHigh := decimal.NewFromFloat64(1.001)
	require.NoError(t, errConversionHigh)

	require.NoError(t,
		p.isPriceChange(highPrice),
		highPrice.String(),
	)
}

func BenchmarkCoinPriceChange(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	p, errCr := NewTimePeriod(
		&ParamsNewTimePeriod{
			Name:                      NamePeriodShort,
			MinimumPriceChanges:       1,
			PercentDeltaIsPriceChange: configuration.DefaultPercentDeltaIsPriceChangeShort,
		},
	)
	require.NoError(b, errCr)

	p.AddPriceChange(decimal.One)

	newPrice, _ := decimal.NewFromFloat64(1.00)

	for n := 0; n < b.N; n++ {
		p.isPriceChange(newPrice)
	}
}

func TestTimePeriod(t *testing.T) {
	p, errCr := NewTimePeriod(
		&ParamsNewTimePeriod{
			MinimumPriceChanges: 1,
		},
	)
	require.NoError(t, errCr)
	require.NotZero(t, p)

	p.AddPriceChange(decimal.One)
	p.AddPriceChange(decimal.Ten)

	time.Sleep(
		time.Duration(11) * time.Millisecond,
	)

	p.AddPriceChange(decimal.One)

	fmt.Println(p)

	require.NotZero(t,
		p.GetPeriodAverage(),
	)
}
