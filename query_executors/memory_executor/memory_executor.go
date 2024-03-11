package memoryexecutor

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

type RowType interface{}
type ExpressionType interface{}

type memoryExecutorContext struct {
	parameters map[string]interface{}
}

func Execute(query parsers.SelectStmt, data []RowType) []RowType {
	ctx := memoryExecutorContext{
		parameters: query.Parameters,
	}

	result := make([]RowType, 0)

	// Apply Filter
	for _, row := range data {
		if ctx.evaluateFilters(query.Filters, row) {
			result = append(result, row)
		}
	}

	// Apply order
	if query.OrderExpressions != nil && len(query.OrderExpressions) > 0 {
		ctx.orderBy(query.OrderExpressions, result)
	}

	// Apply group
	isGroupSelect := query.GroupBy != nil && len(query.GroupBy) > 0
	if isGroupSelect {
		result = ctx.groupBy(query, result)
	}

	// Apply select
	if !isGroupSelect {
		selectedData := make([]RowType, 0)
		for _, row := range result {
			selectedData = append(selectedData, ctx.selectRow(query.SelectItems, row))
		}
		result = selectedData
	}

	// Apply distinct
	if query.Distinct {
		result = deduplicate(result)
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

func (c memoryExecutorContext) selectRow(selectItems []parsers.SelectItem, row RowType) interface{} {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return c.getFieldValue(selectItems[0], row)
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

		newRow[destinationName] = c.getFieldValue(column, row)
	}

	return newRow
}

func (c memoryExecutorContext) evaluateFilters(expr ExpressionType, row RowType) bool {
	if expr == nil {
		return true
	}

	switch typedValue := expr.(type) {
	case parsers.ComparisonExpression:
		leftValue := c.getExpressionParameterValue(typedValue.Left, row)
		rightValue := c.getExpressionParameterValue(typedValue.Right, row)

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
			expressionResult := c.evaluateFilters(expression, row)
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
		resolvedValue := c.getFieldValue(typedValue, row)
		if value, ok := resolvedValue.(bool); ok {
			return value
		}
	}
	return false
}

func (c memoryExecutorContext) getFieldValue(field parsers.SelectItem, row RowType) interface{} {
	if field.Type == parsers.SelectItemTypeArray {
		arrayValue := make([]interface{}, 0)
		for _, selectItem := range field.SelectItems {
			arrayValue = append(arrayValue, c.getFieldValue(selectItem, row))
		}
		return arrayValue
	}

	if field.Type == parsers.SelectItemTypeObject {
		objectValue := make(map[string]interface{})
		for _, selectItem := range field.SelectItems {
			objectValue[selectItem.Alias] = c.getFieldValue(selectItem, row)
		}
		return objectValue
	}

	if field.Type == parsers.SelectItemTypeConstant {
		var typedValue parsers.Constant
		var ok bool
		if typedValue, ok = field.Value.(parsers.Constant); !ok {
			// TODO: Handle error
			logger.Error("parsers.Constant has incorrect Value type")
		}

		if typedValue.Type == parsers.ConstantTypeParameterConstant &&
			c.parameters != nil {
			if key, ok := typedValue.Value.(string); ok {
				return c.parameters[key]
			}
		}

		return typedValue.Value
	}

	rowValue := row
	if array, isArray := row.([]RowType); isArray {
		rowValue = array[0]
	}

	if field.Type == parsers.SelectItemTypeFunctionCall {
		var typedValue parsers.FunctionCall
		var ok bool
		if typedValue, ok = field.Value.(parsers.FunctionCall); !ok {
			// TODO: Handle error
			logger.Error("parsers.Constant has incorrect Value type")
		}

		switch typedValue.Type {
		case parsers.FunctionCallStringEquals:
			return c.strings_StringEquals(typedValue.Arguments, rowValue)
		case parsers.FunctionCallContains:
			return c.strings_Contains(typedValue.Arguments, rowValue)
		case parsers.FunctionCallEndsWith:
			return c.strings_EndsWith(typedValue.Arguments, rowValue)
		case parsers.FunctionCallStartsWith:
			return c.strings_StartsWith(typedValue.Arguments, rowValue)
		case parsers.FunctionCallConcat:
			return c.strings_Concat(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIndexOf:
			return c.strings_IndexOf(typedValue.Arguments, rowValue)
		case parsers.FunctionCallToString:
			return c.strings_ToString(typedValue.Arguments, rowValue)
		case parsers.FunctionCallUpper:
			return c.strings_Upper(typedValue.Arguments, rowValue)
		case parsers.FunctionCallLower:
			return c.strings_Lower(typedValue.Arguments, rowValue)
		case parsers.FunctionCallLeft:
			return c.strings_Left(typedValue.Arguments, rowValue)
		case parsers.FunctionCallLength:
			return c.strings_Length(typedValue.Arguments, rowValue)
		case parsers.FunctionCallLTrim:
			return c.strings_LTrim(typedValue.Arguments, rowValue)
		case parsers.FunctionCallReplace:
			return c.strings_Replace(typedValue.Arguments, rowValue)
		case parsers.FunctionCallReplicate:
			return c.strings_Replicate(typedValue.Arguments, rowValue)
		case parsers.FunctionCallReverse:
			return c.strings_Reverse(typedValue.Arguments, rowValue)
		case parsers.FunctionCallRight:
			return c.strings_Right(typedValue.Arguments, rowValue)
		case parsers.FunctionCallRTrim:
			return c.strings_RTrim(typedValue.Arguments, rowValue)
		case parsers.FunctionCallSubstring:
			return c.strings_Substring(typedValue.Arguments, rowValue)
		case parsers.FunctionCallTrim:
			return c.strings_Trim(typedValue.Arguments, rowValue)

		case parsers.FunctionCallIsDefined:
			return c.typeChecking_IsDefined(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsArray:
			return c.typeChecking_IsArray(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsBool:
			return c.typeChecking_IsBool(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsFiniteNumber:
			return c.typeChecking_IsFiniteNumber(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsInteger:
			return c.typeChecking_IsInteger(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsNull:
			return c.typeChecking_IsNull(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsNumber:
			return c.typeChecking_IsNumber(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsObject:
			return c.typeChecking_IsObject(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsPrimitive:
			return c.typeChecking_IsPrimitive(typedValue.Arguments, rowValue)
		case parsers.FunctionCallIsString:
			return c.typeChecking_IsString(typedValue.Arguments, rowValue)

		case parsers.FunctionCallArrayConcat:
			return c.array_Concat(typedValue.Arguments, rowValue)
		case parsers.FunctionCallArrayLength:
			return c.array_Length(typedValue.Arguments, rowValue)
		case parsers.FunctionCallArraySlice:
			return c.array_Slice(typedValue.Arguments, rowValue)
		case parsers.FunctionCallSetIntersect:
			return c.set_Intersect(typedValue.Arguments, rowValue)
		case parsers.FunctionCallSetUnion:
			return c.set_Union(typedValue.Arguments, rowValue)

		case parsers.FunctionCallIn:
			return c.misc_In(typedValue.Arguments, rowValue)
		}
	}

	value := rowValue

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

func (c memoryExecutorContext) getExpressionParameterValue(
	parameter interface{},
	row RowType,
) interface{} {
	switch typedParameter := parameter.(type) {
	case parsers.SelectItem:
		return c.getFieldValue(typedParameter, row)
	}

	logger.Error("getExpressionParameterValue - got incorrect parameter type")

	return nil
}

func (c memoryExecutorContext) orderBy(orderBy []parsers.OrderExpression, data []RowType) {
	less := func(i, j int) bool {
		for _, order := range orderBy {
			val1 := c.getFieldValue(order.SelectItem, data[i])
			val2 := c.getFieldValue(order.SelectItem, data[j])

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

func (c memoryExecutorContext) groupBy(selectStmt parsers.SelectStmt, data []RowType) []RowType {
	groupedRows := make(map[string][]RowType)
	groupedKeys := make([]string, 0)

	// Group rows by group by columns
	for _, row := range data {
		key := c.generateGroupKey(selectStmt.GroupBy, row)
		if _, ok := groupedRows[key]; !ok {
			groupedKeys = append(groupedKeys, key)
		}
		groupedRows[key] = append(groupedRows[key], row)
	}

	// Aggregate each group
	aggregatedRows := make([]RowType, 0)
	for _, key := range groupedKeys {
		groupRows := groupedRows[key]
		aggregatedRow := c.aggregateGroup(selectStmt, groupRows)
		aggregatedRows = append(aggregatedRows, aggregatedRow)
	}

	return aggregatedRows
}

func (c memoryExecutorContext) generateGroupKey(groupByFields []parsers.SelectItem, row RowType) string {
	var keyBuilder strings.Builder
	for _, column := range groupByFields {
		fieldValue := c.getFieldValue(column, row)
		keyBuilder.WriteString(fmt.Sprintf("%v", fieldValue))
		keyBuilder.WriteString(":")
	}

	return keyBuilder.String()
}

func (c memoryExecutorContext) aggregateGroup(selectStmt parsers.SelectStmt, groupRows []RowType) RowType {
	aggregatedRow := c.selectRow(selectStmt.SelectItems, groupRows)

	return aggregatedRow
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

func deduplicate(slice []RowType) []RowType {
	var result []RowType

	for i := 0; i < len(slice); i++ {
		unique := true
		for j := 0; j < len(result); j++ {
			if compareValues(slice[i], result[j]) == 0 {
				unique = false
				break
			}
		}

		if unique {
			result = append(result, slice[i])
		}
	}

	return result
}
