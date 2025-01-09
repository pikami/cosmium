package memoryexecutor

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
	"golang.org/x/exp/slices"
)

type RowType interface{}
type rowContext struct {
	tables       map[string]RowType
	parameters   map[string]interface{}
	grouppedRows []rowContext
}

func ExecuteQuery(query parsers.SelectStmt, documents []RowType) []RowType {
	currentDocuments := make([]rowContext, 0)
	for _, doc := range documents {
		currentDocuments = append(currentDocuments, resolveFrom(query, doc)...)
	}

	// Handle JOINS
	nextDocuments := make([]rowContext, 0)
	for _, currentDocument := range currentDocuments {
		rowContexts := currentDocument.handleJoin(query)
		nextDocuments = append(nextDocuments, rowContexts...)
	}
	currentDocuments = nextDocuments

	// Apply filters
	nextDocuments = make([]rowContext, 0)
	for _, currentDocument := range currentDocuments {
		if currentDocument.applyFilters(query.Filters) {
			nextDocuments = append(nextDocuments, currentDocument)
		}
	}
	currentDocuments = nextDocuments

	// Apply order
	if len(query.OrderExpressions) > 0 {
		applyOrder(currentDocuments, query.OrderExpressions)
	}

	// Apply group by
	if len(query.GroupBy) > 0 {
		currentDocuments = applyGroupBy(currentDocuments, query.GroupBy)
	}

	// Apply select
	projectedDocuments := applyProjection(currentDocuments, query.SelectItems, query.GroupBy)

	// Apply distinct
	if query.Distinct {
		projectedDocuments = deduplicate(projectedDocuments)
	}

	// Apply result limit
	if query.Count > 0 && len(projectedDocuments) > query.Count {
		projectedDocuments = projectedDocuments[:query.Count]
	}

	return projectedDocuments
}

func resolveFrom(query parsers.SelectStmt, doc RowType) []rowContext {
	initialRow, gotParentContext := doc.(rowContext)
	if !gotParentContext {
		var initialTableName string
		if query.Table.SelectItem.Type == parsers.SelectItemTypeSubQuery {
			initialTableName = query.Table.SelectItem.Value.(parsers.SelectStmt).Table.Value
		}

		if initialTableName == "" {
			initialTableName = query.Table.Value
		}

		initialRow = rowContext{
			parameters: query.Parameters,
			tables: map[string]RowType{
				initialTableName: doc,
			},
		}
	}

	if query.Table.SelectItem.Path != nil || query.Table.SelectItem.Type == parsers.SelectItemTypeSubQuery {
		destinationTableName := query.Table.SelectItem.Alias
		if destinationTableName == "" {
			destinationTableName = query.Table.Value
		}

		selectValue := initialRow.parseArray(query.Table.SelectItem)
		rowContexts := make([]rowContext, len(selectValue))
		for i, newRowData := range selectValue {
			rowContexts[i].parameters = initialRow.parameters
			rowContexts[i].tables = copyMap(initialRow.tables)
			rowContexts[i].tables[destinationTableName] = newRowData
		}
		return rowContexts
	}

	return []rowContext{initialRow}
}

func (r rowContext) handleJoin(query parsers.SelectStmt) []rowContext {
	currentDocuments := []rowContext{r}

	for _, joinItem := range query.JoinItems {
		nextDocuments := make([]rowContext, 0)
		for _, currentDocument := range currentDocuments {
			joinedItems := currentDocument.resolveJoinItemSelect(joinItem.SelectItem)
			for _, joinedItem := range joinedItems {
				tablesCopy := copyMap(currentDocument.tables)
				tablesCopy[joinItem.Table.Value] = joinedItem
				nextDocuments = append(nextDocuments, rowContext{
					parameters: currentDocument.parameters,
					tables:     tablesCopy,
				})
			}
		}
		currentDocuments = nextDocuments
	}

	return currentDocuments
}

