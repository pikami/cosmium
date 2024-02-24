package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_ArrayFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "arr1": []int{1, 2, 3}, "arr2": []int{3, 4, 5}},
		map[string]interface{}{"id": "456", "arr1": []int{4, 5, 6}, "arr2": []int{5, 6, 7, 8}},
		map[string]interface{}{"id": "789", "arr1": []int{7, 8, 9}, "arr2": []int{7, 8, 9, 10, 11}},
	}

	t.Run("Should execute function CONCAT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "Concat",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayConcat,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "arr2"},
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
				map[string]interface{}{"id": "123", "Concat": []interface{}{1, 2, 3, 3, 4, 5}},
				map[string]interface{}{"id": "456", "Concat": []interface{}{4, 5, 6, 5, 6, 7, 8}},
				map[string]interface{}{"id": "789", "Concat": []interface{}{7, 8, 9, 7, 8, 9, 10, 11}},
			},
		)
	})

	t.Run("Should execute function ARRAY_LENGTH()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "Length",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayLength,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr2"},
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
				map[string]interface{}{"id": "123", "Length": 3},
				map[string]interface{}{"id": "456", "Length": 4},
				map[string]interface{}{"id": "789", "Length": 5},
			},
		)
	})

	t.Run("Should execute function ARRAY_SLICE()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "Slice",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArraySlice,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr2"},
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
										Value: 2,
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
				map[string]interface{}{"id": "123", "Slice": []interface{}{4, 5}},
				map[string]interface{}{"id": "456", "Slice": []interface{}{6, 7}},
				map[string]interface{}{"id": "789", "Slice": []interface{}{8, 9}},
			},
		)
	})

	t.Run("Should execute function SET_INTERSECT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "Intersection",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSetIntersect,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "arr2"},
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
				map[string]interface{}{"id": "123", "Intersection": []interface{}{3}},
				map[string]interface{}{"id": "456", "Intersection": []interface{}{5, 6}},
				map[string]interface{}{"id": "789", "Intersection": []interface{}{7, 8, 9}},
			},
		)
	})

	t.Run("Should execute function SET_UNION()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "Union",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSetUnion,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr1"},
									Type: parsers.SelectItemTypeField,
								},
								parsers.SelectItem{
									Path: []string{"c", "arr2"},
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
				map[string]interface{}{"id": "123", "Union": []interface{}{1, 2, 3, 4, 5}},
				map[string]interface{}{"id": "456", "Union": []interface{}{4, 5, 6, 7, 8}},
				map[string]interface{}{"id": "789", "Union": []interface{}{7, 8, 9, 10, 11}},
			},
		)
	})
}
