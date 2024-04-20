package ordering

import (
	"test/configuration"
)

type ParamsOder struct {
	Code     configuration.CoinCODE
	Quantity uint32
}

type IOrdering interface {
	Buy(params ParamsOder) error
	Sell(params ParamsOder) error
}

var _ IOrdering = &OrderingLogOnly{}
