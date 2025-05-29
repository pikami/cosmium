package memoryexecutor

import (
	"math"
	"math/rand"

	"github.com/pikami/cosmium/internal/logger"
	"github.com/pikami/cosmium/parsers"
)

func (r rowContext) math_Abs(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case float64:
		return math.Abs(val)
	case int:
		if val < 0 {
			return -val
		}
		return val
	default:
		logger.DebugLn("math_Abs - got parameters of wrong type")
		return 0
	}
}

func (r rowContext) math_Acos(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Acos - got parameters of wrong type")
		return nil
	}

	if val < -1 || val > 1 {
		logger.DebugLn("math_Acos - value out of domain for acos")
		return nil
	}

	return math.Acos(val) * 180 / math.Pi
}

func (r rowContext) math_Asin(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Asin - got parameters of wrong type")
		return nil
	}

	if val < -1 || val > 1 {
		logger.DebugLn("math_Asin - value out of domain for acos")
		return nil
	}

	return math.Asin(val) * 180 / math.Pi
}

func (r rowContext) math_Atan(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Atan - got parameters of wrong type")
		return nil
	}

	return math.Atan(val) * 180 / math.Pi
}

func (r rowContext) math_Ceiling(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case float64:
		return math.Ceil(val)
	case int:
		return val
	default:
		logger.DebugLn("math_Ceiling - got parameters of wrong type")
		return 0
	}
}

func (r rowContext) math_Cos(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Cos - got parameters of wrong type")
		return nil
	}

	return math.Cos(val)
}

func (r rowContext) math_Cot(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Cot - got parameters of wrong type")
		return nil
	}

	if val == 0 {
		logger.DebugLn("math_Cot - cotangent undefined for zero")
		return nil
	}

	return 1 / math.Tan(val)
}

func (r rowContext) math_Degrees(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Degrees - got parameters of wrong type")
		return nil
	}

	return val * (180 / math.Pi)
}

func (r rowContext) math_Exp(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Exp - got parameters of wrong type")
		return nil
	}

	return math.Exp(val)
}

func (r rowContext) math_Floor(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case float64:
		return math.Floor(val)
	case int:
		return val
	default:
		logger.DebugLn("math_Floor - got parameters of wrong type")
		return 0
	}
}

func (r rowContext) math_IntBitNot(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case int:
		return ^val
	default:
		logger.DebugLn("math_IntBitNot - got parameters of wrong type")
		return nil
	}
}

func (r rowContext) math_Log10(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Log10 - got parameters of wrong type")
		return nil
	}

	if val <= 0 {
		logger.DebugLn("math_Log10 - value must be greater than 0")
		return nil
	}

	return math.Log10(val)
}

func (r rowContext) math_Radians(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Radians - got parameters of wrong type")
		return nil
	}

	return val * (math.Pi / 180.0)
}

func (r rowContext) math_Round(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case float64:
		return math.Round(val)
	case int:
		return val
	default:
		logger.DebugLn("math_Round - got parameters of wrong type")
		return nil
	}
}

func (r rowContext) math_Sign(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

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
		logger.DebugLn("math_Sign - got parameters of wrong type")
		return nil
	}
}

func (r rowContext) math_Sin(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Sin - got parameters of wrong type")
		return nil
	}

	return math.Sin(val)
}

func (r rowContext) math_Sqrt(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Sqrt - got parameters of wrong type")
		return nil
	}

	return math.Sqrt(val)
}

func (r rowContext) math_Square(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Square - got parameters of wrong type")
		return nil
	}

	return math.Pow(val, 2)
}

func (r rowContext) math_Tan(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	val, valIsNumber := numToFloat64(ex)
	if !valIsNumber {
		logger.DebugLn("math_Tan - got parameters of wrong type")
		return nil
	}

	return math.Tan(val)
}

func (r rowContext) math_Trunc(arguments []interface{}) interface{} {
	exItem := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem)

	switch val := ex.(type) {
	case float64:
		return math.Trunc(val)
	case int:
		return float64(val)
	default:
		logger.DebugLn("math_Trunc - got parameters of wrong type")
		return nil
	}
}

func (r rowContext) math_Atn2(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	y, yIsNumber := numToFloat64(ex1)
	x, xIsNumber := numToFloat64(ex2)

	if !yIsNumber || !xIsNumber {
		logger.DebugLn("math_Atn2 - got parameters of wrong type")
		return nil
	}

	return math.Atan2(y, x)
}

func (r rowContext) math_IntAdd(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	ex1Number, ex1IsNumber := numToInt(ex1)
	ex2Number, ex2IsNumber := numToInt(ex2)

	if !ex1IsNumber || !ex2IsNumber {
		logger.DebugLn("math_IntAdd - got parameters of wrong type")
		return nil
	}

	return ex1Number + ex2Number
}

