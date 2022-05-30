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
		return "limit"
	case TypeLimitIoc:
		return "limitIoc"
	case TypeMarket:
		return "market"
	case TypeMarketTop5:
		return "marketTop5"
	case TypeMarketTop10:
		return "marketTop10"
	case TypeMarketOpponent:
		return "marketOpponent"
	default:
		return "unknown"
	}
}

func (o OrderType) Valid() bool {
	if o.String() == "unknown" {
		return false
	}
	return true
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

func (o OrderSide) Valid() bool {
	if o.String() == "unknown" {
		return false
	}
	return true
}