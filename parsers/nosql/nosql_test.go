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
// 	res, err := nosql.Parse("", []byte("select c.id, c._self AS self, c._rid, c._ts FROM c where c.id=\"12345\" AND c.pk=123"))
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
	t.Run("Shoul parse simple SELECT", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c.pk FROM c`,
			parsers.SelectStmt{
				Columns: []parsers.FieldPath{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
			},
		)
	})

	t.Run("Shoul parse SELECT with single WHERE condition", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
			FROM c
			WHERE c.pk=true`,
			parsers.SelectStmt{
				Columns: []parsers.FieldPath{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{Value: "c"},
				Filters: parsers.ComparisonExpression{
					Operation: "=",
					Left:      parsers.FieldPath{Path: []string{"c", "pk"}},
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
				Columns: []parsers.FieldPath{
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
							Left:      parsers.FieldPath{Path: []string{"c", "id"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeString, Value: "12345"},
						},
						parsers.ComparisonExpression{
							Operation: "=",
							Left:      parsers.FieldPath{Path: []string{"c", "pk"}},
							Right:     parsers.Constant{Type: parsers.ConstantTypeInteger, Value: 123},
						},
					},
				},
			},
		)
	})

	t.Run("Shoul correctly parse literals in conditions", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
			FROM c
			WHERE c.boolean=true
				AND c.integer=1
				AND c.float=6.9
				AND c.string="hello"`,
			parsers.SelectStmt{
				Columns: []parsers.FieldPath{{Path: []string{"c", "id"}, Alias: ""}},
				Table:   parsers.Table{Value: "c"},
				Filters: parsers.LogicalExpression{
					Expressions: []interface{}{
						parsers.ComparisonExpression{
							Left:      parsers.FieldPath{Path: []string{"c", "boolean"}},
							Right:     parsers.Constant{Type: 3, Value: true},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.FieldPath{Path: []string{"c", "integer"}},
							Right:     parsers.Constant{Type: 1, Value: 1},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.FieldPath{Path: []string{"c", "float"}},
							Right:     parsers.Constant{Type: 2, Value: 6.9},
							Operation: "=",
						},
						parsers.ComparisonExpression{
							Left:      parsers.FieldPath{Path: []string{"c", "string"}},
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
