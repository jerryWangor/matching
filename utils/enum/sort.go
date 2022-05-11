package enum

// SortDirection 升序=1 倒序=2
type SortDirection int

const (
	SortAsc SortDirection =  iota + 1
	SortDesc
)