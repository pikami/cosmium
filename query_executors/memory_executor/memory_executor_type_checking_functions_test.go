package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_TypeCheckingFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "pk": "aaa"},
		map[string]interface{}{"id": "456"},
		map[string]interface{}{"id": "789", "pk": ""},
	}

	t.Run("Should execute function IS_DEFINED(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "defined",
						Type:  parsers.SelectItemTypeFunctionCall,
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
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "defined": true},
				map[string]interface{}{"id": "456", "defined": false},
				map[string]interface{}{"id": "789", "defined": true},
			},
		)
	})
}
