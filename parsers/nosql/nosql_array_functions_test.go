package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Parse_ArrayFunctions(t *testing.T) {

	t.Run("Should parse function ARRAY_CONCAT()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_CONCAT(c.a1, c.a2) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayConcat,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "a2"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse function ARRAY_LENGTH()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_LENGTH(c.array) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayLength,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "array"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse function ARRAY_SLICE()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_SLICE(c.array, 0, 2) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArraySlice,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "array"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 0,
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 2,
									},
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse function SetIntersect()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT SetIntersect(c.set1, c.set2) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSetIntersect,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "set1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "set2"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse function SetUnion()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT SetUnion(c.set1, c.set2) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSetUnion,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "set1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "set2"},
									Type: parsers.SelectItemTypeField,
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})
}
