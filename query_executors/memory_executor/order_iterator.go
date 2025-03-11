package memoryexecutor

import (
	"sort"

	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type orderIterator struct {
	documents        rowIterator
	orderExpressions []parsers.OrderExpression
	orderedDocs      []rowContext
	docsIndex        int
}

func (oi *orderIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if oi.orderedDocs != nil {
		if oi.docsIndex >= len(oi.orderedDocs) {
			return rowContext{}, datastore.IterEOF
		}
		row := oi.orderedDocs[oi.docsIndex]
		oi.orderedDocs[oi.docsIndex] = rowContext{}
		oi.docsIndex++
		return row, datastore.StatusOk
	}

	oi.orderedDocs = make([]rowContext, 0)
	for {
		row, status := oi.documents.Next()
		if status != datastore.StatusOk {
			break
		}

		oi.orderedDocs = append(oi.orderedDocs, row)
	}
	oi.documents = nil

	less := func(i, j int) bool {
		for _, order := range oi.orderExpressions {
			val1 := oi.orderedDocs[i].resolveSelectItem(order.SelectItem)
			val2 := oi.orderedDocs[j].resolveSelectItem(order.SelectItem)

			cmp := compareValues(val1, val2)
			if cmp != 0 {
				if order.Direction == parsers.OrderDirectionDesc {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return i < j
	}

	sort.SliceStable(oi.orderedDocs, less)

	if len(oi.orderedDocs) == 0 {
		return rowContext{}, datastore.IterEOF
	}

	oi.docsIndex = 1
	return oi.orderedDocs[0], datastore.StatusOk
}
