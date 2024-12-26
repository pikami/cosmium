package memoryexecutor_test

import (
	"reflect"
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func testQueryExecute(
	t *testing.T,
	query parsers.SelectStmt,
	data []memoryexecutor.RowType,
	expectedData []memoryexecutor.RowType,
) {
	result := memoryexecutor.ExecuteQuery(query, data)

	if !reflect.DeepEqual(result, expectedData) {
		t.Errorf("execution result does not match expected data.\nExpected: %+v\nGot: %+v", expectedData, result)
	}
}

func Test_Execute(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "12345", "pk": 123, "_self": "self1", "_rid": "rid1", "_ts": 123456, "isCool": false},
		map[string]interface{}{"id": "67890", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
		map[string]interface{}{
			"id": "456", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true,
			"tags": []map[string]interface{}{
				{"name": "tag-a"},
				{"name": "tag-b"},
			},
		},
		map[string]interface{}{
			"id": "123", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true,
			"tags": []map[string]interface{}{
				{"name": "tag-b"},
				{"name": "tag-c"},
			},
		},
	}

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

	t.Run("Should execute SELECT with GROUP BY", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "pk"}},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"pk": 123},
				map[string]interface{}{"pk": 456},
			},
		)
	})

	t.Run("Should execute IN function", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.SelectItem{
					Type: parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type: parsers.FunctionCallIn,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "id"},
								Type: parsers.SelectItemTypeField,
							},
							testutils.SelectItem_Constant_String("123"),
							testutils.SelectItem_Constant_String("456"),
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

	t.Run("Should execute IN selector", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "name"},
						Type: parsers.SelectItemTypeField,
					},
				},
				Table: parsers.Table{
					Value: "c",
					SelectItem: parsers.SelectItem{
						Path: []string{"c", "tags"},
					},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"name": "tag-a"},
				map[string]interface{}{"name": "tag-b"},
				map[string]interface{}{"name": "tag-b"},
				map[string]interface{}{"name": "tag-c"},
			},
		)
	})
}
