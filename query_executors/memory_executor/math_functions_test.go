package memoryexecutor_test

import (
	"math"
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Execute_MathFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": 1, "value": 0.0},
		map[string]interface{}{"id": 2, "value": 1.0},
		map[string]interface{}{"id": 3, "value": -1.0},
		map[string]interface{}{"id": 4, "value": 0.5},
		map[string]interface{}{"id": 5, "value": -0.5},
		map[string]interface{}{"id": 6, "value": 0.707},
		map[string]interface{}{"id": 7, "value": -0.707},
		map[string]interface{}{"id": 8, "value": 0.866},
		map[string]interface{}{"id": 9, "value": -0.866},
	}

	mockDataInts := []memoryexecutor.RowType{
		map[string]interface{}{"id": 1, "value": -1},
		map[string]interface{}{"id": 2, "value": 0},
		map[string]interface{}{"id": 3, "value": 1},
		map[string]interface{}{"id": 4, "value": 5},
	}

	t.Run("Should execute function ABS(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathAbs,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": 0.0},
				map[string]interface{}{"value": 1.0, "result": 1.0},
				map[string]interface{}{"value": -1.0, "result": 1.0},
				map[string]interface{}{"value": 0.5, "result": 0.5},
				map[string]interface{}{"value": -0.5, "result": 0.5},
				map[string]interface{}{"value": 0.707, "result": 0.707},
				map[string]interface{}{"value": -0.707, "result": 0.707},
				map[string]interface{}{"value": 0.866, "result": 0.866},
				map[string]interface{}{"value": -0.866, "result": 0.866},
			},
		)
	})

	t.Run("Should execute function ACOS(cosine)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathAcos,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Acos(0.0) * 180 / math.Pi},
				map[string]interface{}{"value": 1.0, "result": math.Acos(1.0) * 180 / math.Pi},
				map[string]interface{}{"value": -1.0, "result": math.Acos(-1.0) * 180 / math.Pi},
				map[string]interface{}{"value": 0.5, "result": math.Acos(0.5) * 180 / math.Pi},
				map[string]interface{}{"value": -0.5, "result": math.Acos(-0.5) * 180 / math.Pi},
				map[string]interface{}{"value": 0.707, "result": math.Acos(0.707) * 180 / math.Pi},
				map[string]interface{}{"value": -0.707, "result": math.Acos(-0.707) * 180 / math.Pi},
				map[string]interface{}{"value": 0.866, "result": math.Acos(0.866) * 180 / math.Pi},
				map[string]interface{}{"value": -0.866, "result": math.Acos(-0.866) * 180 / math.Pi},
			},
		)
	})

	t.Run("Should execute function ASIN(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathAsin,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Asin(0.0) * 180 / math.Pi},
				map[string]interface{}{"value": 1.0, "result": math.Asin(1.0) * 180 / math.Pi},
				map[string]interface{}{"value": -1.0, "result": math.Asin(-1.0) * 180 / math.Pi},
				map[string]interface{}{"value": 0.5, "result": math.Asin(0.5) * 180 / math.Pi},
				map[string]interface{}{"value": -0.5, "result": math.Asin(-0.5) * 180 / math.Pi},
				map[string]interface{}{"value": 0.707, "result": math.Asin(0.707) * 180 / math.Pi},
				map[string]interface{}{"value": -0.707, "result": math.Asin(-0.707) * 180 / math.Pi},
				map[string]interface{}{"value": 0.866, "result": math.Asin(0.866) * 180 / math.Pi},
				map[string]interface{}{"value": -0.866, "result": math.Asin(-0.866) * 180 / math.Pi},
			},
		)
	})

	t.Run("Should execute function ATAN(tangent)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathAtan,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Atan(0.0) * 180 / math.Pi},
				map[string]interface{}{"value": 1.0, "result": math.Atan(1.0) * 180 / math.Pi},
				map[string]interface{}{"value": -1.0, "result": math.Atan(-1.0) * 180 / math.Pi},
				map[string]interface{}{"value": 0.5, "result": math.Atan(0.5) * 180 / math.Pi},
				map[string]interface{}{"value": -0.5, "result": math.Atan(-0.5) * 180 / math.Pi},
				map[string]interface{}{"value": 0.707, "result": math.Atan(0.707) * 180 / math.Pi},
				map[string]interface{}{"value": -0.707, "result": math.Atan(-0.707) * 180 / math.Pi},
				map[string]interface{}{"value": 0.866, "result": math.Atan(0.866) * 180 / math.Pi},
				map[string]interface{}{"value": -0.866, "result": math.Atan(-0.866) * 180 / math.Pi},
			},
		)
	})

	t.Run("Should execute function COS(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathCos,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Cos(0.0)},
				map[string]interface{}{"value": 1.0, "result": math.Cos(1.0)},
				map[string]interface{}{"value": -1.0, "result": math.Cos(-1.0)},
				map[string]interface{}{"value": 0.5, "result": math.Cos(0.5)},
				map[string]interface{}{"value": -0.5, "result": math.Cos(-0.5)},
				map[string]interface{}{"value": 0.707, "result": math.Cos(0.707)},
				map[string]interface{}{"value": -0.707, "result": math.Cos(-0.707)},
				map[string]interface{}{"value": 0.866, "result": math.Cos(0.866)},
				map[string]interface{}{"value": -0.866, "result": math.Cos(-0.866)},
			},
		)
	})

	t.Run("Should execute function COT(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathCot,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": nil},
				map[string]interface{}{"value": 1.0, "result": 1 / math.Tan(1.0)},
				map[string]interface{}{"value": -1.0, "result": 1 / math.Tan(-1.0)},
				map[string]interface{}{"value": 0.5, "result": 1 / math.Tan(0.5)},
				map[string]interface{}{"value": -0.5, "result": 1 / math.Tan(-0.5)},
				map[string]interface{}{"value": 0.707, "result": 1 / math.Tan(0.707)},
				map[string]interface{}{"value": -0.707, "result": 1 / math.Tan(-0.707)},
				map[string]interface{}{"value": 0.866, "result": 1 / math.Tan(0.866)},
				map[string]interface{}{"value": -0.866, "result": 1 / math.Tan(-0.866)},
			},
		)
	})

	t.Run("Should execute function Degrees(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathDegrees,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": 0.0 * (180 / math.Pi)},
				map[string]interface{}{"value": 1.0, "result": 1.0 * (180 / math.Pi)},
				map[string]interface{}{"value": -1.0, "result": -1.0 * (180 / math.Pi)},
				map[string]interface{}{"value": 0.5, "result": 0.5 * (180 / math.Pi)},
				map[string]interface{}{"value": -0.5, "result": -0.5 * (180 / math.Pi)},
				map[string]interface{}{"value": 0.707, "result": 0.707 * (180 / math.Pi)},
				map[string]interface{}{"value": -0.707, "result": -0.707 * (180 / math.Pi)},
				map[string]interface{}{"value": 0.866, "result": 0.866 * (180 / math.Pi)},
				map[string]interface{}{"value": -0.866, "result": -0.866 * (180 / math.Pi)},
			},
		)
	})

	t.Run("Should execute function EXP(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathExp,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Exp(0.0)},
				map[string]interface{}{"value": 1.0, "result": math.Exp(1.0)},
				map[string]interface{}{"value": -1.0, "result": math.Exp(-1.0)},
				map[string]interface{}{"value": 0.5, "result": math.Exp(0.5)},
				map[string]interface{}{"value": -0.5, "result": math.Exp(-0.5)},
				map[string]interface{}{"value": 0.707, "result": math.Exp(0.707)},
				map[string]interface{}{"value": -0.707, "result": math.Exp(-0.707)},
				map[string]interface{}{"value": 0.866, "result": math.Exp(0.866)},
				map[string]interface{}{"value": -0.866, "result": math.Exp(-0.866)},
			},
		)
	})

	t.Run("Should execute function FLOOR(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathFloor,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": math.Floor(0.0)},
				map[string]interface{}{"value": 1.0, "result": math.Floor(1.0)},
				map[string]interface{}{"value": -1.0, "result": math.Floor(-1.0)},
				map[string]interface{}{"value": 0.5, "result": math.Floor(0.5)},
				map[string]interface{}{"value": -0.5, "result": math.Floor(-0.5)},
				map[string]interface{}{"value": 0.707, "result": math.Floor(0.707)},
				map[string]interface{}{"value": -0.707, "result": math.Floor(-0.707)},
				map[string]interface{}{"value": 0.866, "result": math.Floor(0.866)},
				map[string]interface{}{"value": -0.866, "result": math.Floor(-0.866)},
			},
		)
	})

	t.Run("Should execute function IntBitNot(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathIntBitNot,
			mockDataInts,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": -1, "result": ^-1},
				map[string]interface{}{"value": 0, "result": ^0},
				map[string]interface{}{"value": 1, "result": ^1},
				map[string]interface{}{"value": 5, "result": ^5},
			},
		)
	})

	t.Run("Should execute function LOG10(value)", func(t *testing.T) {
		testMathFunctionExecute(
			t,
			parsers.FunctionCallMathLog10,
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"value": 0.0, "result": nil},
				map[string]interface{}{"value": 1.0, "result": math.Log10(1.0)},
				map[string]interface{}{"value": -1.0, "result": nil},
				map[string]interface{}{"value": 0.5, "result": math.Log10(0.5)},
				map[string]interface{}{"value": -0.5, "result": nil},
				map[string]interface{}{"value": 0.707, "result": math.Log10(0.707)},
				map[string]interface{}{"value": -0.707, "result": nil},
				map[string]interface{}{"value": 0.866, "result": math.Log10(0.866)},
				map[string]interface{}{"value": -0.866, "result": nil},
			},
		)
	})
}

func testMathFunctionExecute(
	t *testing.T,
	functionCallType parsers.FunctionCallType,
	data []memoryexecutor.RowType,
	expectedData []memoryexecutor.RowType,
) {
	testQueryExecute(
		t,
		parsers.SelectStmt{
			SelectItems: []parsers.SelectItem{
				{
					Path: []string{"c", "value"},
					Type: parsers.SelectItemTypeField,
				},
				{
					Alias: "result",
					Type:  parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type: functionCallType,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "value"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
			Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
		},
		data,
		expectedData,
	)
}
