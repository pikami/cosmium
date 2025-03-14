package jsondatastore

import "github.com/pikami/cosmium/internal/datastore"

type ArrayDocumentIterator struct {
	documents []datastore.Document
	index     int
}

func (i *ArrayDocumentIterator) Next() (datastore.Document, datastore.DataStoreStatus) {
	i.index++
	if i.index >= len(i.documents) {
		return datastore.Document{}, datastore.StatusNotFound
	}

	return i.documents[i.index], datastore.StatusOk
}

func (i *ArrayDocumentIterator) Close() {
	i.documents = []datastore.Document{}
}
