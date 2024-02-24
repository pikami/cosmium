package memoryexecutor

import (
	"math"

	"github.com/pikami/cosmium/parsers"
)

func typeChecking_IsDefined(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	return ex != nil
}

func typeChecking_IsArray(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isArray := ex.([]interface{})
	return isArray
}

func typeChecking_IsBool(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isBool := ex.(bool)
	return isBool
}

func typeChecking_IsFiniteNumber(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	switch num := ex.(type) {
	case int:
		return true
	case float64:
		return !math.IsInf(num, 0) && !math.IsNaN(num)
	default:
		return false
	}
}

func typeChecking_IsInteger(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isInt := ex.(int)
	return isInt
}

func typeChecking_IsNull(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	return ex == nil
}

func typeChecking_IsNumber(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isFloat := ex.(float64)
	_, isInt := ex.(int)
	return isFloat || isInt
}

func typeChecking_IsObject(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isObject := ex.(map[string]interface{})
	return isObject
}

func typeChecking_IsPrimitive(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	switch ex.(type) {
	case bool, string, float64, int, nil:
		return true
	default:
		return false
	}
}

func typeChecking_IsString(arguments []interface{}, queryParameters map[string]interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := getFieldValue(exItem, queryParameters, row)

	_, isStr := ex.(string)
	return isStr
}
