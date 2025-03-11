package memoryexecutor

import "github.com/pikami/cosmium/internal/datastore"

type limitIterator struct {
	documents rowTypeIterator
	limit     int
	count     int
}

func (li *limitIterator) Next() (RowType, datastore.DataStoreStatus) {
	if li.count >= li.limit {
		li.documents = nil
		return rowContext{}, datastore.IterEOF
	}

	li.count++
	return li.documents.Next()
}
