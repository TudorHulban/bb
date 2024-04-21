package strategies

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"test/apperrors"
	"test/helpers"
	"test/ordering"
	"test/timeperiod"

	"github.com/govalues/decimal"
)

const nameStrategyDropSudden = "Drop Sudden"

type StrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	PercentDrop     decimal.Decimal

	numberSimultaneousOrders        *uint32
	allowedNumberSimultaneousOrders uint16
}

type ParamsNewStrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	PercentDrop     decimal.Decimal

	AllowedNumberSimultaneousOrders uint16
}

func NewStrategyDropSudden(params *ParamsNewStrategyDropSudden) (*StrategyDropSudden, error) {
	if !params.PercentDrop.IsPos() {
		return nil,
			apperrors.ErrorInvalidInput{
				InputName: "DropPercent",
			}
	}

	var zero uint32

	return &StrategyDropSudden{
			PriceBeforeDrop:                 params.PriceBeforeDrop,
			PercentDrop:                     params.PercentDrop,
			allowedNumberSimultaneousOrders: params.AllowedNumberSimultaneousOrders,

			numberSimultaneousOrders: &zero,
		},
		nil
}

func (s *StrategyDropSudden) SetPrice(price decimal.Decimal) error {
	if !price.IsPos() {
		return apperrors.ErrorInvalidInput{
			InputName: "SetPrice",
		}
	}

	s.PriceBeforeDrop = price

	return nil
}

func (s *StrategyDropSudden) AddPriceChange(params *ParamsAddPriceChangeBuy) (ordering.Action, error) {
	difference, errSubtract := s.PriceBeforeDrop.Sub(params.PriceNow)
	if errSubtract != nil {
		return ordering.DoNothing, errSubtract
	}

	if difference.IsNeg() {
		return ordering.DoNothing,
			errors.New("price went up")
	}

	if errPercent := helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: s.PriceBeforeDrop,
			PriceNew: params.PriceNow,
			Delta:    s.PercentDrop,
		},
	); errPercent != nil {
		return ordering.DoNothing,
			errors.New("no sudden drop")
	}

	return ordering.Buy, nil
}

func (s *StrategyDropSudden) IncrementSimultaneousOrders() {
	atomic.AddUint32(s.numberSimultaneousOrders, 1)
}

func (s *StrategyDropSudden) DecrementSimultaneousOrders() {
	atomic.AddUint32(s.numberSimultaneousOrders, ^uint32(0))
}

func (s *StrategyDropSudden) CanPlaceOrder() bool {
	return s.allowedNumberSimultaneousOrders >= uint16(*s.numberSimultaneousOrders)
}

func (s StrategyDropSudden) String() string {
	return strings.Join(
		[]string{
			"",
			fmt.Sprintf("Strategy '%s'", nameStrategyDropSudden),
			fmt.Sprintf("Current orders trigerred: %d.", *s.numberSimultaneousOrders),
			fmt.Sprintf("Drop sudden percent: %s%%.", s.PercentDrop.String()),
			fmt.Sprintf(
				"Current average price for %s period: %s.",
				timeperiod.NamePeriodMedium,
				s.PriceBeforeDrop.String(),
			),
			"",
		},
		"\n",
	)
}
