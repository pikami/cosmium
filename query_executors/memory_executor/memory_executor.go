package memoryexecutor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pikami/cosmium/parsers"
)

type RowType interface{}
type ExpressionType interface{}

func Execute(query parsers.SelectStmt, data []RowType) []RowType {
	result := make([]RowType, 0)

	// Apply Filter
	for _, row := range data {
		// Check if the row satisfies the filter conditions
		if evaluateFilters(query.Filters, query.Parameters, row) {
			result = append(result, row)
		}
	}

	// Apply order
	if query.OrderExpressions != nil && len(query.OrderExpressions) > 0 {
		orderBy(query.OrderExpressions, result)
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

	// Apply select
	selectedData := make([]RowType, 0)
	for _, row := range result {
		selectedData = append(selectedData, selectRow(query.SelectItems, row))
	}

	return selectedData
}

func selectRow(selectItems []parsers.SelectItem, row RowType) interface{} {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return getFieldValue(selectItems[0], row)
	}

	// Construct a new row based on the selected columns
	newRow := make(map[string]interface{})
	for index, column := range selectItems {
		destinationName := column.Alias
		if destinationName == "" {
			if len(column.Path) < 1 {
				destinationName = fmt.Sprintf("$%d", index+1)
			} else {
				destinationName = column.Path[len(column.Path)-1]
			}
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
		case "!=":
			return leftValue != rightValue
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
	case parsers.Constant:
		if value, ok := typedValue.Value.(bool); ok {
			return value
		}
		// TODO: Check if we should do something if it is not a boolean constant
		return false
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

func orderBy(orderBy []parsers.OrderExpression, data []RowType) {
	less := func(i, j int) bool {
		for _, order := range orderBy {
			val1 := getFieldValue(order.SelectItem, data[i])
			val2 := getFieldValue(order.SelectItem, data[j])

			cmp := compareValues(val1, val2)
			if cmp != 0 {
				if order.Direction == parsers.OrderDirectionDesc {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return i < j
	}

	sort.SliceStable(data, less)
}

func compareValues(val1, val2 interface{}) int {
	switch val1 := val1.(type) {
	case int:
		val2 := val2.(int)
		if val1 < val2 {
			return -1
		} else if val1 > val2 {
			return 1
		}
		return 0
	case string:
		val2 := val2.(string)
		return strings.Compare(val1, val2)
	// TODO: Add more types
	default:
		return 0
	}
}
