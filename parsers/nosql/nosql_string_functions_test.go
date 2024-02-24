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

	t.Run("Should parse function ToString()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ToString(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallToString,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function UPPER()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT UPPER(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallUpper,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function LOWER()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT LOWER(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLower,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function LEFT()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT LEFT(c.id, 5) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLeft,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 5,
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

	t.Run("Should parse function LENGTH()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT LENGTH(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLength,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function LTRIM()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT LTRIM(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function REPLACE()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT REPLACE(c.id, "old", "new") FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReplace,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "old",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "new",
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

	t.Run("Should parse function REPLICATE()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT REPLICATE(c.id, 3) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReplicate,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 3,
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

	t.Run("Should parse function REVERSE()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT REVERSE(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReverse,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function RIGHT()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT RIGHT(c.id, 3) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallRight,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 3,
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

	t.Run("Should parse function RTRIM()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT RTRIM(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallRTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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

	t.Run("Should parse function SUBSTRING()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT SUBSTRING(c.id, 1, 5) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSubstring,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 1,
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 5,
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

	t.Run("Should parse function TRIM()", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT TRIM(c.id) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "id"},
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
