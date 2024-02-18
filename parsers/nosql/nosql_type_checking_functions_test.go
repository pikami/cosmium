package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Execute_TypeCheckingFunctions(t *testing.T) {

	t.Run("Should parse function IS_DEFINED", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_DEFINED(c.id) FROM c WHERE IS_DEFINED(c.pk)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsDefined,
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
				Filters: parsers.SelectItem{
					Type: parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type: parsers.FunctionCallIsDefined,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "pk"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})
}
