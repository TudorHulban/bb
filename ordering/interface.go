package ordering

type IOrdering interface {
	Buy() error
	Sell() error
}

var _ IOrdering = &OrderingLogOnly{}
