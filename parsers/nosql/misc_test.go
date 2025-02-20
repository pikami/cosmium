package nosql_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/pikami/cosmium/parsers"
	"github.com/pikami/cosmium/parsers/nosql"
	testutils "github.com/pikami/cosmium/test_utils"
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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

	t.Run("Should parse SELECT with GROUP BY", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c["pk"] FROM c GROUP BY c.id, c.pk`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				GroupBy: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
			},
		)
	})

	t.Run("Should parse IN function", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id FROM c WHERE c.id IN ("123", "456")`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Filters: parsers.SelectItem{
					Type: parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type: parsers.FunctionCallIn,
						Arguments: []interface{}{
							parsers.SelectItem{
								Path: []string{"c", "id"},
								Type: parsers.SelectItemTypeField,
							},
							testutils.SelectItem_Constant_String("123"),
							testutils.SelectItem_Constant_String("456"),
						},
					},
				},
			},
		)
	})

	t.Run("Should parse IN function with function call", func(t *testing.T) {
		testQueryParse(
			t,
			`Select c.id FROM c WHERE (ToString(c.id) IN ("123", "456"))`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					testutils.SelectItem_Path("c", "id"),
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Filters: parsers.SelectItem{
					Type: parsers.SelectItemTypeFunctionCall,
					Value: parsers.FunctionCall{
						Type: parsers.FunctionCallIn,
						Arguments: []interface{}{
							parsers.SelectItem{
								Type: parsers.SelectItemTypeFunctionCall,
								Value: parsers.FunctionCall{
									Type: parsers.FunctionCallToString,
									Arguments: []interface{}{
										testutils.SelectItem_Path("c", "id"),
									},
								},
							},
							testutils.SelectItem_Constant_String("123"),
							testutils.SelectItem_Constant_String("456"),
						},
					},
				},
			},
		)
	})

	t.Run("Should parse IN selector", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id FROM c IN c.tags`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{
						Path: []string{"c", "id"},
						Type: parsers.SelectItemTypeField,
					},
				},
				Table: parsers.Table{
					Value:      "c",
					SelectItem: testutils.SelectItem_Path("c", "tags"),
					IsInSelect: true,
				},
			},
		)
	})
}
