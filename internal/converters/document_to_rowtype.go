package converters

import (
	"github.com/pikami/cosmium/internal/datastore"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

type DocumentToRowTypeIterator struct {
	documents datastore.DocumentIterator
}

func NewDocumentToRowTypeIterator(documents datastore.DocumentIterator) *DocumentToRowTypeIterator {
	return &DocumentToRowTypeIterator{
		documents: documents,
	}
}

func (di *DocumentToRowTypeIterator) Next() (memoryexecutor.RowType, datastore.DataStoreStatus) {
	return di.documents.Next()
}
