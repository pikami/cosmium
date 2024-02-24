package memoryexecutor_test

import (
	"math"
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_TypeCheckingFunctions(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "1", "obj": nil},
		map[string]interface{}{"id": "2", "obj": "str"},
		map[string]interface{}{"id": "3", "obj": 1},
		map[string]interface{}{"id": "4", "obj": 1.2},
		map[string]interface{}{"id": "5", "obj": true},
		map[string]interface{}{"id": "6", "obj": []interface{}{1, 2, 3}},
		map[string]interface{}{"id": "7", "obj": map[string]interface{}{"a": "a"}},
		map[string]interface{}{"id": "8", "obj": math.Inf(1)},
	}

	t.Run("Should execute function IS_DEFINED(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsDefined",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsDefined,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsDefined": false},
				map[string]interface{}{"id": "2", "IsDefined": true},
				map[string]interface{}{"id": "3", "IsDefined": true},
				map[string]interface{}{"id": "4", "IsDefined": true},
				map[string]interface{}{"id": "5", "IsDefined": true},
				map[string]interface{}{"id": "6", "IsDefined": true},
				map[string]interface{}{"id": "7", "IsDefined": true},
				map[string]interface{}{"id": "8", "IsDefined": true},
			},
		)
	})

	t.Run("Should execute function IS_ARRAY(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsArray",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsArray,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsArray": false},
				map[string]interface{}{"id": "2", "IsArray": false},
				map[string]interface{}{"id": "3", "IsArray": false},
				map[string]interface{}{"id": "4", "IsArray": false},
				map[string]interface{}{"id": "5", "IsArray": false},
				map[string]interface{}{"id": "6", "IsArray": true},
				map[string]interface{}{"id": "7", "IsArray": false},
				map[string]interface{}{"id": "8", "IsArray": false},
			},
		)
	})

	t.Run("Should execute function IS_BOOL(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsBool",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsBool,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsBool": false},
				map[string]interface{}{"id": "2", "IsBool": false},
				map[string]interface{}{"id": "3", "IsBool": false},
				map[string]interface{}{"id": "4", "IsBool": false},
				map[string]interface{}{"id": "5", "IsBool": true},
				map[string]interface{}{"id": "6", "IsBool": false},
				map[string]interface{}{"id": "7", "IsBool": false},
				map[string]interface{}{"id": "8", "IsBool": false},
			},
		)
	})

	t.Run("Should execute function IS_FINITENUMBER(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsFiniteNumber",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsFiniteNumber,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsFiniteNumber": false},
				map[string]interface{}{"id": "2", "IsFiniteNumber": false},
				map[string]interface{}{"id": "3", "IsFiniteNumber": true},
				map[string]interface{}{"id": "4", "IsFiniteNumber": true},
				map[string]interface{}{"id": "5", "IsFiniteNumber": false},
				map[string]interface{}{"id": "6", "IsFiniteNumber": false},
				map[string]interface{}{"id": "7", "IsFiniteNumber": false},
				map[string]interface{}{"id": "8", "IsFiniteNumber": false},
			},
		)
	})

	t.Run("Should execute function IS_INTEGER(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsInteger",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsInteger,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsInteger": false},
				map[string]interface{}{"id": "2", "IsInteger": false},
				map[string]interface{}{"id": "3", "IsInteger": true},
				map[string]interface{}{"id": "4", "IsInteger": false},
				map[string]interface{}{"id": "5", "IsInteger": false},
				map[string]interface{}{"id": "6", "IsInteger": false},
				map[string]interface{}{"id": "7", "IsInteger": false},
				map[string]interface{}{"id": "8", "IsInteger": false},
			},
		)
	})

	t.Run("Should execute function IS_NULL(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsNull",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsNull,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsNull": true},
				map[string]interface{}{"id": "2", "IsNull": false},
				map[string]interface{}{"id": "3", "IsNull": false},
				map[string]interface{}{"id": "4", "IsNull": false},
				map[string]interface{}{"id": "5", "IsNull": false},
				map[string]interface{}{"id": "6", "IsNull": false},
				map[string]interface{}{"id": "7", "IsNull": false},
				map[string]interface{}{"id": "8", "IsNull": false},
			},
		)
	})

	t.Run("Should execute function IS_NUMBER(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsNumber",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsNumber,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsNumber": false},
				map[string]interface{}{"id": "2", "IsNumber": false},
				map[string]interface{}{"id": "3", "IsNumber": true},
				map[string]interface{}{"id": "4", "IsNumber": true},
				map[string]interface{}{"id": "5", "IsNumber": false},
				map[string]interface{}{"id": "6", "IsNumber": false},
				map[string]interface{}{"id": "7", "IsNumber": false},
				map[string]interface{}{"id": "8", "IsNumber": true},
			},
		)
	})

	t.Run("Should execute function IS_OBJECT(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsObject",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsObject,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsObject": false},
				map[string]interface{}{"id": "2", "IsObject": false},
				map[string]interface{}{"id": "3", "IsObject": false},
				map[string]interface{}{"id": "4", "IsObject": false},
				map[string]interface{}{"id": "5", "IsObject": false},
				map[string]interface{}{"id": "6", "IsObject": false},
				map[string]interface{}{"id": "7", "IsObject": true},
				map[string]interface{}{"id": "8", "IsObject": false},
			},
		)
	})

	t.Run("Should execute function IS_PRIMITIVE(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsPrimitive",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsPrimitive,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsPrimitive": true},
				map[string]interface{}{"id": "2", "IsPrimitive": true},
				map[string]interface{}{"id": "3", "IsPrimitive": true},
				map[string]interface{}{"id": "4", "IsPrimitive": true},
				map[string]interface{}{"id": "5", "IsPrimitive": true},
				map[string]interface{}{"id": "6", "IsPrimitive": false},
				map[string]interface{}{"id": "7", "IsPrimitive": false},
				map[string]interface{}{"id": "8", "IsPrimitive": true},
			},
		)
	})

	t.Run("Should execute function IS_STRING(path)", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{
						Alias: "IsString",
						Type:  parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsString,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "obj"},
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
				map[string]interface{}{"id": "1", "IsString": false},
				map[string]interface{}{"id": "2", "IsString": true},
				map[string]interface{}{"id": "3", "IsString": false},
				map[string]interface{}{"id": "4", "IsString": false},
				map[string]interface{}{"id": "5", "IsString": false},
				map[string]interface{}{"id": "6", "IsString": false},
				map[string]interface{}{"id": "7", "IsString": false},
				map[string]interface{}{"id": "8", "IsString": false},
			},
		)
	})
}
