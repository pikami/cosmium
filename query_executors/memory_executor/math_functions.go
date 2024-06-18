package memoryexecutor

import (
	"math"
	"math/rand"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

func (c memoryExecutorContext) math_Abs(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		return math.Abs(val)
	case int:
		if val < 0 {
			return -val
		}
		return val
	default:
		logger.Debug("math_Abs - got parameters of wrong type")
		return 0
	}
}

func (c memoryExecutorContext) math_Acos(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Acos - got parameters of wrong type")
		return nil
	}

	if val < -1 || val > 1 {
		logger.Debug("math_Acos - value out of domain for acos")
		return nil
	}

	return math.Acos(val) * 180 / math.Pi
}

func (c memoryExecutorContext) math_Asin(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Asin - got parameters of wrong type")
		return nil
	}

	if val < -1 || val > 1 {
		logger.Debug("math_Asin - value out of domain for acos")
		return nil
	}

	return math.Asin(val) * 180 / math.Pi
}

func (c memoryExecutorContext) math_Atan(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Atan - got parameters of wrong type")
		return nil
	}

	return math.Atan(val) * 180 / math.Pi
}

func (c memoryExecutorContext) math_Ceiling(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		return math.Ceil(val)
	case int:
		return val
	default:
		logger.Debug("math_Ceiling - got parameters of wrong type")
		return 0
	}
}

func (c memoryExecutorContext) math_Cos(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Cos - got parameters of wrong type")
		return nil
	}

	return math.Cos(val)
}

func (c memoryExecutorContext) math_Cot(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Cot - got parameters of wrong type")
		return nil
	}

	if val == 0 {
		logger.Debug("math_Cot - cotangent undefined for zero")
		return nil
	}

	return 1 / math.Tan(val)
}

func (c memoryExecutorContext) math_Degrees(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Degrees - got parameters of wrong type")
		return nil
	}

	return val * (180 / math.Pi)
}

func (c memoryExecutorContext) math_Exp(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Exp - got parameters of wrong type")
		return nil
	}

	return math.Exp(val)
}

func (c memoryExecutorContext) math_Floor(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		return math.Floor(val)
	case int:
		return val
	default:
		logger.Debug("math_Floor - got parameters of wrong type")
		return 0
	}
}

func (c memoryExecutorContext) math_IntBitNot(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case int:
		return ^val
	default:
		logger.Debug("math_IntBitNot - got parameters of wrong type")
		return nil
	}
}

func (c memoryExecutorContext) math_Log10(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Log10 - got parameters of wrong type")
		return nil
	}

	if val <= 0 {
		logger.Debug("math_Log10 - value must be greater than 0")
		return nil
	}

	return math.Log10(val)
}

func (c memoryExecutorContext) math_Radians(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Radians - got parameters of wrong type")
		return nil
	}

	return val * (math.Pi / 180.0)
}

func (c memoryExecutorContext) math_Round(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		return math.Round(val)
	case int:
		return val
	default:
		logger.Debug("math_Round - got parameters of wrong type")
		return nil
	}
}

func (c memoryExecutorContext) math_Sign(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		if val > 0 {
			return 1
		} else if val < 0 {
			return -1
		} else {
			return 0
		}
	case int:
		if val > 0 {
			return 1
		} else if val < 0 {
			return -1
		} else {
			return 0
		}
	default:
		logger.Debug("math_Sign - got parameters of wrong type")
		return nil
	}
}

func (c memoryExecutorContext) math_Sin(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Sin - got parameters of wrong type")
		return nil
	}

	return math.Sin(val)
}

func (c memoryExecutorContext) math_Sqrt(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Sqrt - got parameters of wrong type")
		return nil
	}

	return math.Sqrt(val)
}

func (c memoryExecutorContext) math_Square(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Square - got parameters of wrong type")
		return nil
	}

	return math.Pow(val, 2)
}

func (c memoryExecutorContext) math_Tan(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.Debug("math_Tan - got parameters of wrong type")
		return nil
	}

	return math.Tan(val)
}

func (c memoryExecutorContext) math_Trunc(arguments []interface{}, row RowType) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem, row)

	switch val := ex.(type) {
	case float64:
		return math.Trunc(val)
	case int:
		return float64(val)
	default:
		logger.Debug("math_Trunc - got parameters of wrong type")
		return nil
	}
}

func (c memoryExecutorContext) math_Atn2(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	y, yIsNumber := numToFloat64(ex1)
	x, xIsNumber := numToFloat64(ex2)

	if !yIsNumber || !xIsNumber {
		logger.Debug("math_Atn2 - got parameters of wrong type")
		return nil
	}

	return math.Atan2(y, x)
}

func (c memoryExecutorContext) math_IntAdd(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	ex1Number, ex1IsNumber := numToInt(ex1)
	ex2Number, ex2IsNumber := numToInt(ex2)

	if !ex1IsNumber || !ex2IsNumber {
		logger.Debug("math_IntAdd - got parameters of wrong type")
		return nil
	}

	return ex1Number + ex2Number
}

