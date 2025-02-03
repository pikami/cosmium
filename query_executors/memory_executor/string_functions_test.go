package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Execute_StringFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "pk": "aaa", "str": "hello", "rng_type": true, "str2": "   hello "},
		map[string]interface{}{"id": "456", "pk": "bbb", "str": "world", "rng_type": 159, "str2": "  world   "},
		map[string]interface{}{"id": "789", "pk": "AAA", "str": "cool world", "str2": "          cool world  "},
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
								testutils.SelectItem_Constant_String("aaa"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_String("aaa"),
								nil,
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_String(" "),
								testutils.SelectItem_Constant_Int(123),
								parsers.SelectItem{
									Path: []string{"c", "pk"},
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
								testutils.SelectItem_Constant_String("2"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_String("3"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_String("1"),
								testutils.SelectItem_Constant_Bool(true),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
								testutils.SelectItem_Constant_String("o"),
								testutils.SelectItem_Constant_Int(4),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "123", "str": "true"},
				map[string]interface{}{"id": "456", "str": "159"},
				map[string]interface{}{"id": "789", "str": ""},
			},
		)
	})

	t.Run("Should execute function LEFT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "left",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLeft,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_Int(3),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "left": "hel"},
				map[string]interface{}{"str": "world", "left": "wor"},
				map[string]interface{}{"str": "cool world", "left": "coo"},
			},
		)
	})

	t.Run("Should execute function LENGTH()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "length",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLength,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
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
				map[string]interface{}{"str": "hello", "length": 5},
				map[string]interface{}{"str": "world", "length": 5},
				map[string]interface{}{"str": "cool world", "length": 10},
			},
		)
	})

	t.Run("Should execute function LTRIM()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "ltrimmed",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallLTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str2"},
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
				map[string]interface{}{"str": "hello", "ltrimmed": "hello "},
				map[string]interface{}{"str": "world", "ltrimmed": "world   "},
				map[string]interface{}{"str": "cool world", "ltrimmed": "cool world  "},
			},
		)
	})

	t.Run("Should execute function REPLACE()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "replaced",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReplace,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_String("world"),
								testutils.SelectItem_Constant_String("universe"),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "replaced": "hello"},
				map[string]interface{}{"str": "world", "replaced": "universe"},
				map[string]interface{}{"str": "cool world", "replaced": "cool universe"},
			},
		)
	})

	t.Run("Should execute function REPLICATE()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "replicated",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReplicate,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_Int(3),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "replicated": "hellohellohello"},
				map[string]interface{}{"str": "world", "replicated": "worldworldworld"},
				map[string]interface{}{"str": "cool world", "replicated": "cool worldcool worldcool world"},
			},
		)
	})

	t.Run("Should execute function REVERSE()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "reversed",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallReverse,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
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
				map[string]interface{}{"str": "hello", "reversed": "olleh"},
				map[string]interface{}{"str": "world", "reversed": "dlrow"},
				map[string]interface{}{"str": "cool world", "reversed": "dlrow looc"},
			},
		)
	})

	t.Run("Should execute function RIGHT()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "right",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallRight,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_Int(3),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "right": "llo"},
				map[string]interface{}{"str": "world", "right": "rld"},
				map[string]interface{}{"str": "cool world", "right": "rld"},
			},
		)
	})

	t.Run("Should execute function RTRIM()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "rtrimmed",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallRTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str2"},
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
				map[string]interface{}{"str": "hello", "rtrimmed": "   hello"},
				map[string]interface{}{"str": "world", "rtrimmed": "  world"},
				map[string]interface{}{"str": "cool world", "rtrimmed": "          cool world"},
			},
		)
	})

	t.Run("Should execute function SUBSTRING()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "substring",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallSubstring,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str"},
									Type: parsers.SelectItemTypeField,
								},
								testutils.SelectItem_Constant_Int(2),
								testutils.SelectItem_Constant_Int(4),
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"str": "hello", "substring": "llo"},
				map[string]interface{}{"str": "world", "substring": "rld"},
				map[string]interface{}{"str": "cool world", "substring": "ol w"},
			},
		)
	})

	t.Run("Should execute function TRIM()", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "str"},
						Type: parsers.SelectItemTypeField,
					},
					{
						Alias: "trimmed",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallTrim,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "str2"},
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
				map[string]interface{}{"str": "hello", "trimmed": "hello"},
				map[string]interface{}{"str": "world", "trimmed": "world"},
				map[string]interface{}{"str": "cool world", "trimmed": "cool world"},
			},
		)
	})
}
