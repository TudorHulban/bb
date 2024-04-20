package coin

import (
	"fmt"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestTimePeriod(t *testing.T) {
	p := NewTimePeriod(
		&ParamsNewTimePeriod{
			MinimumPriceChanges: 1,
		},
	)
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
