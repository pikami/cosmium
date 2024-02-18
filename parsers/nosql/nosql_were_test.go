package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Parse_Were(t *testing.T) {

	t.Run("Should parse SELECT with single WHERE condition", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
		FROM c
		WHERE c.isCool=true`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.ComparisonExpression{
					Operation: "=",
					Left:      parsers.SelectItem{Path: []string{"c", "isCool"}},
					Right: parsers.SelectItem{
						Type:  parsers.SelectItemTypeConstant,
						Value: parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
					},
				},
			},
		)
	})

	t.Run("Should parse SELECT with multiple WHERE conditions", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id, c._self AS self, c._rid, c._ts
		FROM c
		WHERE c.id="12345" OR c.pk=123`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "_self"}, Alias: "self"},
					{Path: []string{"c", "_rid"}},
					{Path: []string{"c", "_ts"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Operation: parsers.LogicalExpressionTypeOr,
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "id"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "12345"},
							},
						},
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "pk"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 123},
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse SELECT with grouped WHERE conditions", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
		FROM c
		WHERE c.isCool=true AND (c.id = "123" OR c.id = "456")`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Operation: parsers.LogicalExpressionTypeAnd,
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "isCool"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
							},
						},
						parsers.LogicalExpression{
							Operation: parsers.LogicalExpressionTypeOr,
							Expressions: []interface{}{
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right: parsers.SelectItem{
										Type:  parsers.SelectItemTypeConstant,
										Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "123"},
									},
								},
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right: parsers.SelectItem{
										Type:  parsers.SelectItemTypeConstant,
										Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "456"},
									},
								},
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should correctly parse literals in conditions", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
		FROM c
		WHERE c.boolean=true
			AND c.integer=1
			AND c.float=6.9
			AND c.string="hello"
			AND c.param=@param_id1`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{{Path: []string{"c", "id"}, Alias: ""}},
				Table:       parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Left: parsers.SelectItem{Path: []string{"c", "boolean"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
							},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left: parsers.SelectItem{Path: []string{"c", "integer"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 1},
							},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left: parsers.SelectItem{Path: []string{"c", "float"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeFloat, Value: 6.9},
							},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left: parsers.SelectItem{Path: []string{"c", "string"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeString, Value: "hello"},
							},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left: parsers.SelectItem{Path: []string{"c", "param"}},
							Right: parsers.SelectItem{
								Type:  parsers.SelectItemTypeConstant,
								Value: parsers.Constant{Type: parsers.ConstantTypeParameterConstant, Value: "@param_id1"},
							},
							Operation: "=",
						},
					},
					Operation: parsers.LogicalExpressionTypeAnd,
				},
			},
		)
	})
}
