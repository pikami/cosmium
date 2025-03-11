package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

func ExecuteQuery(query parsers.SelectStmt, documents rowTypeIterator) []RowType {
	resultIter := executeQuery(query, &rowTypeToRowContextIterator{documents: documents, query: query})
	result := make([]RowType, 0)
	for {
		row, status := resultIter.Next()
		if status != datastore.StatusOk {
			break
		}

		result = append(result, row)
	}
	return result
}

func executeQuery(query parsers.SelectStmt, documents rowIterator) rowTypeIterator {
	// Resolve FROM
	var iter rowIterator = &fromIterator{
		documents: documents,
		table:     query.Table,
	}

	// Apply JOIN
	if len(query.JoinItems) > 0 {
		iter = &joinIterator{
			documents: iter,
			query:     query,
		}
	}

	// Apply WHERE
	if query.Filters != nil {
		iter = &filterIterator{
			documents: iter,
			filters:   query.Filters,
		}
	}

	// Apply ORDER BY
	if len(query.OrderExpressions) > 0 {
		iter = &orderIterator{
			documents:        iter,
			orderExpressions: query.OrderExpressions,
		}
	}

	// Apply GROUP BY
	if len(query.GroupBy) > 0 {
		iter = &groupByIterator{
			documents: iter,
			groupBy:   query.GroupBy,
		}
	}

	// Apply SELECT
	var projectedIterator rowTypeIterator = &projectIterator{
		documents:   iter,
		selectItems: query.SelectItems,
		groupBy:     query.GroupBy,
	}

	// Apply DISTINCT
	if query.Distinct {
		projectedIterator = &distinctIterator{
			documents: projectedIterator,
		}
	}

	// Apply OFFSET
	if query.Offset > 0 {
		projectedIterator = &offsetIterator{
			documents: projectedIterator,
			offset:    query.Offset,
		}
	}

	// Apply LIMIT
	if query.Count > 0 {
		projectedIterator = &limitIterator{
			documents: projectedIterator,
			limit:     query.Count,
		}
	}

	return projectedIterator
}
