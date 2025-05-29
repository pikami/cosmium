package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Execute_Arithmetics(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": 1, "a": 420},
		map[string]interface{}{"id": 2, "a": 6.9},
		map[string]interface{}{"id": 3},
	}

	t.Run("Should execute simple arithmetics", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Type:  parsers.SelectItemTypeBinaryExpression,
						Alias: "result",
						Value: parsers.BinaryExpression{
							Operation: "+",
							Left:      testutils.SelectItem_Path("c", "a"),
							Right: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left:      testutils.SelectItem_Constant_Float(2.0),
									Right:     testutils.SelectItem_Constant_Int(3),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": 1, "result": 426.0},
				map[string]interface{}{"id": 2, "result": 12.9},
				map[string]interface{}{"id": 3, "result": nil},
			},
		)
	})

	t.Run("Should execute arithmetics in WHERE clause", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					testutils.SelectItem_Path("c", "id"),
					{
						Alias: "result",
						Type:  parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left:      testutils.SelectItem_Path("c", "a"),
							Right:     testutils.SelectItem_Constant_Int(2),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Filters: parsers.ComparisonExpression{
					Operation: ">",
					Left: parsers.SelectItem{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left:      testutils.SelectItem_Path("c", "a"),
							Right:     testutils.SelectItem_Constant_Int(2),
						},
					},
					Right: testutils.SelectItem_Constant_Int(500),
				},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": 1, "result": 840.0},
			},
		)
	})
}
