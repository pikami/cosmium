package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type fromIterator struct {
	documents   rowIterator
	table       parsers.Table
	buffer      []rowContext
	bufferIndex int
}

func (fi *fromIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if fi.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	// Return from buffer if available
	if fi.bufferIndex < len(fi.buffer) {
		result := fi.buffer[fi.bufferIndex]
		fi.buffer[fi.bufferIndex] = rowContext{}
		fi.bufferIndex++
		return result, datastore.StatusOk
	}

	// Resolve next row from documents
	row, status := fi.documents.Next()
	if status != datastore.StatusOk {
		fi.documents = nil
		return row, status
	}

	if fi.table.SelectItem.Path != nil || fi.table.SelectItem.Type == parsers.SelectItemTypeSubQuery {
		destinationTableName := fi.table.SelectItem.Alias
		if destinationTableName == "" {
			destinationTableName = fi.table.Value
		}
		if destinationTableName == "" {
			destinationTableName = resolveDestinationColumnName(fi.table.SelectItem, 0, row.parameters)
		}

		if fi.table.IsInSelect || fi.table.SelectItem.Type == parsers.SelectItemTypeSubQuery {
			selectValue := row.parseArray(fi.table.SelectItem)
			rowContexts := make([]rowContext, len(selectValue))
			for i, newRowData := range selectValue {
				rowContexts[i].parameters = row.parameters
				rowContexts[i].tables = copyMap(row.tables)
				rowContexts[i].tables[destinationTableName] = newRowData
			}

			fi.buffer = rowContexts
			fi.bufferIndex = 0
			return fi.Next()
		}

		if len(fi.table.SelectItem.Path) > 0 {
			sourceTableName := fi.table.SelectItem.Path[0]
			sourceTableData := row.tables[sourceTableName]
			if sourceTableData == nil {
				// When source table is not found, assume it's root document
				row.tables[sourceTableName] = row.tables["$root"]
			}
		}

		newRowData := row.resolveSelectItem(fi.table.SelectItem)
		row.tables[destinationTableName] = newRowData
		return row, status
	}

	return row, status
}
