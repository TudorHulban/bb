package ordering

type Action uint8

func (a Action) String() string {
	switch a {
	case DoNothing:
		return "DoNothing"

	case Buy:
		return "Buy"

	case Sell:
		return "Sell"

	default:
		return "Unknown Action"
	}
}

const (
	DoNothing        = 0
	Buy       Action = 1
	Sell      Action = 2
)

type Order struct {
	Action
	Quantity uint32
}
