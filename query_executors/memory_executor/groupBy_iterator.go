package memoryexecutor

import (
	"fmt"
	"strings"

	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type groupByIterator struct {
	documents   rowIterator
	groupBy     []parsers.SelectItem
	groupedRows []rowContext
}

func (gi *groupByIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if gi.groupedRows != nil {
		if len(gi.groupedRows) == 0 {
			return rowContext{}, datastore.IterEOF
		}
		row := gi.groupedRows[0]
		gi.groupedRows = gi.groupedRows[1:]
		return row, datastore.StatusOk
	}

	documents := make([]rowContext, 0)
	for {
		row, status := gi.documents.Next()
		if status != datastore.StatusOk {
			break
		}

		documents = append(documents, row)
	}
	gi.documents = nil

	groupedRows := make(map[string][]rowContext)
	groupedKeys := make([]string, 0)

	for _, row := range documents {
		key := row.generateGroupByKey(gi.groupBy)
		if _, ok := groupedRows[key]; !ok {
			groupedKeys = append(groupedKeys, key)
		}
		groupedRows[key] = append(groupedRows[key], row)
	}

	gi.groupedRows = make([]rowContext, 0)
	for _, key := range groupedKeys {
		gi.groupedRows = append(gi.groupedRows, rowContext{
			tables:       groupedRows[key][0].tables,
			parameters:   groupedRows[key][0].parameters,
			grouppedRows: groupedRows[key],
		})
	}

	return gi.Next()
}

func (r rowContext) generateGroupByKey(groupBy []parsers.SelectItem) string {
	var keyBuilder strings.Builder
	for _, selectItem := range groupBy {
		value := r.resolveSelectItem(selectItem)
		keyBuilder.WriteString(fmt.Sprintf("%v", value))
		keyBuilder.WriteString(":")
	}
	return keyBuilder.String()
}