func (r rowContext) math_IntBitAnd(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	ex1Int, ex1IsInt := numToInt(ex1)
	ex2Int, ex2IsInt := numToInt(ex2)

	if !ex1IsInt || !ex2IsInt {
		logger.DebugLn("math_IntBitAnd - got parameters of wrong type")
		return nil
	}

	return ex1Int & ex2Int
}

func (r rowContext) math_IntBitLeftShift(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := numToInt(ex1)
	num2, num2IsInt := numToInt(ex2)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntBitLeftShift - got parameters of wrong type")
		return nil
	}

	return num1 << uint(num2)
}

func (r rowContext) math_IntBitOr(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntBitOr - got parameters of wrong type")
		return nil
	}

	return num1 | num2
}

func (r rowContext) math_IntBitRightShift(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := numToInt(ex1)
	num2, num2IsInt := numToInt(ex2)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntBitRightShift - got parameters of wrong type")
		return nil
	}

	return num1 >> uint(num2)
}

func (r rowContext) math_IntBitXor(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntBitXor - got parameters of wrong type")
		return nil
	}

	return num1 ^ num2
}

func (r rowContext) math_IntDiv(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt || num2 == 0 {
		logger.DebugLn("math_IntDiv - got parameters of wrong type or divide by zero")
		return nil
	}

	return num1 / num2
}

func (r rowContext) math_IntMul(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntMul - got parameters of wrong type")
		return nil
	}

	return num1 * num2
}

func (r rowContext) math_IntSub(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt {
		logger.DebugLn("math_IntSub - got parameters of wrong type")
		return nil
	}

	return num1 - num2
}

func (r rowContext) math_IntMod(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	num1, num1IsInt := ex1.(int)
	num2, num2IsInt := ex2.(int)

	if !num1IsInt || !num2IsInt || num2 == 0 {
		logger.DebugLn("math_IntMod - got parameters of wrong type or divide by zero")
		return nil
	}

	return num1 % num2
}

func (r rowContext) math_Power(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	exItem2 := arguments[1].(parsers.SelectItem)
	ex1 := r.resolveSelectItem(exItem1)
	ex2 := r.resolveSelectItem(exItem2)

	base, baseIsNumber := numToFloat64(ex1)
	exponent, exponentIsNumber := numToFloat64(ex2)

	if !baseIsNumber || !exponentIsNumber {
		logger.DebugLn("math_Power - got parameters of wrong type")
		return nil
	}

	return math.Pow(base, exponent)
}

func (r rowContext) math_Log(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem1)

	var base float64 = math.E
	if len(arguments) > 1 {
		exItem2 := arguments[1].(parsers.SelectItem)
		baseValueObject := r.resolveSelectItem(exItem2)
		baseValue, baseValueIsNumber := numToFloat64(baseValueObject)

		if !baseValueIsNumber {
			logger.DebugLn("math_Log - base parameter must be a numeric value")
			return nil
		}

		if baseValue > 0 && baseValue != 1 {
			base = baseValue
		} else {
			logger.DebugLn("math_Log - base must be greater than 0 and not equal to 1")
			return nil
		}
	}

	num, numIsNumber := numToFloat64(ex)
	if !numIsNumber || num <= 0 {
		logger.DebugLn("math_Log - parameter must be a positive numeric value")
		return nil
	}

	return math.Log(num) / math.Log(base)
}

func (r rowContext) math_NumberBin(arguments []interface{}) interface{} {
	exItem1 := arguments[0].(parsers.SelectItem)
	ex := r.resolveSelectItem(exItem1)

	binSize := 1.0

	if len(arguments) > 1 {
		exItem2 := arguments[1].(parsers.SelectItem)
		binSizeValueObject := r.resolveSelectItem(exItem2)
		binSizeValue, binSizeValueIsNumber := numToFloat64(binSizeValueObject)

		if !binSizeValueIsNumber {
			logger.DebugLn("math_NumberBin - base parameter must be a numeric value")
			return nil
		}

		if binSizeValue != 0 {
			binSize = binSizeValue
		} else {
			logger.DebugLn("math_NumberBin - base must not be equal to 0")
			return nil
		}
	}

	num, numIsNumber := numToFloat64(ex)
	if !numIsNumber {
		logger.DebugLn("math_NumberBin - parameter must be a numeric value")
		return nil
	}

	return math.Floor(num/binSize) * binSize
}

func (r rowContext) math_Pi() interface{} {
	return math.Pi
}

func (r rowContext) math_Rand() interface{} {
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
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}
