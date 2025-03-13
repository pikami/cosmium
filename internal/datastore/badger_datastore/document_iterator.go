package badgerdatastore

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/vmihailenco/msgpack/v5"
)

type BadgerDocumentIterator struct {
	txn    *badger.Txn
	it     *badger.Iterator
	prefix string
}

func NewBadgerDocumentIterator(txn *badger.Txn, prefix string) *BadgerDocumentIterator {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(prefix)

	it := txn.NewIterator(opts)
	it.Rewind()

	return &BadgerDocumentIterator{
		txn:    txn,
		it:     it,
		prefix: prefix,
	}
}

func (i *BadgerDocumentIterator) Next() (datastore.Document, datastore.DataStoreStatus) {
	if !i.it.Valid() {
		i.it.Close()
		return datastore.Document{}, datastore.IterEOF
	}

	item := i.it.Item()
	val, err := item.ValueCopy(nil)
	if err != nil {
		logger.ErrorLn("Error while copying value:", err)
		return datastore.Document{}, datastore.Unknown
	}

	current := &datastore.Document{}
	err = msgpack.Unmarshal(val, &current)
	if err != nil {
		logger.ErrorLn("Error while decoding value:", err)
		return datastore.Document{}, datastore.Unknown
	}

	i.it.Next()

	return *current, datastore.StatusOk
}

func (i *BadgerDocumentIterator) Close() {
	i.it.Close()
	i.txn.Discard()
}
