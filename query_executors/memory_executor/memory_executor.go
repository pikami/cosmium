package memoryexecutor

import (
	"fmt"
	"reflect"
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
		if evaluateFilters(query.Filters, query.Parameters, row) {
			result = append(result, row)
		}
	}

	// Apply order
	if query.OrderExpressions != nil && len(query.OrderExpressions) > 0 {
		orderBy(query.OrderExpressions, query.Parameters, result)
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
		selectedData = append(selectedData, selectRow(query.SelectItems, query.Parameters, row))
	}

	return selectedData
}

func selectRow(selectItems []parsers.SelectItem, queryParameters map[string]interface{}, row RowType) interface{} {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return getFieldValue(selectItems[0], queryParameters, row)
	}

	// Construct a new row based on the selected columns
	newRow := make(map[string]interface{})
	for index, column := range selectItems {
		destinationName := column.Alias
		if destinationName == "" {
			if len(column.Path) > 0 {
				destinationName = column.Path[len(column.Path)-1]
			} else {
				destinationName = fmt.Sprintf("$%d", index+1)
			}
		}

		newRow[destinationName] = getFieldValue(column, queryParameters, row)
	}

	return newRow
}

func evaluateFilters(expr ExpressionType, queryParameters map[string]interface{}, row RowType) bool {
	if expr == nil {
		return true
	}

	switch typedValue := expr.(type) {
	case parsers.ComparisonExpression:
		leftValue := getExpressionParameterValue(typedValue.Left, queryParameters, row)
		rightValue := getExpressionParameterValue(typedValue.Right, queryParameters, row)

		cmp := compareValues(leftValue, rightValue)
		switch typedValue.Operation {
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
	case parsers.LogicalExpression:
		var result bool
		for i, expression := range typedValue.Expressions {
			expressionResult := evaluateFilters(expression, queryParameters, row)
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
		return false
	case parsers.SelectItem:
		resolvedValue := getFieldValue(typedValue, queryParameters, row)
		if value, ok := resolvedValue.(bool); ok {
			return value
		}
	}
	return false
}

func getFieldValue(field parsers.SelectItem, queryParameters map[string]interface{}, row RowType) interface{} {
	if field.Type == parsers.SelectItemTypeArray {
		arrayValue := make([]interface{}, 0)
		for _, selectItem := range field.SelectItems {
			arrayValue = append(arrayValue, getFieldValue(selectItem, queryParameters, row))
		}
		return arrayValue
	}

	if field.Type == parsers.SelectItemTypeObject {
		objectValue := make(map[string]interface{})
		for _, selectItem := range field.SelectItems {
			objectValue[selectItem.Alias] = getFieldValue(selectItem, queryParameters, row)
		}
		return objectValue
	}

	if field.Type == parsers.SelectItemTypeConstant {
		var typedValue parsers.Constant
		var ok bool
		if typedValue, ok = field.Value.(parsers.Constant); !ok {
			// TODO: Handle error
			fmt.Println("parsers.Constant has incorrect Value type")
		}

		if typedValue.Type == parsers.ConstantTypeParameterConstant &&
			queryParameters != nil {
			if key, ok := typedValue.Value.(string); ok {
				return queryParameters[key]
			}
		}

		return typedValue.Value
	}

	if field.Type == parsers.SelectItemTypeFunctionCall {
		var typedValue parsers.FunctionCall
		var ok bool
		if typedValue, ok = field.Value.(parsers.FunctionCall); !ok {
			// TODO: Handle error
			fmt.Println("parsers.Constant has incorrect Value type")
		}

		switch typedValue.Type {
		case parsers.FunctionCallStringEquals:
			return strings_StringEquals(typedValue.Arguments, queryParameters, row)
		case parsers.FunctionCallContains:
			return strings_Contains(typedValue.Arguments, queryParameters, row)
		case parsers.FunctionCallConcat:
			return strings_Concat(typedValue.Arguments, queryParameters, row)
		case parsers.FunctionCallIsDefined:
			return typeChecking_IsDefined(typedValue.Arguments, queryParameters, row)
		}
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
	queryParameters map[string]interface{},
	row RowType,
) interface{} {
	switch typedParameter := parameter.(type) {
	case parsers.SelectItem:
		return getFieldValue(typedParameter, queryParameters, row)
	}

	fmt.Println("getExpressionParameterValue - got incorrect parameter type")

	return nil
}

func orderBy(orderBy []parsers.OrderExpression, queryParameters map[string]interface{}, data []RowType) {
	less := func(i, j int) bool {
		for _, order := range orderBy {
			val1 := getFieldValue(order.SelectItem, queryParameters, data[i])
			val2 := getFieldValue(order.SelectItem, queryParameters, data[j])

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
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return 1
	}

	switch val1 := val1.(type) {
	case int:
		val2 := val2.(int)
		if val1 < val2 {
			return -1
		} else if val1 > val2 {
			return 1
		}
		return 0
	case float64:
		val2 := val2.(float64)
		if val1 < val2 {
			return -1
		} else if val1 > val2 {
			return 1
		}
		return 0
	case string:
		val2 := val2.(string)
		return strings.Compare(val1, val2)
	case bool:
		val2 := val2.(bool)
		if val1 == val2 {
			return 0
		} else if val1 {
			return 1
		} else {
			return -1
		}
	// TODO: Add more types
	default:
		if reflect.DeepEqual(val1, val2) {
			return 0
		}
		return 1
	}
}
