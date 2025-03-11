package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type rowTypeToRowContextIterator struct {
	documents rowTypeIterator
	query     parsers.SelectStmt
}

func (di *rowTypeToRowContextIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if di.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	doc, status := di.documents.Next()
	if status != datastore.StatusOk {
		di.documents = nil
		return rowContext{}, status
	}

	var initialTableName string
	if di.query.Table.SelectItem.Type == parsers.SelectItemTypeSubQuery {
		initialTableName = di.query.Table.SelectItem.Value.(parsers.SelectStmt).Table.Value
	}

	if initialTableName == "" {
		initialTableName = di.query.Table.Value
	}

	if initialTableName == "" {
		initialTableName = resolveDestinationColumnName(di.query.Table.SelectItem, 0, di.query.Parameters)
	}

	return rowContext{
		parameters: di.query.Parameters,
		tables: map[string]RowType{
			initialTableName: doc,
			"$root":          doc,
		},
	}, status
}
