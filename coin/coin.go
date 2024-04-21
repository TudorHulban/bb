package coin

import (
	"test/apperrors"
	"test/configuration"
	"test/helpers"
	"test/ordering"
	"test/strategies"
	"test/timeperiod"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/govalues/decimal"
)

type Coin struct {
	ordering       ordering.IOrdering
	strategiesBuy  []strategies.IStrategyBuy
	strategiesSell []strategies.IStrategySell

	code configuration.CoinCODE

	periodShort  *timeperiod.TimePeriod
	periodMedium *timeperiod.TimePeriod

	averageBuyPrice helpers.Price

	*quantities
}

type ParamsNewCoin struct {
	Ordering ordering.IOrdering
	Code     configuration.CoinCODE `valid:"required"`

	MinimumPriceChangesShortPeriod  uint32 `valid:"required"`
	MinimumPriceChangesMediumPeriod uint32 `valid:"required"`

	DefaultQuantityBuy uint32 `valid:"required"`
	MaximumQuantity    uint32 `valid:"required"`

	MinimumDurationTimeframeShort  time.Duration
	MinimumDurationTimeframeMedium time.Duration

	PercentDeltaIsPriceChangeShort  float64 `valid:"required"`
	PercentDeltaIsPriceChangeMedium float64 `valid:"required"`
}

// NewCoin - if minimum duration is zero, no deletion would occur for old element.
func NewCoin(params *ParamsNewCoin, options ...OptionCoin) (*Coin, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errVa,
			}
	}

	periodShort, errCrShort := timeperiod.NewTimePeriod(
		&timeperiod.ParamsNewTimePeriod{
			Name: timeperiod.NamePeriodShort,

			MinimumPriceChanges:       params.MinimumPriceChangesShortPeriod,
			MinimumDurationTimeframe:  params.MinimumDurationTimeframeShort,
			PercentDeltaIsPriceChange: params.PercentDeltaIsPriceChangeShort,
		},
	)
	if errCrShort != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errCrShort,
			}
	}

	periodMedium, errCrMedium := timeperiod.NewTimePeriod(
		&timeperiod.ParamsNewTimePeriod{
			Name: timeperiod.NamePeriodMedium,

			MinimumPriceChanges:       params.MinimumPriceChangesMediumPeriod,
			MinimumDurationTimeframe:  params.MinimumDurationTimeframeMedium,
			PercentDeltaIsPriceChange: params.PercentDeltaIsPriceChangeMedium,
		},
	)
	if errCrMedium != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errCrMedium,
			}
	}

	quant, errCrQuantities := NewQuantities(
		&ParamNewQuantities{
			QuantityBuy:     params.DefaultQuantityBuy,
			QuantityMaximum: params.MaximumQuantity,
		},
	)
	if errCrQuantities != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewCoin",
				Issue:  errCrQuantities,
			}
	}

	c := Coin{
		code:       params.Code,
		quantities: quant,

		periodShort:  periodShort,
		periodMedium: periodMedium,

		ordering: params.Ordering,
	}

	for _, option := range options {
		option(&c)
	}

	return &c, nil
}

func (c *Coin) AddPriceChange(price decimal.Decimal) {
	c.periodShort.AddPriceChange(price)
	c.periodMedium.AddPriceChange(price)

	c.validatePriceChangeBuy(price)

	time.Sleep(1 * time.Second)

	c.validatePriceChangeSell(price)
}

func (c *Coin) AddPriceChanges(prices []decimal.Decimal) {
	for _, price := range prices {
		c.AddPriceChange(price)
	}
}

func (c *Coin) AddPriceChangesFloat(prices []float64) error {
	for _, price := range prices {
		priceDecimal, errDecimal := decimal.NewFromFloat64(price)
		if errDecimal != nil {
			return errDecimal
		}

		c.AddPriceChange(priceDecimal)
	}

	return nil
}

func (c *Coin) setAveragePriceBuy() {

}
