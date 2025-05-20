package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Execute_Expressions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "age": 10, "isCool": true},
		map[string]interface{}{"id": "456", "age": 20, "isCool": false},
		map[string]interface{}{"id": "789", "age": 30, "isCool": true},
	}

	t.Run("Should execute comparison expressions in SELECT", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "isAdult",
						Type:  parsers.SelectItemTypeExpression,
						Value: parsers.ComparisonExpression{
							Operation: ">=",
							Left:      testutils.SelectItem_Path("c", "age"),
							Right:     testutils.SelectItem_Constant_Int(18),
						},
					},
					{
						Alias: "isNotCool",
						Type:  parsers.SelectItemTypeExpression,
						Value: parsers.ComparisonExpression{
							Operation: "!=",
							Left:      testutils.SelectItem_Path("c", "isCool"),
							Right:     testutils.SelectItem_Constant_Bool(true),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "isAdult": false, "isNotCool": false},
				map[string]interface{}{"id": "456", "isAdult": true, "isNotCool": true},
				map[string]interface{}{"id": "789", "isAdult": true, "isNotCool": false},
			},
		)
	})

	t.Run("Should execute logical expressions in SELECT", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "isCoolAndAdult",
						Type:  parsers.SelectItemTypeExpression,
						Value: parsers.LogicalExpression{
							Operation: parsers.LogicalExpressionTypeAnd,
							Expressions: []interface{}{
								testutils.SelectItem_Path("c", "isCool"),
								parsers.ComparisonExpression{
									Operation: ">=",
									Left:      testutils.SelectItem_Path("c", "age"),
									Right:     testutils.SelectItem_Constant_Int(18),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "isCoolAndAdult": false},
				map[string]interface{}{"id": "456", "isCoolAndAdult": false},
				map[string]interface{}{"id": "789", "isCoolAndAdult": true},
			},
		)
	})
}
