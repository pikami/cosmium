package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_AggregateFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "number": 123, "key": "a"},
		map[string]interface{}{"id": "456", "number": 456, "key": "a"},
		map[string]interface{}{"id": "789", "number": 789, "key": "b"},
		map[string]interface{}{"id": "no-number", "key": "b"},
	}

	t.Run("Should execute function AVG()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
					{
						Alias: "avg",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateAvg,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"key": "a", "avg": 289.5},
				map[string]interface{}{"key": "b", "avg": 789.0},
			},
		)
	})

	t.Run("Should execute function AVG() without GROUP BY clause", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "avg",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateAvg,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"avg": 456.0},
			},
		)
	})

	t.Run("Should execute function COUNT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
					{
						Alias: "cnt",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateCount,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"key": "a", "cnt": 2},
				map[string]interface{}{"key": "b", "cnt": 1},
			},
		)
	})

	t.Run("Should execute function MAX()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
					{
						Alias: "max",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateMax,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"key": "a", "max": 456.0},
				map[string]interface{}{"key": "b", "max": 789.0},
			},
		)
	})

	t.Run("Should execute function MIN()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
					{
						Alias: "min",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateMin,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"key": "a", "min": 123.0},
				map[string]interface{}{"key": "b", "min": 789.0},
			},
		)
	})

	t.Run("Should execute function SUM()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
					{
						Alias: "sum",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateSum,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "key"}},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"key": "a", "sum": 579.0},
				map[string]interface{}{"key": "b", "sum": 789.0},
			},
		)
	})
}
