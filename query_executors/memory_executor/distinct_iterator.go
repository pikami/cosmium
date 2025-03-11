package memoryexecutor

import "github.com/pikami/cosmium/internal/datastore"

type distinctIterator struct {
	documents rowTypeIterator
	seenDocs  []RowType
}

func (di *distinctIterator) Next() (RowType, datastore.DataStoreStatus) {
	if di.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	for {
		row, status := di.documents.Next()
		if status != datastore.StatusOk {
			di.documents = nil
			return rowContext{}, status
		}

		if !di.seen(row) {
			di.seenDocs = append(di.seenDocs, row)
			return row, status
		}
	}
}

func (di *distinctIterator) seen(row RowType) bool {
	for _, seenRow := range di.seenDocs {
		if compareValues(seenRow, row) == 0 {
			return true
		}
	}
	return false
}
