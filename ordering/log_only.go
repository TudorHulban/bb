package ordering

import "log"

type OrderingLogOnly struct{}

func NewOrderingLogOnly() *OrderingLogOnly {
	return &OrderingLogOnly{}
}

func (o OrderingLogOnly) Buy(params ParamsOder) error {
	log.Printf(
		"buy: %s in quantity: %d",
		params.Code,
		params.Quantity,
	)

	return nil
}

func (o OrderingLogOnly) Sell(params ParamsOder) error {
	log.Printf(
		"sell: %s in quantity: %d",
		params.Code,
		params.Quantity,
	)

	return nil
}
