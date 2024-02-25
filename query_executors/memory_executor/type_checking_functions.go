package memoryexecutor

import (
	"math"

	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) typeChecking_IsDefined(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	return ex != nil
}

func (c memoryExecutorContext) typeChecking_IsArray(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isArray := ex.([]interface{})
	return isArray
}

func (c memoryExecutorContext) typeChecking_IsBool(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isBool := ex.(bool)
	return isBool
}

func (c memoryExecutorContext) typeChecking_IsFiniteNumber(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch num := ex.(type) {
	case int:
		return true
	case float64:
		return !math.IsInf(num, 0) && !math.IsNaN(num)
	default:
		return false
	}
}

func (c memoryExecutorContext) typeChecking_IsInteger(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isInt := ex.(int)
	return isInt
}

func (c memoryExecutorContext) typeChecking_IsNull(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	return ex == nil
}

func (c memoryExecutorContext) typeChecking_IsNumber(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isFloat := ex.(float64)
	_, isInt := ex.(int)
	return isFloat || isInt
}

func (c memoryExecutorContext) typeChecking_IsObject(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isObject := ex.(map[string]interface{})
	return isObject
}

func (c memoryExecutorContext) typeChecking_IsPrimitive(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch ex.(type) {
	case bool, string, float64, int, nil:
		return true
	default:
		return false
	}
}

func (c memoryExecutorContext) typeChecking_IsString(arguments []interface{}, row RowType) bool {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	_, isStr := ex.(string)
	return isStr
}
