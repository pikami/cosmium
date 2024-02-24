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

	t.Run("Should parse function IS_ARRAY", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_ARRAY(c.arr) FROM c WHERE IS_ARRAY(c.arr)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsArray,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "arr"},
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
						Type: parsers.FunctionCallIsArray,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "arr"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_BOOL", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_BOOL(c.flag) FROM c WHERE IS_BOOL(c.flag)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsBool,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "flag"},
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
						Type: parsers.FunctionCallIsBool,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "flag"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_FINITE_NUMBER", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_FINITE_NUMBER(c.value) FROM c WHERE IS_FINITE_NUMBER(c.amount)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsFiniteNumber,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "value"},
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
						Type: parsers.FunctionCallIsFiniteNumber,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "amount"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_INTEGER", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_INTEGER(c.number) FROM c WHERE IS_INTEGER(c.number)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsInteger,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "number"},
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
						Type: parsers.FunctionCallIsInteger,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "number"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_NULL", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_NULL(c.value) FROM c WHERE IS_NULL(c.value)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsNull,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "value"},
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
						Type: parsers.FunctionCallIsNull,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "value"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_NUMBER", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_NUMBER(c.value) FROM c WHERE IS_NUMBER(c.value)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsNumber,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "value"},
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
						Type: parsers.FunctionCallIsNumber,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "value"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_OBJECT", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_OBJECT(c.obj) FROM c WHERE IS_OBJECT(c.obj)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
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
				Filters: parsers.SelectItem{
					Type: parsers.SelectItemTypeFunctionCall,
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
		)
	})

	t.Run("Should parse function IS_PRIMITIVE", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_PRIMITIVE(c.value) FROM c WHERE IS_PRIMITIVE(c.value)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsPrimitive,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "value"},
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
						Type: parsers.FunctionCallIsPrimitive,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "value"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse function IS_STRING", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT IS_STRING(c.value) FROM c WHERE IS_STRING(c.value)`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeFunctionCall,
						Value: parsers.FunctionCall{
							Type: parsers.FunctionCallIsString,
							Arguments: []interface{}{
								parsers.SelectItem{
									Path: []string{"c", "value"},
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
						Type: parsers.FunctionCallIsString,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "value"},
								Type: parsers.SelectItemTypeField,
							},
						},
					},
				},
			},
		)
	})
}
