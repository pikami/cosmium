package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Execute_MathFunctions(t *testing.T) {
	t.Run("Should parse function ABS(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ABS(c.value) FROM c`,
			parsers.FunctionCallMathAbs,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function ACOS(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ACOS(c.value) FROM c`,
			parsers.FunctionCallMathAcos,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function ASIN(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ASIN(c.value) FROM c`,
			parsers.FunctionCallMathAsin,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function ATAN(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ATAN(c.value) FROM c`,
			parsers.FunctionCallMathAtan,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function CEILING(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT CEILING(c.value) FROM c`,
			parsers.FunctionCallMathCeiling,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function COS(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT COS(c.value) FROM c`,
			parsers.FunctionCallMathCos,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function COT(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT COT(c.value) FROM c`,
			parsers.FunctionCallMathCot,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function DEGREES(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT DEGREES(c.value) FROM c`,
			parsers.FunctionCallMathDegrees,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function EXP(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT EXP(c.value) FROM c`,
			parsers.FunctionCallMathExp,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function FLOOR(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT FLOOR(c.value) FROM c`,
			parsers.FunctionCallMathFloor,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitNot(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitNot(c.value) FROM c`,
			parsers.FunctionCallMathIntBitNot,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function LOG10(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT LOG10(c.value) FROM c`,
			parsers.FunctionCallMathLog10,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function RADIANS(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT RADIANS(c.value) FROM c`,
			parsers.FunctionCallMathRadians,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function ROUND(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ROUND(c.value) FROM c`,
			parsers.FunctionCallMathRound,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function SIGN(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT SIGN(c.value) FROM c`,
			parsers.FunctionCallMathSign,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function SIN(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT SIN(c.value) FROM c`,
			parsers.FunctionCallMathSin,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function SQRT(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT SQRT(c.value) FROM c`,
			parsers.FunctionCallMathSqrt,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function SQUARE(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT SQUARE(c.value) FROM c`,
			parsers.FunctionCallMathSquare,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function TAN(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT TAN(c.value) FROM c`,
			parsers.FunctionCallMathTan,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function TRUNC(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT TRUNC(c.value) FROM c`,
			parsers.FunctionCallMathTrunc,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function ATN2(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT ATN2(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathAtn2,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntAdd(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntAdd(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntAdd,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitAnd(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitAnd(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntBitAnd,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitLeftShift(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitLeftShift(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntBitLeftShift,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitOr(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitOr(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntBitOr,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitRightShift(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitRightShift(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntBitRightShift,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntBitXor(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntBitXor(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntBitXor,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntDiv(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntDiv(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntDiv,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntMod(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntMod(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntMod,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntMul(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntMul(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntMul,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function IntSub(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT IntSub(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathIntSub,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function POWER(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT POWER(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathPower,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function LOG(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT LOG(c.value) FROM c`,
			parsers.FunctionCallMathLog,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function LOG(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT LOG(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathLog,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function NumberBin(ex)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT NumberBin(c.value) FROM c`,
			parsers.FunctionCallMathNumberBin,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function NumberBin(ex1, ex2)", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT NumberBin(c.value, c.secondValue) FROM c`,
			parsers.FunctionCallMathNumberBin,
			[]interface{}{
				parsers.SelectItem{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				parsers.SelectItem{
					Path: []string{"c", "secondValue"},
					Type: parsers.SelectItemTypeField,
				},
			},
			"c",
		)
	})

	t.Run("Should parse function PI()", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT PI() FROM c`,
			parsers.FunctionCallMathPi,
			[]interface{}{},
			"c",
		)
	})

	t.Run("Should parse function RAND()", func(t *testing.T) {
		testMathFunctionParse(
			t,
			`SELECT RAND() FROM c`,
			parsers.FunctionCallMathRand,
			[]interface{}{},
			"c",
		)
	})
}

func testMathFunctionParse(
	t *testing.T,
	query string,
	expectedFunctionType parsers.FunctionCallType,
	expectedArguments []interface{},
	expectedTable string,
) {
	testQueryParse(
		t,
		query,
		parsers.SelectStmt{
			SelectItems: []parsers.SelectItem{
				{
					Type: parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type:      expectedFunctionType,
						Arguments: expectedArguments,
					},
				},
			},
			Table: parsers.Table{Value: expectedTable},
		},
	)
}
