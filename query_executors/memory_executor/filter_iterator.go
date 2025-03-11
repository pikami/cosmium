package memoryexecutor

import (
	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

type filterIterator struct {
	documents rowIterator
	filters   interface{}
}

func (fi *filterIterator) Next() (rowContext, datastore.DataStoreStatus) {
	if fi.documents == nil {
		return rowContext{}, datastore.IterEOF
	}

	for {
		row, status := fi.documents.Next()
		if status != datastore.StatusOk {
			fi.documents = nil
			return rowContext{}, status
		}

		if fi.evaluateFilters(row) {
			return row, status
		}
	}
}

func (fi *filterIterator) evaluateFilters(row rowContext) bool {
	if fi.filters == nil {
		return true
	}

	switch typedFilters := fi.filters.(type) {
	case parsers.ComparisonExpression:
		return row.filters_ComparisonExpression(typedFilters)
	case parsers.LogicalExpression:
		return row.filters_LogicalExpression(typedFilters)
	case parsers.Constant:
		if value, ok := typedFilters.Value.(bool); ok {
			return value
		}
		return false
	case parsers.SelectItem:
		resolvedValue := row.resolveSelectItem(typedFilters)
		if value, ok := resolvedValue.(bool); ok {
			if typedFilters.Invert {
				return !value
			}

			return value
		}
	}

	return false
}

func (r rowContext) applyFilters(filters interface{}) bool {
	if filters == nil {
		return true
	}

	switch typedFilters := filters.(type) {
	case parsers.ComparisonExpression:
		return r.filters_ComparisonExpression(typedFilters)
	case parsers.LogicalExpression:
		return r.filters_LogicalExpression(typedFilters)
	case parsers.Constant:
		if value, ok := typedFilters.Value.(bool); ok {
			return value
		}
		return false
	case parsers.SelectItem:
		resolvedValue := r.resolveSelectItem(typedFilters)
		if value, ok := resolvedValue.(bool); ok {
			if typedFilters.Invert {
				return !value
			}

			return value
		}
	}

	return false
}

func (r rowContext) filters_ComparisonExpression(expression parsers.ComparisonExpression) bool {
	leftExpression, leftExpressionOk := expression.Left.(parsers.SelectItem)
	rightExpression, rightExpressionOk := expression.Right.(parsers.SelectItem)

	if !leftExpressionOk || !rightExpressionOk {
		logger.ErrorLn("ComparisonExpression has incorrect Left or Right type")
		return false
	}

	leftValue := r.resolveSelectItem(leftExpression)
	rightValue := r.resolveSelectItem(rightExpression)

	cmp := compareValues(leftValue, rightValue)
	switch expression.Operation {
	case "=":
		return cmp == 0
	case "!=":
		return cmp != 0
	case "<":
		return cmp < 0
	case ">":
		return cmp > 0
	case "<=":
		return cmp <= 0
	case ">=":
		return cmp >= 0
	}

	return false
}

func (r rowContext) filters_LogicalExpression(expression parsers.LogicalExpression) bool {
	var result bool
	for i, subExpression := range expression.Expressions {
		expressionResult := r.applyFilters(subExpression)
		if i == 0 {
			result = expressionResult
		}

		switch expression.Operation {
		case parsers.LogicalExpressionTypeAnd:
			result = result && expressionResult
			if !result {
				return false
			}
		case parsers.LogicalExpressionTypeOr:
			result = result || expressionResult
			if result {
				return true
			}
		}
	}
	return result
}
