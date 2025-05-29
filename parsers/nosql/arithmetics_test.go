package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Parse_Arithmetics(t *testing.T) {
	t.Run("Should parse multiplication before addition", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.a + c.b * c.c FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "+",
							Left:      testutils.SelectItem_Path("c", "a"),
							Right: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left:      testutils.SelectItem_Path("c", "b"),
									Right:     testutils.SelectItem_Path("c", "c"),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse division before subtraction", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.x - c.y / c.z FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "-",
							Left:      testutils.SelectItem_Path("c", "x"),
							Right: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "/",
									Left:      testutils.SelectItem_Path("c", "y"),
									Right:     testutils.SelectItem_Path("c", "z"),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle complex mixed operations", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.a + c.b * c.c - c.d / c.e FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "-",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "+",
									Left:      testutils.SelectItem_Path("c", "a"),
									Right: parsers.SelectItem{
										Type: parsers.SelectItemTypeBinaryExpression,
										Value: parsers.BinaryExpression{
											Operation: "*",
											Left:      testutils.SelectItem_Path("c", "b"),
											Right:     testutils.SelectItem_Path("c", "c"),
										},
									},
								},
							},
							Right: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "/",
									Left:      testutils.SelectItem_Path("c", "d"),
									Right:     testutils.SelectItem_Path("c", "e"),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should respect parentheses overriding precedence", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT (c.a + c.b) * c.c FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "+",
									Left:      testutils.SelectItem_Path("c", "a"),
									Right:     testutils.SelectItem_Path("c", "b"),
								},
							},
							Right: testutils.SelectItem_Path("c", "c"),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle nested parentheses", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT ((c.a + c.b) * c.c) - c.d FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "-",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left: parsers.SelectItem{
										Type: parsers.SelectItemTypeBinaryExpression,
										Value: parsers.BinaryExpression{
											Operation: "+",
											Left:      testutils.SelectItem_Path("c", "a"),
											Right:     testutils.SelectItem_Path("c", "b"),
										},
									},
									Right: testutils.SelectItem_Path("c", "c"),
								},
							},
							Right: testutils.SelectItem_Path("c", "d"),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should be left associative for same precedence operators", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.a - c.b - c.c FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "-",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "-",
									Left:      testutils.SelectItem_Path("c", "a"),
									Right:     testutils.SelectItem_Path("c", "b"),
								},
							},
							Right: testutils.SelectItem_Path("c", "c"),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should be left associative with multiplication and division", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.a * c.b / c.c FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "/",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left:      testutils.SelectItem_Path("c", "a"),
									Right:     testutils.SelectItem_Path("c", "b"),
								},
							},
							Right: testutils.SelectItem_Path("c", "c"),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle math with constants", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT 10 + 20 * 5 FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "+",
							Left:      testutils.SelectItem_Constant_Int(10),
							Right: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left:      testutils.SelectItem_Constant_Int(20),
									Right:     testutils.SelectItem_Constant_Int(5),
								},
							},
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle math with floating point numbers", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.price * 1.08 FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left:      testutils.SelectItem_Path("c", "price"),
							Right:     testutils.SelectItem_Constant_Float(1.08),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle parentheses around single value", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT (c.value) FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					testutils.SelectItem_Path("c", "value"),
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle function calls in math expressions", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT LENGTH(c.name) * 2 + 10 FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "+",
							Left: parsers.SelectItem{
								Type: parsers.SelectItemTypeBinaryExpression,
								Value: parsers.BinaryExpression{
									Operation: "*",
									Left: parsers.SelectItem{
										Type: parsers.SelectItemTypeFunctionCall,
										Value: parsers.FunctionCall{
											Type:      parsers.FunctionCallLength,
											Arguments: []interface{}{testutils.SelectItem_Path("c", "name")},
										},
									},
									Right: testutils.SelectItem_Constant_Int(2),
								},
							},
							Right: testutils.SelectItem_Constant_Int(10),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle multiple select items with math", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.a + c.b, c.x * c.y FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "+",
							Left:      testutils.SelectItem_Path("c", "a"),
							Right:     testutils.SelectItem_Path("c", "b"),
						},
					},
					{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left:      testutils.SelectItem_Path("c", "x"),
							Right:     testutils.SelectItem_Path("c", "y"),
						},
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should handle math in WHERE clause", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id FROM c WHERE c.price * 1.08 > 100`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					testutils.SelectItem_Path("c", "id"),
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Filters: parsers.ComparisonExpression{
					Operation: ">",
					Left: parsers.SelectItem{
						Type: parsers.SelectItemTypeBinaryExpression,
						Value: parsers.BinaryExpression{
							Operation: "*",
							Left:      testutils.SelectItem_Path("c", "price"),
							Right:     testutils.SelectItem_Constant_Float(1.08),
						},
					},
					Right: testutils.SelectItem_Constant_Int(100),
				},
			},
		)
	})
}
