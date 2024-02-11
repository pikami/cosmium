package memoryexecutor

import (
	"github.com/pikami/cosmium/parsers"
)

type RowType interface{}
type ExpressionType interface{}

func Execute(query parsers.SelectStmt, data []RowType) []RowType {
	result := make([]RowType, 0)

	// Iterate over each row in the data
	for _, row := range data {
		// Check if the row satisfies the filter conditions
		if evaluateFilters(query.Filters, row) {
			// Construct a new row based on the selected columns
			newRow := make(map[string]interface{})
			for _, column := range query.Columns {
				destinationName := column.Alias
				if destinationName == "" {
					destinationName = column.Path[len(column.Path)-1]
				}

				newRow[destinationName] = getFieldValue(column, row)
			}
			// Add the new row to the result
			result = append(result, newRow)
		}
	}

	return result
}

// Helper function to evaluate filter conditions recursively
func evaluateFilters(expr ExpressionType, row RowType) bool {
	if expr == nil {
		return true
	}

	switch typedValue := expr.(type) {
	case parsers.ComparisonExpression:
		leftValue := getExpressionParameterValue(typedValue.Left, row)
		rightValue := getExpressionParameterValue(typedValue.Right, row)

		switch typedValue.Operation {
		case "=":
			return leftValue == rightValue
			// Handle other comparison operators as needed
		}
	case parsers.LogicalExpression:
		var result bool
		for i, expression := range typedValue.Expressions {
			expressionResult := evaluateFilters(expression, row)
			if i == 0 {
				result = expressionResult
			}

			switch typedValue.Operation {
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
	return false
}

func getFieldValue(field parsers.FieldPath, row RowType) interface{} {
	value := row
	for _, pathSegment := range field.Path[1:] {
		if nestedValue, ok := value.(map[string]interface{}); ok {
			value = nestedValue[pathSegment]
		} else {
			return nil
		}
	}
	return value
}

func getExpressionParameterValue(parameter interface{}, row RowType) interface{} {
	switch typedParameter := parameter.(type) {
	case parsers.FieldPath:
		return getFieldValue(typedParameter, row)
	case parsers.Constant:
		return typedParameter.Value
	}
	// TODO: Handle error
	return nil
}
