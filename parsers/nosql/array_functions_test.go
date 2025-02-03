package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	testutils "github.com/pikami/cosmium/test_utils"
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function ARRAY_CONTAINS()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_CONTAINS(c.a1, "value") FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "a1"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_String("value"),
								nil,
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function ARRAY_CONTAINS() with partial match", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_CONTAINS(["a", "b"], "value", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								parsers.SelectItem{
									Type: parsers.SelectItemTypeArray,
									SelectItems: []parsers.SelectItem{
										testutils.SelectItem_Constant_String("a"),
										testutils.SelectItem_Constant_String("b"),
									},
								},
								testutils.SelectItem_Constant_String("value"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function ARRAY_CONTAINS_ANY()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_CONTAINS_ANY(["a", "b"], "value", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								parsers.SelectItem{
									Type: parsers.SelectItemTypeArray,
									SelectItems: []parsers.SelectItem{
										testutils.SelectItem_Constant_String("a"),
										testutils.SelectItem_Constant_String("b"),
									},
								},
								testutils.SelectItem_Constant_String("value"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse function ARRAY_CONTAINS_ALL()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ARRAY_CONTAINS_ALL(["a", "b"], "value", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								parsers.SelectItem{
									Type: parsers.SelectItemTypeArray,
									SelectItems: []parsers.SelectItem{
										testutils.SelectItem_Constant_String("a"),
										testutils.SelectItem_Constant_String("b"),
									},
								},
								testutils.SelectItem_Constant_String("value"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_Int(0),
								testutils.SelectItem_Constant_Int(2),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})
}
