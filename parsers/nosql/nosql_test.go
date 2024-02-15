package nosql_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
)

// For Parser Debugging
// func Test_ParseTest(t *testing.T) {
// 	// select c.id, c._self, c._rid, c._ts, [c[\"pk\"]] as _partitionKeyValue from c
// 	res, err := nosql.Parse("", []byte("SELECT VALUE c.id FROM c"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	result, err := json.MarshalIndent(res, "", "   ")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Printf("output:\n%v\n", string(result))
// }

func testQueryParse(t *testing.T, query string, expectedQuery parsers.SelectStmt) {
	parsedQuery, err := nosql.Parse("", []byte(query))
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(parsedQuery, expectedQuery) {
		t.Errorf("parsed query does not match expected structure.\nExpected: %+v\nGot: %+v", expectedQuery, parsedQuery)
	}
}

func Test_Parse(t *testing.T) {
	t.Run("Should parse simple SELECT", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c["pk"] FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse SELECT TOP", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT TOP 1 c.id FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Count: 1,
			},
		)
	})

	t.Run("Should parse SELECT VALUE", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT VALUE c.id FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}, IsTopLevel: true},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse SELECT array", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT [c.id, c.pk] as arr FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "arr",
						Type:  parsers.SelectItemTypeArray,
						SelectItems: []parsers.SelectItem{
							{Path: []string{"c", "id"}},
							{Path: []string{"c", "pk"}},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Should parse SELECT object", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT { id: c.id, _pk: c.pk } AS obj FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Alias: "obj",
						Type:  parsers.SelectItemTypeObject,
						SelectItems: []parsers.SelectItem{
							{Alias: "id", Path: []string{"c", "id"}},
							{Alias: "_pk", Path: []string{"c", "pk"}},
						},
					},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

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
					Right:     parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
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
							Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "12345"},
						},
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.SelectItem{Path: []string{"c", "pk"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 123},
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
							Right:     parsers.Constant{Type: parsers.ConstantTypeBoolean, Value: true},
						},
						parsers.LogicalExpression{
							Operation: parsers.LogicalExpressionTypeOr,
							Expressions: []interface{}{
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "123"},
								},
								parsers.ComparisonExpression{
									Operation: "=",
									Left:      parsers.SelectItem{Path: []string{"c", "id"}},
									Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "456"},
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
				AND c.string="hello"`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{{Path: []string{"c", "id"}, Alias: ""}},
				Table:       parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Left:      parsers.SelectItem{Path: []string{"c", "boolean"}},
							Right:     parsers.Constant{Type: 3, Value: true},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.SelectItem{Path: []string{"c", "integer"}},
							Right:     parsers.Constant{Type: 1, Value: 1},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.SelectItem{Path: []string{"c", "float"}},
							Right:     parsers.Constant{Type: 2, Value: 6.9},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.SelectItem{Path: []string{"c", "string"}},
							Right:     parsers.Constant{Type: 0, Value: "hello"},
							Operation: "=",
						},
					},
					Operation: parsers.LogicalExpressionTypeAnd,
				},
			},
		)
	})
}
