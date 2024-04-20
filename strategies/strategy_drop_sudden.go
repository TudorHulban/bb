package strategies

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"test/apperrors"
	"test/helpers"
	"test/ordering"

	"github.com/govalues/decimal"
)

const nameStrategy = "Drop Sudden"

type StrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal

	numberSimultaneousOrders        *uint32
	allowedNumberSimultaneousOrders uint16
}

type ParamsNewStrategyDropSudden struct {
	PriceBeforeDrop decimal.Decimal
	DropPercent     decimal.Decimal

	AllowedNumberSimultaneousOrders uint16
}

func NewStrategyDropSudden(params ParamsNewStrategyDropSudden) (*StrategyDropSudden, error) {
	if !params.DropPercent.IsPos() {
		return nil,
			apperrors.ErrorInvalidInput{
				InputName: "DropPercent",
			}
	}

	var zero uint32

	return &StrategyDropSudden{
			PriceBeforeDrop:                 params.PriceBeforeDrop,
			DropPercent:                     params.DropPercent,
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

func (s *StrategyDropSudden) AddPriceChange(params *ParamsAddPriceChange) (ordering.Action, error) {
	difference, errSubtract := s.PriceBeforeDrop.Sub(params.PriceNow)
	if errSubtract != nil {
		return ordering.DoNothing, errSubtract
	}

	if difference == decimal.Zero {
		return ordering.DoNothing,
			errors.New("no price change")
	}

	if errPercent := helpers.PriceChangeByPercent(
		&helpers.ParamsPriceChangeByPercent{
			PriceOld: s.PriceBeforeDrop,
			PriceNew: params.PriceNow,
			Delta:    s.DropPercent,
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
			fmt.Sprintf("Strategy %s", nameStrategy),
			fmt.Sprintf("Current orders trigerred %d", *s.numberSimultaneousOrders),
			fmt.Sprintf("drop percent: %s%%.", s.DropPercent.String()),
			fmt.Sprintf("current price before drop: %s.", s.PriceBeforeDrop.String()),
			"",
		},
		"\n",
	)
}
