package enum

type OrderAction string

const (
	OrderCreate OrderAction = "create"
	OrderCancel OrderAction = "cancel"
)