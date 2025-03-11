package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type joinIterator struct {
	documents rowIterator
	query     parsers.SelectStmt
	buffer    []rowContext
}

func (ji *joinIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if ji.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	if len(ji.buffer) > 0 {
		row := ji.buffer[0]
		ji.buffer = ji.buffer[1:]
		return row, datastore.StatusOk
	}

	doc, status := ji.documents.Next()
	if status != datastore.StatusOk {
		ji.documents = nil
		return rowContext{}, status
	}

	ji.buffer = []rowContext{doc}
	for _, joinItem := range ji.query.JoinItems {
		nextDocuments := make([]rowContext, 0)
		for _, row := range ji.buffer {
			joinedItems := row.resolveJoinItemSelect(joinItem.SelectItem)
			for _, joinedItem := range joinedItems {
				tablesCopy := copyMap(row.tables)
				tablesCopy[joinItem.Table.Value] = joinedItem
				nextDocuments = append(nextDocuments, rowContext{
					parameters: row.parameters,
					tables:     tablesCopy,
				})
			}
		}
		ji.buffer = nextDocuments
	}

	return ji.Next()
}

func (r rowContext) resolveJoinItemSelect(selectItem parsers.SelectItem) []RowType {
	if selectItem.Path != nil || selectItem.Type == parsers.SelectItemTypeSubQuery {
		selectValue := r.parseArray(selectItem)
		documents := make([]RowType, len(selectValue))
		for i, newRowData := range selectValue {
			documents[i] = newRowData
		}
		return documents
	}

	return []RowType{}
}
