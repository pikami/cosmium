package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_StringFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "pk": "aaa", "str": "hello", "rng_type": true},
		map[string]interface{}{"id": "456", "pk": "bbb", "str": "world", "rng_type": 159},
		map[string]interface{}{"id": "789", "pk": "AAA", "str": "cool world"},
	}

	t.Run("Should execute function STRINGEQUALS(ex1, ex2, ignoreCase)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "stringEquals",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallStringEquals,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "pk"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "aaa",
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "stringEquals": true},
				map[string]interface{}{"id": "456", "stringEquals": false},
				map[string]interface{}{"id": "789", "stringEquals": true},
			},
		)
	})

	t.Run("Should execute function STRINGEQUALS(ex1, ex2)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "stringEquals",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallStringEquals,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "pk"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "aaa",
									},
								},
								nil,
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "stringEquals": true},
				map[string]interface{}{"id": "456", "stringEquals": false},
				map[string]interface{}{"id": "789", "stringEquals": false},
			},
		)
	})

	t.Run("Should execute function CONCAT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "concat",
						Type:  parsers.SelectItemTypeFunctionCall,
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
										Value: " ",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: 123,
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"concat": "123 123aaa"},
				map[string]interface{}{"concat": "456 123bbb"},
				map[string]interface{}{"concat": "789 123AAA"},
			},
		)
	})

	t.Run("Should execute function CONTAINS()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "contains",
						Type:  parsers.SelectItemTypeFunctionCall,
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
										Value: "2",
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "contains": true},
				map[string]interface{}{"id": "456", "contains": false},
				map[string]interface{}{"id": "789", "contains": false},
			},
		)
	})

	t.Run("Should execute function ENDSWITH()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "ends",
						Type:  parsers.SelectItemTypeFunctionCall,
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
										Value: "3",
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "ends": true},
				map[string]interface{}{"id": "456", "ends": false},
				map[string]interface{}{"id": "789", "ends": false},
			},
		)
	})

	t.Run("Should execute function STARTSWITH()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "starts",
						Type:  parsers.SelectItemTypeFunctionCall,
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
										Value: "1",
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "starts": true},
				map[string]interface{}{"id": "456", "starts": false},
				map[string]interface{}{"id": "789", "starts": false},
			},
		)
	})

	t.Run("Should execute function INDEX_OF()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "index",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIndexOf,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeString,
										Value: "o",
									},
								},
								parsers.SelectItem{
									Type: parsers.SelectItemTypeConstant,
									Value: parsers.Constant{
										Type:  parsers.ConstantTypeInteger,
										Value: 4,
									},
								},
							},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "index": 4},
				map[string]interface{}{"str": "world", "index": -1},
				map[string]interface{}{"str": "cool world", "index": 6},
			},
		)
	})

	t.Run("Should execute function ToString()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "str",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallToString,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "rng_type"},
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
				map[string]interface{}{"id": "123", "str": "true"},
				map[string]interface{}{"id": "456", "str": "159"},
				map[string]interface{}{"id": "789", "str": ""},
			},
		)
	})
}