func (c memoryExecutorContext) math_IntBitAnd(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	ex1Int, ex1IsInt := numToInt(ex1)
	ex2Int, ex2IsInt := numToInt(ex2)

	if !ex1IsInt || !ex2IsInt {
		logger.Debug("math_IntBitAnd - got parameters of wrong type")
		return nil
	}

	return ex1Int & ex2Int
}

func (c memoryExecutorContext) math_IntBitLeftShift(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := numToInt(ex1)
	num2, num2IsInt := numToInt(ex2)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntBitLeftShift - got parameters of wrong type")
		return nil
	}

	return num1 << uint(num2)
}

func (c memoryExecutorContext) math_IntBitOr(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntBitOr - got parameters of wrong type")
		return nil
	}

	return num1 | num2
}

func (c memoryExecutorContext) math_IntBitRightShift(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := numToInt(ex1)
	num2, num2IsInt := numToInt(ex2)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntBitRightShift - got parameters of wrong type")
		return nil
	}

	return num1 >> uint(num2)
}

func (c memoryExecutorContext) math_IntBitXor(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntBitXor - got parameters of wrong type")
		return nil
	}

	return num1 ^ num2
}

func (c memoryExecutorContext) math_IntDiv(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt || num2 == 0 {
		logger.Debug("math_IntDiv - got parameters of wrong type or divide by zero")
		return nil
	}

	return num1 / num2
}

func (c memoryExecutorContext) math_IntMul(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntMul - got parameters of wrong type")
		return nil
	}

	return num1 * num2
}

func (c memoryExecutorContext) math_IntSub(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.Debug("math_IntSub - got parameters of wrong type")
		return nil
	}

	return num1 - num2
}

func (c memoryExecutorContext) math_IntMod(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt || num2 == 0 {
		logger.Debug("math_IntMod - got parameters of wrong type or divide by zero")
		return nil
	}

	return num1 % num2
}

func (c memoryExecutorContext) math_Power(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := c.getFieldValue(exItem1, row)
	ex2 := c.getFieldValue(exItem2, row)

	base, baseIsNumber := numToFloat64(ex1)
	exponent, exponentIsNumber := numToFloat64(ex2)

	if !baseIsNumber || !exponentIsNumber {
		logger.Debug("math_Power - got parameters of wrong type")
		return nil
	}

	return math.Pow(base, exponent)
}

func (c memoryExecutorContext) math_Log(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem1, row)

	var base float64 = math.E
	if len(arguments) > 1 {
		exItem2 := arguments[1].(parsers.SelectItem)
		baseValueObject := c.getFieldValue(exItem2, row)
		baseValue, baseValueIsNumber := numToFloat64(baseValueObject)

		if !baseValueIsNumber {
			logger.Debug("math_Log - base parameter must be a numeric value")
			return nil
		}

		if baseValue > 0 && baseValue != 1 {
			base = baseValue
		} else {
			logger.Debug("math_Log - base must be greater than 0 and not equal to 1")
			return nil
		}
	}

	num, numIsNumber := numToFloat64(ex)
	if !numIsNumber || num <= 0 {
		logger.Debug("math_Log - parameter must be a positive numeric value")
		return nil
	}

	return math.Log(num) / math.Log(base)
}

func (c memoryExecutorContext) math_NumberBin(arguments []interface{}, row RowType) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	ex := c.getFieldValue(exItem1, row)

	binSize := 1.0

	if len(arguments) > 1 {
		exItem2 := arguments[1].(parsers.SelectItem)
		binSizeValueObject := c.getFieldValue(exItem2, row)
		binSizeValue, binSizeValueIsNumber := numToFloat64(binSizeValueObject)

		if !binSizeValueIsNumber {
			logger.Debug("math_NumberBin - base parameter must be a numeric value")
			return nil
		}

		if binSizeValue != 0 {
			binSize = binSizeValue
		} else {
			logger.Debug("math_NumberBin - base must not be equal to 0")
			return nil
		}
	}

	num, numIsNumber := numToFloat64(ex)
	if !numIsNumber {
		logger.Debug("math_NumberBin - parameter must be a numeric value")
		return nil
	}

	return math.Floor(num/binSize) * binSize
}

func (c memoryExecutorContext) math_Pi() interface{} {
	return math.Pi
}

func (c memoryExecutorContext) math_Rand() interface{} {
	return rand.Float64()
}

func numToInt(ex interface{}) (int, bool) {
	switch val := ex.(type) {
	case float64:
		return int(val), true
	case int:
		return val, true
	default:
		return 0, false
	}
}

func numToFloat64(num interface{}) (float64, bool) {
	switch val := num.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	default:
		return 0, false
	}
}
