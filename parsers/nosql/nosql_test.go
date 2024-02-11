package nosql_test

import (
	"log"
	"reflect"
	"testing"

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

func testQueryParse(t *testing.T, query string, expectedQuery nosql.SelectStmt) {
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
			nosql.SelectStmt{
				Columns: []nosql.FieldPath{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: nosql.Table{Value: "c"},
			},
		)
	})

	t.Run("Shoul parse SELECT with single WHERE condition", func(t *testing.T) {
		testQueryParse(
			t,
			`select c.id
			FROM c
			WHERE c.pk=true`,
			nosql.SelectStmt{
				Columns: []nosql.FieldPath{
					{Path: []string{"c", "id"}},
				},
				Table: nosql.Table{Value: "c"},
				Filters: nosql.ComparisonExpression{
					Operation: "=",
					Left:      nosql.FieldPath{Path: []string{"c", "pk"}},
					Right:     nosql.Constant{Type: nosql.ConstantTypeBoolean, Value: true},
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
			nosql.SelectStmt{
				Columns: []nosql.FieldPath{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "_self"}, Alias: "self"},
					{Path: []string{"c", "_rid"}},
					{Path: []string{"c", "_ts"}},
				},
				Table: nosql.Table{Value: "c"},
				Filters: nosql.LogicalExpression{
					Operation: nosql.LogicalExpressionTypeOr,
					Expressions: []interface{}{
						nosql.ComparisonExpression{
							Operation: "=",
							Left:      nosql.FieldPath{Path: []string{"c", "id"}},
							Right:     nosql.Constant{Type: nosql.ConstantTypeString, Value: "12345"},
						},
						nosql.ComparisonExpression{
							Operation: "=",
							Left:      nosql.FieldPath{Path: []string{"c", "pk"}},
							Right:     nosql.Constant{Type: nosql.ConstantTypeInteger, Value: 123},
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
			nosql.SelectStmt{
				Columns: []nosql.FieldPath{{Path: []string{"c", "id"}, Alias: ""}},
				Table:   nosql.Table{Value: "c"},
				Filters: nosql.LogicalExpression{
					Expressions: []interface{}{
						nosql.ComparisonExpression{
							Left:      nosql.FieldPath{Path: []string{"c", "boolean"}},
							Right:     nosql.Constant{Type: 3, Value: true},
							Operation: "=",
						},
						nosql.ComparisonExpression{
							Left:      nosql.FieldPath{Path: []string{"c", "integer"}},
							Right:     nosql.Constant{Type: 1, Value: 1},
							Operation: "=",
						},
						nosql.ComparisonExpression{
							Left:      nosql.FieldPath{Path: []string{"c", "float"}},
							Right:     nosql.Constant{Type: 2, Value: 6.9},
							Operation: "=",
						},
						nosql.ComparisonExpression{
							Left:      nosql.FieldPath{Path: []string{"c", "string"}},
							Right:     nosql.Constant{Type: 0, Value: "hello"},
							Operation: "=",
						},
					},
					Operation: nosql.LogicalExpressionTypeAnd,
				},
			},
		)
	})
}