func (r rowContext) resolveJoinItemSelect(selectItem parsers.SelectItem) []RowType {
	if selectItem.Path != nil || selectItem.Type == parsers.SelectItemTypeSubQuery {
		selectValue := r.parseArray(selectItem)
		documents := make([]RowType, len(selectValue))
		for i, newRowData := range selectValue {
			documents[i] = newRowData
		}
		return documents
	}

	return []RowType{}
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

func applyOrder(documents []rowContext, orderExpressions []parsers.OrderExpression) {
	less := func(i, j int) bool {
		for _, order := range orderExpressions {
			val1 := documents[i].resolveSelectItem(order.SelectItem)
			val2 := documents[j].resolveSelectItem(order.SelectItem)

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

	sort.SliceStable(documents, less)
}

func applyGroupBy(documents []rowContext, groupBy []parsers.SelectItem) []rowContext {
	groupedRows := make(map[string][]rowContext)
	groupedKeys := make([]string, 0)

	for _, row := range documents {
		key := row.generateGroupByKey(groupBy)
		if _, ok := groupedRows[key]; !ok {
			groupedKeys = append(groupedKeys, key)
		}
		groupedRows[key] = append(groupedRows[key], row)
	}

	grouppedRows := make([]rowContext, 0)
	for _, key := range groupedKeys {
		grouppedRowContext := rowContext{
			tables:       groupedRows[key][0].tables,
			parameters:   groupedRows[key][0].parameters,
			grouppedRows: groupedRows[key],
		}
		grouppedRows = append(grouppedRows, grouppedRowContext)
	}

	return grouppedRows
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

func applyProjection(documents []rowContext, selectItems []parsers.SelectItem, groupBy []parsers.SelectItem) []RowType {
	if len(documents) == 0 {
		return []RowType{}
	}

	if hasAggregateFunctions(selectItems) && len(groupBy) == 0 {
		// When can have aggregate functions without GROUP BY clause,
		// we should aggregate all rows in that case
		rowContext := rowContext{
			tables:       documents[0].tables,
			parameters:   documents[0].parameters,
			grouppedRows: documents,
		}
		return []RowType{rowContext.applyProjection(selectItems)}
	}

	projectedDocuments := make([]RowType, len(documents))
	for index, row := range documents {
		projectedDocuments[index] = row.applyProjection(selectItems)
	}

	return projectedDocuments
}

func (r rowContext) applyProjection(selectItems []parsers.SelectItem) RowType {
	// When the first value is top level, select it instead
	if len(selectItems) > 0 && selectItems[0].IsTopLevel {
		return r.resolveSelectItem(selectItems[0])
	}

	// Construct a new row based on the selected columns
	row := make(map[string]interface{})
	for index, selectItem := range selectItems {
		destinationName := selectItem.Alias
		if destinationName == "" {
			if len(selectItem.Path) > 0 {
				destinationName = selectItem.Path[len(selectItem.Path)-1]
			} else {
				destinationName = fmt.Sprintf("$%d", index+1)
			}

			if destinationName[0] == '@' {
				destinationName = r.parameters[destinationName].(string)
			}
		}

		row[destinationName] = r.resolveSelectItem(selectItem)
	}

	return row
}

func (r rowContext) resolveSelectItem(selectItem parsers.SelectItem) interface{} {
	if selectItem.Type == parsers.SelectItemTypeArray {
		return r.selectItem_SelectItemTypeArray(selectItem)
	}

	if selectItem.Type == parsers.SelectItemTypeObject {
		return r.selectItem_SelectItemTypeObject(selectItem)
	}

	if selectItem.Type == parsers.SelectItemTypeConstant {
		return r.selectItem_SelectItemTypeConstant(selectItem)
	}

	if selectItem.Type == parsers.SelectItemTypeSubQuery {
		return r.selectItem_SelectItemTypeSubQuery(selectItem)
	}

	if selectItem.Type == parsers.SelectItemTypeFunctionCall {
		if typedFunctionCall, ok := selectItem.Value.(parsers.FunctionCall); ok {
			return r.selectItem_SelectItemTypeFunctionCall(typedFunctionCall)
		}

		logger.ErrorLn("parsers.SelectItem has incorrect Value type (expected parsers.FunctionCall)")
		return nil
	}

	return r.selectItem_SelectItemTypeField(selectItem)
}

func (r rowContext) selectItem_SelectItemTypeArray(selectItem parsers.SelectItem) interface{} {
	arrayValue := make([]interface{}, 0)
	for _, subSelectItem := range selectItem.SelectItems {
		arrayValue = append(arrayValue, r.resolveSelectItem(subSelectItem))
	}
	return arrayValue
}

func (r rowContext) selectItem_SelectItemTypeObject(selectItem parsers.SelectItem) interface{} {
	objectValue := make(map[string]interface{})
	for _, subSelectItem := range selectItem.SelectItems {
		objectValue[subSelectItem.Alias] = r.resolveSelectItem(subSelectItem)
	}
	return objectValue
}

func (r rowContext) selectItem_SelectItemTypeConstant(selectItem parsers.SelectItem) interface{} {
	var typedValue parsers.Constant
	var ok bool
	if typedValue, ok = selectItem.Value.(parsers.Constant); !ok {
		// TODO: Handle error
		logger.ErrorLn("parsers.Constant has incorrect Value type")
	}

	if typedValue.Type == parsers.ConstantTypeParameterConstant &&
		r.parameters != nil {
		if key, ok := typedValue.Value.(string); ok {
			return r.parameters[key]
		}
	}

	return typedValue.Value
}

func (r rowContext) selectItem_SelectItemTypeSubQuery(selectItem parsers.SelectItem) interface{} {
	subQuery := selectItem.Value.(parsers.SelectStmt)
	subQueryResult := ExecuteQuery(
		subQuery,
		[]RowType{r},
	)

	if subQuery.Exists {
		return len(subQueryResult) > 0
	}

	return subQueryResult
}

func (r rowContext) selectItem_SelectItemTypeFunctionCall(functionCall parsers.FunctionCall) interface{} {
	switch functionCall.Type {
	case parsers.FunctionCallStringEquals:
		return r.strings_StringEquals(functionCall.Arguments)
	case parsers.FunctionCallContains:
		return r.strings_Contains(functionCall.Arguments)
	case parsers.FunctionCallEndsWith:
		return r.strings_EndsWith(functionCall.Arguments)
	case parsers.FunctionCallStartsWith:
		return r.strings_StartsWith(functionCall.Arguments)
	case parsers.FunctionCallConcat:
		return r.strings_Concat(functionCall.Arguments)
	case parsers.FunctionCallIndexOf:
		return r.strings_IndexOf(functionCall.Arguments)
	case parsers.FunctionCallToString:
		return r.strings_ToString(functionCall.Arguments)
	case parsers.FunctionCallUpper:
		return r.strings_Upper(functionCall.Arguments)
	case parsers.FunctionCallLower:
		return r.strings_Lower(functionCall.Arguments)
	case parsers.FunctionCallLeft:
		return r.strings_Left(functionCall.Arguments)
	case parsers.FunctionCallLength:
		return r.strings_Length(functionCall.Arguments)
	case parsers.FunctionCallLTrim:
		return r.strings_LTrim(functionCall.Arguments)
	case parsers.FunctionCallReplace:
		return r.strings_Replace(functionCall.Arguments)
	case parsers.FunctionCallReplicate:
		return r.strings_Replicate(functionCall.Arguments)
	case parsers.FunctionCallReverse:
		return r.strings_Reverse(functionCall.Arguments)
	case parsers.FunctionCallRight:
		return r.strings_Right(functionCall.Arguments)
	case parsers.FunctionCallRTrim:
		return r.strings_RTrim(functionCall.Arguments)
	case parsers.FunctionCallSubstring:
		return r.strings_Substring(functionCall.Arguments)
	case parsers.FunctionCallTrim:
		return r.strings_Trim(functionCall.Arguments)

	case parsers.FunctionCallIsDefined:
		return r.typeChecking_IsDefined(functionCall.Arguments)
	case parsers.FunctionCallIsArray:
		return r.typeChecking_IsArray(functionCall.Arguments)
	case parsers.FunctionCallIsBool:
		return r.typeChecking_IsBool(functionCall.Arguments)
	case parsers.FunctionCallIsFiniteNumber:
		return r.typeChecking_IsFiniteNumber(functionCall.Arguments)
	case parsers.FunctionCallIsInteger:
		return r.typeChecking_IsInteger(functionCall.Arguments)
	case parsers.FunctionCallIsNull:
		return r.typeChecking_IsNull(functionCall.Arguments)
	case parsers.FunctionCallIsNumber:
		return r.typeChecking_IsNumber(functionCall.Arguments)
	case parsers.FunctionCallIsObject:
		return r.typeChecking_IsObject(functionCall.Arguments)
	case parsers.FunctionCallIsPrimitive:
		return r.typeChecking_IsPrimitive(functionCall.Arguments)
	case parsers.FunctionCallIsString:
		return r.typeChecking_IsString(functionCall.Arguments)

	case parsers.FunctionCallArrayConcat:
		return r.array_Concat(functionCall.Arguments)
	case parsers.FunctionCallArrayContains:
		return r.array_Contains(functionCall.Arguments)
	case parsers.FunctionCallArrayContainsAny:
		return r.array_Contains_Any(functionCall.Arguments)
	case parsers.FunctionCallArrayContainsAll:
		return r.array_Contains_All(functionCall.Arguments)
	case parsers.FunctionCallArrayLength:
		return r.array_Length(functionCall.Arguments)
	case parsers.FunctionCallArraySlice:
		return r.array_Slice(functionCall.Arguments)
	case parsers.FunctionCallSetIntersect:
		return r.set_Intersect(functionCall.Arguments)
	case parsers.FunctionCallSetUnion:
		return r.set_Union(functionCall.Arguments)

	case parsers.FunctionCallMathAbs:
		return r.math_Abs(functionCall.Arguments)
	case parsers.FunctionCallMathAcos:
		return r.math_Acos(functionCall.Arguments)
	case parsers.FunctionCallMathAsin:
		return r.math_Asin(functionCall.Arguments)
	case parsers.FunctionCallMathAtan:
		return r.math_Atan(functionCall.Arguments)
	case parsers.FunctionCallMathCeiling:
		return r.math_Ceiling(functionCall.Arguments)
	case parsers.FunctionCallMathCos:
		return r.math_Cos(functionCall.Arguments)
	case parsers.FunctionCallMathCot:
		return r.math_Cot(functionCall.Arguments)
	case parsers.FunctionCallMathDegrees:
		return r.math_Degrees(functionCall.Arguments)
	case parsers.FunctionCallMathExp:
		return r.math_Exp(functionCall.Arguments)
	case parsers.FunctionCallMathFloor:
		return r.math_Floor(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitNot:
		return r.math_IntBitNot(functionCall.Arguments)
	case parsers.FunctionCallMathLog10:
		return r.math_Log10(functionCall.Arguments)
	case parsers.FunctionCallMathRadians:
		return r.math_Radians(functionCall.Arguments)
	case parsers.FunctionCallMathRound:
		return r.math_Round(functionCall.Arguments)
	case parsers.FunctionCallMathSign:
		return r.math_Sign(functionCall.Arguments)
	case parsers.FunctionCallMathSin:
		return r.math_Sin(functionCall.Arguments)
	case parsers.FunctionCallMathSqrt:
		return r.math_Sqrt(functionCall.Arguments)
	case parsers.FunctionCallMathSquare:
		return r.math_Square(functionCall.Arguments)
	case parsers.FunctionCallMathTan:
		return r.math_Tan(functionCall.Arguments)
	case parsers.FunctionCallMathTrunc:
		return r.math_Trunc(functionCall.Arguments)
	case parsers.FunctionCallMathAtn2:
		return r.math_Atn2(functionCall.Arguments)
	case parsers.FunctionCallMathIntAdd:
		return r.math_IntAdd(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitAnd:
		return r.math_IntBitAnd(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitLeftShift:
		return r.math_IntBitLeftShift(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitOr:
		return r.math_IntBitOr(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitRightShift:
		return r.math_IntBitRightShift(functionCall.Arguments)
	case parsers.FunctionCallMathIntBitXor:
		return r.math_IntBitXor(functionCall.Arguments)
	case parsers.FunctionCallMathIntDiv:
		return r.math_IntDiv(functionCall.Arguments)
	case parsers.FunctionCallMathIntMod:
		return r.math_IntMod(functionCall.Arguments)
	case parsers.FunctionCallMathIntMul:
		return r.math_IntMul(functionCall.Arguments)
	case parsers.FunctionCallMathIntSub:
		return r.math_IntSub(functionCall.Arguments)
	case parsers.FunctionCallMathPower:
		return r.math_Power(functionCall.Arguments)
	case parsers.FunctionCallMathLog:
		return r.math_Log(functionCall.Arguments)
	case parsers.FunctionCallMathNumberBin:
		return r.math_NumberBin(functionCall.Arguments)
	case parsers.FunctionCallMathPi:
		return r.math_Pi()
	case parsers.FunctionCallMathRand:
		return r.math_Rand()

	case parsers.FunctionCallAggregateAvg:
		return r.aggregate_Avg(functionCall.Arguments)
	case parsers.FunctionCallAggregateCount:
		return r.aggregate_Count(functionCall.Arguments)
	case parsers.FunctionCallAggregateMax:
		return r.aggregate_Max(functionCall.Arguments)
	case parsers.FunctionCallAggregateMin:
		return r.aggregate_Min(functionCall.Arguments)
	case parsers.FunctionCallAggregateSum:
		return r.aggregate_Sum(functionCall.Arguments)

	case parsers.FunctionCallIn:
		return r.misc_In(functionCall.Arguments)
	}

	logger.Errorf("Unknown function call type: %v", functionCall.Type)
	return nil
}

func (r rowContext) selectItem_SelectItemTypeField(selectItem parsers.SelectItem) interface{} {
	value := r.tables[selectItem.Path[0]]

	if len(selectItem.Path) > 1 {
		for _, pathSegment := range selectItem.Path[1:] {
			if pathSegment[0] == '@' {
				pathSegment = r.parameters[pathSegment].(string)
			}

			switch nestedValue := value.(type) {
			case map[string]interface{}:
				value = nestedValue[pathSegment]
			case map[string]RowType:
				value = nestedValue[pathSegment]
			case []int, []string, []interface{}:
				slice := reflect.ValueOf(nestedValue)
				if arrayIndex, err := strconv.Atoi(pathSegment); err == nil && slice.Len() > arrayIndex {
					value = slice.Index(arrayIndex).Interface()
				} else {
					return nil
				}
			default:
				return nil
			}
		}
	}

	return value
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

func deduplicate[T RowType | interface{}](slice []T) []T {
	var result []T
	result = make([]T, 0)

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

func copyMap[T RowType | []RowType](originalMap map[string]T) map[string]T {
	targetMap := make(map[string]T)

	for k, v := range originalMap {
		targetMap[k] = v
	}

	return targetMap
}
