package memoryexecutor

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pikami/cosmium/internal/datastore"
	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

type RowType interface{}
type rowContext struct {
	tables       map[string]RowType
	parameters   map[string]interface{}
	grouppedRows []rowContext
}

type rowIterator interface {
	Next() (rowContext, datastore.DataStoreStatus)
}

type rowTypeIterator interface {
	Next() (RowType, datastore.DataStoreStatus)
}

func resolveDestinationColumnName(selectItem parsers.SelectItem, itemIndex int, queryParameters map[string]interface{}) string {
	if selectItem.Alias != "" {
		return selectItem.Alias
	}

	destinationName := fmt.Sprintf("$%d", itemIndex+1)
	if len(selectItem.Path) > 0 {
		destinationName = selectItem.Path[len(selectItem.Path)-1]
	}

	if destinationName[0] == '@' {
		destinationName = queryParameters[destinationName].(string)
	}

	return destinationName
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

	if selectItem.Type == parsers.SelectItemTypeExpression {
		if typedExpression, ok := selectItem.Value.(parsers.ComparisonExpression); ok {
			return r.filters_ComparisonExpression(typedExpression)
		}

		if typedExpression, ok := selectItem.Value.(parsers.LogicalExpression); ok {
			return r.filters_LogicalExpression(typedExpression)
		}

		logger.ErrorLn("parsers.SelectItem has incorrect Value type (expected parsers.ComparisonExpression)")
		return nil
	}

	if selectItem.Type == parsers.SelectItemTypeBinaryExpression {
		if typedSelectItem, ok := selectItem.Value.(parsers.BinaryExpression); ok {
			return r.selectItem_SelectItemTypeBinaryExpression(typedSelectItem)
		}

		logger.ErrorLn("parsers.SelectItem has incorrect Value type (expected parsers.BinaryExpression)")
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
	subQueryResult := executeQuery(
		subQuery,
		NewRowArrayIterator([]rowContext{r}),
	)

	if subQuery.Exists {
		_, status := subQueryResult.Next()
		return status == datastore.StatusOk
	}

	allDocuments := make([]RowType, 0)
	for {
		row, status := subQueryResult.Next()
		if status != datastore.StatusOk {
			break
		}
		allDocuments = append(allDocuments, row)
	}

	return allDocuments
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

	case parsers.FunctionCallIif:
		return r.misc_Iif(functionCall.Arguments)

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

func (r rowContext) selectItem_SelectItemTypeBinaryExpression(binaryExpression parsers.BinaryExpression) interface{} {
	if binaryExpression.Left == nil || binaryExpression.Right == nil {
		logger.Debug("parsers.BinaryExpression has nil Left or Right value")
		return nil
	}

	leftValue := r.resolveSelectItem(binaryExpression.Left.(parsers.SelectItem))
	rightValue := r.resolveSelectItem(binaryExpression.Right.(parsers.SelectItem))

	if leftValue == nil || rightValue == nil {
		return nil
	}

	leftNumber, leftIsNumber := numToFloat64(leftValue)
	rightNumber, rightIsNumber := numToFloat64(rightValue)

	if !leftIsNumber || !rightIsNumber {
		logger.Debug("Binary expression operands are not numbers, returning nil")
		return nil
	}

	switch binaryExpression.Operation {
	case "+":
		return leftNumber + rightNumber
	case "-":
		return leftNumber - rightNumber
	case "*":
		return leftNumber * rightNumber
	case "/":
		if rightNumber == 0 {
			logger.Debug("Division by zero in binary expression")
			return nil
		}
		return leftNumber / rightNumber
	default:
		return nil
	}
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
			case datastore.Document:
				value = nestedValue[pathSegment]
			case map[string]datastore.Document:
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

func compareValues(val1, val2 interface{}) int {
	// Handle nil values
	if val1 == nil && val2 == nil {
		return 0
	} else if val1 == nil {
		return -1
	} else if val2 == nil {
		return 1
	}

	// Handle number values
	val1Number, val1IsNumber := numToFloat64(val1)
	val2Number, val2IsNumber := numToFloat64(val2)
	if val1IsNumber && val2IsNumber {
		if val1Number < val2Number {
			return -1
		} else if val1Number > val2Number {
			return 1
		}
		return 0
	}

	// Handle different types
	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return 1
	}

	switch val1 := val1.(type) {
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

func copyMap[T RowType | []RowType](originalMap map[string]T) map[string]T {
	targetMap := make(map[string]T)

	for k, v := range originalMap {
		targetMap[k] = v
	}

	return targetMap
}
