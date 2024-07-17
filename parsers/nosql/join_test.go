package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
)

func Test_Parse_Join(t *testing.T) {

	t.Run("Should parse simple JOIN", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, c["pk"] FROM c JOIN cc IN c["tags"]`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"c", "pk"}},
				},
				Table: parsers.Table{Value: "c"},
				JoinItems: []parsers.JoinItem{
					{
						Table: parsers.Table{
							Value: "cc",
						},
						SelectItem: parsers.SelectItem{
							Path: []string{"c", "tags"},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse JOIN VALUE", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT VALUE cc FROM c JOIN cc IN c["tags"]`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"cc"}, IsTopLevel: true},
				},
				Table: parsers.Table{Value: "c"},
				JoinItems: []parsers.JoinItem{
					{
						Table: parsers.Table{
							Value: "cc",
						},
						SelectItem: parsers.SelectItem{
							Path: []string{"c", "tags"},
						},
					},
				},
			},
		)
	})
}
