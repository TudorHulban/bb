package strategies

import "test/ordering"

type IStrategy interface {
	AddPriceChange() (ordering.Order, error)
}
