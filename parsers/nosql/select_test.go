package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Parse_Select(t *testing.T) {

	t.Run("Should parse simple SELECT", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c["pk"] FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse SELECT with query parameters as accessor", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c[@param] FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "@param"}},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse SELECT DISTINCT", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT DISTINCT c.id FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table:    parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Distinct: true,
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Count: 1,
			},
		)
	})

	t.Run("Should parse SELECT OFFSET", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id FROM c OFFSET 3 LIMIT 5`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table:  parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				Count:  5,
				Offset: 3,
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse SELECT *", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT * FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c"}, IsTopLevel: true},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse SELECT c", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c FROM c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c"}, IsTopLevel: false},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})

	t.Run("Should parse SELECT with alias", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT
			  c.id AS aliasWithAs,
			  c.pk aliasWithoutAs
			FROM root c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Alias: "aliasWithAs", Path: []string{"c", "id"}},
					{Alias: "aliasWithoutAs", Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{
					Value:      "c",
					SelectItem: parsers.SelectItem{Alias: "c", Path: []string{"root"}},
				},
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
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
			},
		)
	})
}
