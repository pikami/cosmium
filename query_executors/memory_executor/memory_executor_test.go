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
}
