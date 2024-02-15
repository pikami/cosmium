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
		if evaluateFilters(query.Filters, query.Parameters, row) {
			result = append(result, selectRow(query.SelectItems, row))
		}
	}

	// Apply result limit
	if query.Count > 0 {
		count := func() int {
			if len(result) < query.Count {
				return len(result)
			}
			return query.Count
		}()
		result = result[:count]
	}

	return result
}

func selectRow(selectItems []parsers.SelectItem, row RowType) interface{} {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return getFieldValue(selectItems[0], row)
	}

	// Construct a new row based on the selected columns
	newRow := make(map[string]interface{})
	for _, column := range selectItems {
		destinationName := column.Alias
		if destinationName == "" {
			destinationName = column.Path[len(column.Path)-1]
		}

		newRow[destinationName] = getFieldValue(column, row)
	}

	return newRow
}

// Helper function to evaluate filter conditions recursively
func evaluateFilters(expr ExpressionType, Parameters map[string]interface{}, row RowType) bool {
	if expr == nil {
		return true
	}

	switch typedValue := expr.(type) {
	case parsers.ComparisonExpression:
		leftValue := getExpressionParameterValue(typedValue.Left, Parameters, row)
		rightValue := getExpressionParameterValue(typedValue.Right, Parameters, row)

		switch typedValue.Operation {
		case "=":
			return leftValue == rightValue
			// Handle other comparison operators as needed
		}
	case parsers.LogicalExpression:
		var result bool
		for i, expression := range typedValue.Expressions {
			expressionResult := evaluateFilters(expression, Parameters, row)
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

func getFieldValue(field parsers.SelectItem, row RowType) interface{} {
	if field.Type == parsers.SelectItemTypeArray {
		arrayValue := make([]interface{}, 0)
		for _, selectItem := range field.SelectItems {
			arrayValue = append(arrayValue, getFieldValue(selectItem, row))
		}
		return arrayValue
	}

	if field.Type == parsers.SelectItemTypeObject {
		objectValue := make(map[string]interface{})
		for _, selectItem := range field.SelectItems {
			objectValue[selectItem.Alias] = getFieldValue(selectItem, row)
		}
		return objectValue
	}

	value := row
	if len(field.Path) > 1 {
		for _, pathSegment := range field.Path[1:] {
			if nestedValue, ok := value.(map[string]interface{}); ok {
				value = nestedValue[pathSegment]
			} else {
				return nil
			}
		}
	}
	return value
}

func getExpressionParameterValue(
	parameter interface{},
	Parameters map[string]interface{},
	row RowType,
) interface{} {
	switch typedParameter := parameter.(type) {
	case parsers.SelectItem:
		return getFieldValue(typedParameter, row)
	case parsers.Constant:
		if typedParameter.Type == parsers.ConstantTypeParameterConstant &&
			Parameters != nil {
			if key, ok := typedParameter.Value.(string); ok {
				return Parameters[key]
			}
		}

		return typedParameter.Value
	}
	// TODO: Handle error
	return nil
}
