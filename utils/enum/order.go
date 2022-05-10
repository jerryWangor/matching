package enum

type OrderAction string

const (
	ActionCreate OrderAction = "create"
	ActionCancel OrderAction = "cancel"
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