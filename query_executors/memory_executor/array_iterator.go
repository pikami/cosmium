package memoryexecutor

import "github.com/pikami/cosmium/internal/datastore"

type rowArrayIterator struct {
	documents []rowContext
	index     int
}

func NewRowArrayIterator(documents []rowContext) *rowArrayIterator {
	return &rowArrayIterator{
		documents: documents,
		index:     -1,
	}
}

func (i *rowArrayIterator) Next() (rowContext, datastore.DataStoreStatus) {
	i.index++
	if i.index >= len(i.documents) {
		return rowContext{}, datastore.IterEOF
	}

	row := i.documents[i.index]
	i.documents[i.index] = rowContext{} // Help GC reclaim memory

	return row, datastore.StatusOk
}
