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

	t.Run("Should parse SELECT with ORDER BY", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c["pk"] FROM c ORDER BY c.id DESC, c.pk`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
				OrderExpressions: []parsers.OrderExpression{
					{
						SelectItem: parsers.SelectItem{Path: []string{"c", "id"}},
						Direction:  parsers.OrderDirectionDesc,
					},
					{
						SelectItem: parsers.SelectItem{Path: []string{"c", "pk"}},
						Direction:  parsers.OrderDirectionAsc,
					},
				},
			},
		)
	})
}
