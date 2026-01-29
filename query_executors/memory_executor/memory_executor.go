package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/parsers"
)

type ExecuteQueryResult struct {
	Rows         []RowType
	HasMorePages bool
}

func ExecuteQuery(
	query parsers.SelectStmt,
	documents rowTypeIterator,
	offset int,
	limit int,
) ExecuteQueryResult {
	resultIter := executeQuery(query, &rowTypeToRowContextIterator{documents: documents, query: query})

	result := &ExecuteQueryResult{
		Rows:         make([]RowType, 0),
		HasMorePages: false,
	}

	for i := 0; i < offset; i++ {
		_, status := resultIter.Next()
		if status != datastore.StatusOk {
			break
		}
	}

	for i := 0; i < limit; i++ {
		row, status := resultIter.Next()
		if status != datastore.StatusOk {
			break
		}

		result.Rows = append(result.Rows, row)
	}

	_, status := resultIter.Next()
	if status == datastore.StatusOk {
		result.HasMorePages = true
	}

	return *result
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
