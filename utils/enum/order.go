package enum

type OrderAction int
type OrderType int
type OrderSide int

const (
	ActionCreate OrderAction = iota
	ActionCancel
)

const (
	TypeLimit OrderType = iota
	TypeLimitIoc
	TypeMarket
	TypeMarketTop5
	TypeMarketTop10
	TypeMarketOpponent
)

const (
	SideBuy OrderSide = iota
	SideSell
)

func (o OrderAction) String() string {
	switch o {
	case ActionCreate:
		return "create"
	case ActionCancel:
		return "cancel"
	default:
		return "unknown"
	}
}

func (o OrderAction) Valid() bool {
	if o.String() == "unknown" {
		return false
	}
	return true
}

func (o OrderType) String() string {
	switch o {
	case TypeLimit:
		return "create"
	case TypeLimitIoc:
		return "cancel"
	case TypeMarket:
		return "cancel"
	case TypeMarketTop5:
		return "cancel"
	case TypeMarketTop10:
		return "cancel"
	case TypeMarketOpponent:
		return "cancel"
	default:
		return "unknown"
	}
}

func (o OrderSide) String() string {
	switch o {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	default:
		return "unknown"
	}
}