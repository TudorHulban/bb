package ordering

type OrderingLogOnly struct{}

func NewOrderingLogOnly() *OrderingLogOnly {
	return &OrderingLogOnly{}
}

func (o OrderingLogOnly) Buy() error {
	return nil
}

func (o OrderingLogOnly) Sell() error {
	return nil
}
