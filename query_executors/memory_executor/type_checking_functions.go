package memoryexecutor

import (
	"math"

	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) typeChecking_IsDefined(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	return ex != nil
}

func (r rowContext) typeChecking_IsArray(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isArray := ex.([]interface{})
	return isArray
}

func (r rowContext) typeChecking_IsBool(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isBool := ex.(bool)
	return isBool
}

func (r rowContext) typeChecking_IsFiniteNumber(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch num := ex.(type) {
	case int:
		return true
	case float64:
		return !math.IsInf(num, 0) && !math.IsNaN(num)
	default:
		return false
	}
}

func (r rowContext) typeChecking_IsInteger(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isInt := ex.(int)
	return isInt
}

func (r rowContext) typeChecking_IsNull(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	return ex == nil
}

func (r rowContext) typeChecking_IsNumber(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isFloat := ex.(float64)
	_, isInt := ex.(int)
	return isFloat || isInt
}

func (r rowContext) typeChecking_IsObject(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isObject := ex.(map[string]interface{})
	return isObject
}

func (r rowContext) typeChecking_IsPrimitive(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch ex.(type) {
	case bool, string, float64, int, nil:
		return true
	default:
		return false
	}
}

func (r rowContext) typeChecking_IsString(arguments []interface{}) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	_, isStr := ex.(string)
	return isStr
}
