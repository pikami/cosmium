package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_Select(t *testing.T) {
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

	t.Run("Should execute SELECT DISTINCT", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "pk"}},
				},
				Table:    parsers.Table{Value: "c"},
				Distinct: true,
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"pk": 123},
				map[string]interface{}{"pk": 456},
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

	t.Run("Should execute SELECT OFFSET", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table:  parsers.Table{Value: "c"},
				Count:  2,
				Offset: 1,
				OrderExpressions: []parsers.OrderExpression{
					{
						SelectItem: parsers.SelectItem{Path: []string{"c", "id"}},
						Direction:  parsers.OrderDirectionDesc,
					},
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "67890", "pk": 456},
				map[string]interface{}{"id": "456", "pk": 456},
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
}
