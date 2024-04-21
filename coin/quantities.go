package coin

import (
	"sync/atomic"
	"test/apperrors"

	"github.com/asaskevich/govalidator"
)

type quantities struct {
	currentQuantity atomic.Uint32

	quantityBuy     uint32
	quantityMaximum uint32
}

type ParamNewQuantities struct {
	QuantityBuy     uint32 `valid:"required"`
	QuantityMaximum uint32 `valid:"required"`
}

func NewQuantities(params *ParamNewQuantities) (*quantities, error) {
	if _, errVa := govalidator.ValidateStruct(params); errVa != nil {
		return nil,
			apperrors.ErrValidation{
				Caller: "NewQuantities",
				Issue:  errVa,
			}
	}

	return &quantities{
			quantityBuy:     params.QuantityBuy,
			quantityMaximum: params.QuantityMaximum,
		},
		nil
}

func (q *quantities) canBuy() bool {
	return q.currentQuantity.Load()+q.quantityBuy <= q.quantityMaximum
}

func (q *quantities) canSell() bool {
	return q.currentQuantity.Load() > 0
}
