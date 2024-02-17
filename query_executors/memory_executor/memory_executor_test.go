package memoryexecutor_test

import (
	"reflect"
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func testQueryExecute(
	t *testing.T,
	query parsers.SelectStmt,
	data []memoryexecutor.RowType,
	expectedData []memoryexecutor.RowType,
) {
	result := memoryexecutor.Execute(query, data)

	if !reflect.DeepEqual(result, expectedData) {
		t.Errorf("execution result does not match expected data.\nExpected: %+v\nGot: %+v", expectedData, result)
	}
}

func Test_Execute(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "12345", "pk": 123, "_self": "self1", "_rid": "rid1", "_ts": 123456, "isCool": false},
		map[string]interface{}{"id": "67890", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
		map[string]interface{}{"id": "456", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
		map[string]interface{}{"id": "123", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
	}

	t.Run("Should execute simple SELECT", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "12345", "pk": 123},
				map[string]interface{}{"id": "67890", "pk": 456},
				map[string]interface{}{"id": "456", "pk": 456},
				map[string]interface{}{"id": "123", "pk": 456},
			},
		)
	})

	t.Run("Should execute SELECT TOP", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
				Count: 1,
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "12345", "pk": 123},
			},
		)
	})

	t.Run("Should execute SELECT VALUE", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}, IsTopLevel: true},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				"12345",
				"67890",
				"456",
				"123",
			},
		)
	})

	t.Run("Should execute SELECT *", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c"}, IsTopLevel: true},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			mockData,
		)
	})

	t.Run("Should execute SELECT array", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "arr",
						Type:  parsers.SelectItemTypeArray,
						SelectItems: []parsers.SelectItem{
							{Path: []string{"c", "id"}},
							{Path: []string{"c", "pk"}},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"arr": []interface{}{"12345", 123}},
				map[string]interface{}{"arr": []interface{}{"67890", 456}},
				map[string]interface{}{"arr": []interface{}{"456", 456}},
				map[string]interface{}{"arr": []interface{}{"123", 456}},
			},
		)
	})

	t.Run("Should execute SELECT object", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "obj",
						Type:  parsers.SelectItemTypeObject,
						SelectItems: []parsers.SelectItem{
							{Alias: "id", Path: []string{"c", "id"}},
							{Alias: "_pk", Path: []string{"c", "pk"}},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"obj": map[string]interface{}{"id": "12345", "_pk": 123}},
				map[string]interface{}{"obj": map[string]interface{}{"id": "67890", "_pk": 456}},
				map[string]interface{}{"obj": map[string]interface{}{"id": "456", "_pk": 456}},
				map[string]interface{}{"obj": map[string]interface{}{"id": "123", "_pk": 456}},
			},
		)
	})

	t.Run("Should execute SELECT with ORDER BY", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
				OrderExpressions: []parsers.OrderExpression{
					{
						SelectItem: parsers.SelectItem{Path: []string{"c", "pk"}},
						Direction:  parsers.OrderDirectionAsc,
					},
					{
						SelectItem: parsers.SelectItem{Path: []string{"c", "id"}},
						Direction:  parsers.OrderDirectionDesc,
					},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "12345", "pk": 123},
				map[string]interface{}{"id": "67890", "pk": 456},
				map[string]interface{}{"id": "456", "pk": 456},
				map[string]interface{}{"id": "123", "pk": 456},
			},
		)
	})

	t.Run("Should execute SELECT with single WHERE condition", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.ComparisonExpression{
					Operation: "=",
					Left:      parsers.SelectItem{Path: []string{"c", "isCool"}},
					Right:     parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "67890"},
				map[string]interface{}{"id": "456"},
				map[string]interface{}{"id": "123"},
			},
		)
	})

	t.Run("Should execute SELECT with WHERE condition with defined parameter constant", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.ComparisonExpression{
					Operation: "=",
					Left:      parsers.SelectItem{Path: []string{"c", "id"}},
					Right:     parsers.Constant{Type: parsers.ConstantTypeParameterConstant, Value: "@param_id"},
				},
				Parameters: map[string]interface{}{
					"@param_id": "456",
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "456"},
			},
		)
	})

	t.Run("Should execute SELECT with multiple WHERE conditions", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "_self"}, Alias: "self"},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Operation: parsers.LogicalExpressionTypeAnd,
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "id"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "67890"},
						},
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "pk"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 456},
						},
					},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "67890", "self": "self2"},
			},
		)
	})

	t.Run("Should execute SELECT with grouped WHERE conditions", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Operation: parsers.LogicalExpressionTypeAnd,
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "isCool"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
						},
						parsers.LogicalExpression{
							Operation: parsers.LogicalExpressionTypeOr,
							Expressions: []interface{}{
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "123"},
								},
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "456"},
								},
							},
						},
					},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "456"},
				map[string]interface{}{"id": "123"},
			},
		)
	})
}
