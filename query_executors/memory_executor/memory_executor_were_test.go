package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_Where(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "12345", "pk": 123, "_self": "self1", "_rid": "rid1", "_ts": 123456, "isCool": false},
		map[string]interface{}{"id": "67890", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
		map[string]interface{}{"id": "456", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
		map[string]interface{}{"id": "123", "pk": 456, "_self": "self2", "_rid": "rid2", "_ts": 789012, "isCool": true},
	}

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
					Right: parsers.SelectItem{
						Type:  parsers.SelectItemTypeConstant,
						Value: parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
					},
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
					Right: parsers.SelectItem{
						Type:  parsers.SelectItemTypeConstant,
						Value: parsers.Constant{Type: parsers.ConstantTypeParameterConstant, Value: "@param_id"},
					},
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
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "67890"},
							},
						},
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "pk"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 456},
							},
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
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
							},
						},
						parsers.LogicalExpression{
							Operation: parsers.LogicalExpressionTypeOr,
							Expressions: []interface{}{
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right: parsers.SelectItem{
										Type:  parsers.SelectItemTypeConstant,
										Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "123"},
									},
								},
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right: parsers.SelectItem{
										Type:  parsers.SelectItemTypeConstant,
										Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "456"},
									},
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
