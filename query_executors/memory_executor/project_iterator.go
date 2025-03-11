package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
	"golang.org/x/exp/slices"
)

type projectIterator struct {
	documents   rowIterator
	selectItems []parsers.SelectItem
	groupBy     []parsers.SelectItem
}

func (pi *projectIterator) Next() (RowType, datastore.DataStoreStatus) {
	if pi.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	row, status := pi.documents.Next()
	if status != datastore.StatusOk {
		pi.documents = nil
		return rowContext{}, status
	}

	if hasAggregateFunctions(pi.selectItems) && len(pi.groupBy) == 0 {
		// When can have aggregate functions without GROUP BY clause,
		// we should aggregate all rows in that case.
		allDocuments := []rowContext{row}
		for {
			row, status := pi.documents.Next()
			if status != datastore.StatusOk {
				break
			}

			allDocuments = append(allDocuments, row)
		}

		if len(allDocuments) == 0 {
			return rowContext{}, datastore.IterEOF
		}

		aggRow := rowContext{
			tables:       row.tables,
			parameters:   row.parameters,
			grouppedRows: allDocuments,
		}

		return aggRow.applyProjection(pi.selectItems), datastore.StatusOk
	}

	return row.applyProjection(pi.selectItems), datastore.StatusOk
}

func (r rowContext) applyProjection(selectItems []parsers.SelectItem) RowType {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return r.resolveSelectItem(selectItems[0])
	}

	// Construct a new row based on the selected columns
	row := make(map[string]interface{})
	for index, selectItem := range selectItems {
		destinationName := resolveDestinationColumnName(selectItem, index, r.parameters)

		row[destinationName] = r.resolveSelectItem(selectItem)
	}

	return row
}

func hasAggregateFunctions(selectItems []parsers.SelectItem) bool {
	if selectItems == nil {
		return false
	}

	for _, selectItem := range selectItems {
		if selectItem.Type == parsers.SelectItemTypeFunctionCall {
			if typedValue, ok := selectItem.Value.(parsers.FunctionCall); ok && slices.Contains[[]parsers.FunctionCallType](parsers.AggregateFunctions, typedValue.Type) {
				return true
			}
		}

		if hasAggregateFunctions(selectItem.SelectItems) {
			return true
		}
	}

	return false
}
