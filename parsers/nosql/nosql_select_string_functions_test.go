package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Execute_StringFunctions(t *testing.T) {

	t.Run("Should execute function STRINGEQUALS(ex1, ex2, ignoreCase)", func(t *testing.T) {
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

	t.Run("Should execute function STRINGEQUALS(ex1, ex2)", func(t *testing.T) {
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
}
