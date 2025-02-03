package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "Concat": []interface{}{1, 2, 3, 3, 4, 5}},
				map[string]interface{}{"id": "456", "Concat": []interface{}{4, 5, 6, 5, 6, 7, 8}},
				map[string]interface{}{"id": "789", "Concat": []interface{}{7, 8, 9, 7, 8, 9, 10, 11}},
			},
		)
	})

	t.Run("Should execute function ARRAY_CONTAINS()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				Parameters: map[string]interface{}{
					"@categories":                []interface{}{"coats", "jackets", "sweatshirts"},
					"@objectArray":               []interface{}{map[string]interface{}{"category": "shirts", "color": "blue"}},
					"@fullMatchObject":           map[string]interface{}{"category": "shirts", "color": "blue"},
					"@partialMatchObject":        map[string]interface{}{"category": "shirts"},
					"@missingPartialMatchObject": map[string]interface{}{"category": "shorts", "color": "blue"},
				},
				SelectItems: []parsers.SelectItem{
					{
						Alias: "ContainsItem",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@categories"),
								testutils.SelectItem_Constant_String("coats"),
							},
						},
					},
					{
						Alias: "MissingItem",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@categories"),
								testutils.SelectItem_Constant_String("hoodies"),
							},
						},
					},
					{
						Alias: "ContainsFullMatchObject",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@objectArray"),
								testutils.SelectItem_Constant_Parameter("@fullMatchObject"),
							},
						},
					},
					{
						Alias: "MissingFullMatchObject",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@objectArray"),
								testutils.SelectItem_Constant_Parameter("@partialMatchObject"),
							},
						},
					},
					{
						Alias: "ContainsPartialMatchObject",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@objectArray"),
								testutils.SelectItem_Constant_Parameter("@partialMatchObject"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
					{
						Alias: "MissingPartialMatchObject",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContains,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@objectArray"),
								testutils.SelectItem_Constant_Parameter("@missingPartialMatchObject"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
			},
			[]memoryexecutor.RowType{map[string]interface{}{"id": "123"}},
			[]memoryexecutor.RowType{
				map[string]interface{}{
					"ContainsItem":               true,
					"MissingItem":                false,
					"ContainsFullMatchObject":    true,
					"MissingFullMatchObject":     false,
					"ContainsPartialMatchObject": true,
					"MissingPartialMatchObject":  false,
				},
			},
		)
	})

	t.Run("Should execute function ARRAY_CONTAINS_ANY()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				Parameters: map[string]interface{}{
					"@mixedArray": []interface{}{1, true, "3", []int{1, 2, 3}},
					"@numbers":    []interface{}{1, 2, 3, 4},
					"@emptyArray": []interface{}{},
					"@arr123":     []interface{}{1, 2, 3},
				},
				SelectItems: []parsers.SelectItem{
					{
						Alias: "matchesEntireArray",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@mixedArray"),
								testutils.SelectItem_Constant_Int(1),
								testutils.SelectItem_Constant_Bool(true),
								testutils.SelectItem_Constant_String("3"),
								testutils.SelectItem_Constant_Parameter("@arr123"),
							},
						},
					},
					{
						Alias: "matchesSomeValues",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(2),
								testutils.SelectItem_Constant_Int(3),
								testutils.SelectItem_Constant_Int(4),
								testutils.SelectItem_Constant_Int(5),
							},
						},
					},
					{
						Alias: "matchSingleValue",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(1),
							},
						},
					},
					{
						Alias: "noMatches",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(5),
								testutils.SelectItem_Constant_Int(6),
								testutils.SelectItem_Constant_Int(7),
								testutils.SelectItem_Constant_Int(8),
							},
						},
					},
					{
						Alias: "emptyArray",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAny,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@emptyArray"),
								testutils.SelectItem_Constant_Int(1),
								testutils.SelectItem_Constant_Int(2),
								testutils.SelectItem_Constant_Int(3),
							},
						},
					},
				},
			},
			[]memoryexecutor.RowType{map[string]interface{}{"id": "123"}},
			[]memoryexecutor.RowType{
				map[string]interface{}{
					"matchesEntireArray": true,
					"matchesSomeValues":  true,
					"matchSingleValue":   true,
					"noMatches":          false,
					"emptyArray":         false,
				},
			},
		)
	})

	t.Run("Should execute function ARRAY_CONTAINS_ALL()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				Parameters: map[string]interface{}{
					"@mixedArray": []interface{}{1, true, "3", []interface{}{1, 2, 3}},
					"@numbers":    []interface{}{1, 2, 3, 4},
					"@emptyArray": []interface{}{},
					"@arr123":     []interface{}{1, 2, 3},
				},
				SelectItems: []parsers.SelectItem{
					{
						Alias: "matchesEntireArray",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@mixedArray"),
								testutils.SelectItem_Constant_Int(1),
								testutils.SelectItem_Constant_Bool(true),
								testutils.SelectItem_Constant_String("3"),
								testutils.SelectItem_Constant_Parameter("@arr123"),
							},
						},
					},
					{
						Alias: "matchesSomeValues",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(2),
								testutils.SelectItem_Constant_Int(3),
								testutils.SelectItem_Constant_Int(4),
								testutils.SelectItem_Constant_Int(5),
							},
						},
					},
					{
						Alias: "matchSingleValue",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(1),
							},
						},
					},
					{
						Alias: "noMatches",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@numbers"),
								testutils.SelectItem_Constant_Int(5),
								testutils.SelectItem_Constant_Int(6),
								testutils.SelectItem_Constant_Int(7),
								testutils.SelectItem_Constant_Int(8),
							},
						},
					},
					{
						Alias: "emptyArray",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallArrayContainsAll,
							Arguments: []interface{}{
								testutils.SelectItem_Constant_Parameter("@emptyArray"),
								testutils.SelectItem_Constant_Int(1),
								testutils.SelectItem_Constant_Int(2),
								testutils.SelectItem_Constant_Int(3),
							},
						},
					},
				},
			},
			[]memoryexecutor.RowType{map[string]interface{}{"id": "123"}},
			[]memoryexecutor.RowType{
				map[string]interface{}{
					"matchesEntireArray": true,
					"matchesSomeValues":  false,
					"matchSingleValue":   true,
					"noMatches":          false,
					"emptyArray":         false,
				},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_Int(1),
								testutils.SelectItem_Constant_Int(2),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
