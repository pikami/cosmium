package memoryexecutor

import "github.com/pikami/cosmium/internal/datastore"

type offsetIterator struct {
	documents rowTypeIterator
	offset    int
	skipped   bool
}

func (oi *offsetIterator) Next() (RowType, datastore.DataStoreStatus) {
	if oi.skipped {
		return oi.documents.Next()
	}

	for i := 0; i < oi.offset; i++ {
		oi.documents.Next()
	}

	oi.skipped = true
	return oi.Next()
}
