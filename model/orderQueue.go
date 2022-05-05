package model

import (
	"container/list"
	"matching/enum"
)

type OrderQueue struct {
	sortBy     enum.SortDirection
	parentList *list.List
	elementMap map[string]*list.Element
}