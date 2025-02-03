package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Parse_AggregateFunctions(t *testing.T) {

	t.Run("Should parse function AVG()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT AVG(c.a1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateAvg,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function COUNT()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT COUNT(c.a1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateCount,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function MAX()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT MAX(c.a1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateMax,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function MIN()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT MIN(c.a1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateMin,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function SUM()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT SUM(c.a1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallAggregateSum,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})
}
