package ordering

type Action uint8

const (
	Buy  Action = 1
	Sell Action = 2
)

type Order struct {
	Action
	Quantity uint32
}
