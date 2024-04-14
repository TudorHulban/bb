package helpers

import (
	"test/apperrors"
	"testing"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"
)

func TestPriceChangeByPercent(t *testing.T) {
	type testCase struct {
		description string
		params      ParamsPriceChangeByPercent
		check       func(err error, t *testing.T)
	}

	tests := []testCase{
		{
			"1. error zero - zero",
			ParamsPriceChangeByPercent{
				Delta: decimal.One,
			},
			func(err error, t *testing.T) {
				require.ErrorAs(t,
					err, &apperrors.ErrorInvalidInputs{},
				)
			},
		},
		{
			"2. error zero - valid",
			ParamsPriceChangeByPercent{
				PriceNew: decimal.Hundred,
				Delta:    decimal.One,
			},
			func(err error, t *testing.T) {
				require.NoError(t, err)
			},
		},
		{
			"3. error valid - zero",
			ParamsPriceChangeByPercent{
				PriceOld: decimal.Hundred,
				Delta:    decimal.One,
			},
			func(err error, t *testing.T) {
				require.NoError(t, err)
			},
		},
		{
			"4. valid - no price change",
			ParamsPriceChangeByPercent{
				PriceOld: decimal.One,
				PriceNew: decimal.One,
				Delta:    decimal.One,
			},
			func(err error, t *testing.T) {
				require.Error(t, err)
			},
		},
		{
			"5. valid - no price change",
			ParamsPriceChangeByPercent{
				PriceOld: decimal.One,
				PriceNew: decimal.Two,
				Delta:    decimal.Thousand,
			},
			func(err error, t *testing.T) {
				require.Error(t, err)
			},
		},
		{
			"6. valid - price change",
			ParamsPriceChangeByPercent{
				PriceOld: decimal.One,
				PriceNew: decimal.Two,
				Delta:    decimal.Hundred,
			},
			func(err error, t *testing.T) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				tc.check(
					PriceChangeByPercent(&tc.params),
					t,
				)
			},
		)
	}
}
