package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Execute_StringFunctions(t *testing.T) {

	t.Run("Should parse function STRINGEQUALS(ex1, ex2, ignoreCase)", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT STRINGEQUALS(c.id, "123", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallStringEquals,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeBoolean,
										Value: true,
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

	t.Run("Should parse function STRINGEQUALS(ex1, ex2)", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT STRINGEQUALS(c.id, "123") FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallStringEquals,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								nil,
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse function CONCAT()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT CONCAT(c.id, "123", c.pk) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallConcat,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								parsers.SelectItem{
									Path: []string{"c", "pk"},
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

	t.Run("Should parse function CONTAINS()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT CONTAINS(c.id, "123", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallContains,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeBoolean,
										Value: true,
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

	t.Run("Should parse function ENDSWITH()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ENDSWITH(c.id, "123", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallEndsWith,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeBoolean,
										Value: true,
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

	t.Run("Should parse function STARTSWITH()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT STARTSWITH(c.id, "123", true) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallStartsWith,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "123",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeBoolean,
										Value: true,
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

	t.Run("Should parse function INDEX_OF()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT INDEX_OF(c.id, "2", 1) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIndexOf,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "2",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 1,
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
}
