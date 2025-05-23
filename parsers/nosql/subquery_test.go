package nosql_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Parse_SubQuery(t *testing.T) {

	t.Run("Should parse FROM subquery", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id FROM (SELECT VALUE cc["info"] FROM cc) AS c`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
				},
				Table: parsers.Table{
					Value: "c",
					SelectItem: parsers.SelectItem{
						Alias: "c",
						Type:  parsers.SelectItemTypeSubQuery,
						Value: parsers.SelectStmt{
							Table: parsers.Table{SelectItem: testutils.SelectItem_Path("cc")},
							SelectItems: []parsers.SelectItem{
								{Path: []string{"cc", "info"}, IsTopLevel: true},
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse JOIN subquery", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id, cc.name FROM c JOIN (SELECT tag.name FROM tag IN c.tags) AS cc`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"cc", "name"}},
				},
				Table: parsers.Table{SelectItem: testutils.SelectItem_Path("c")},
				JoinItems: []parsers.JoinItem{
					{
						Table: parsers.Table{
							Value: "cc",
						},
						SelectItem: parsers.SelectItem{
							Alias: "cc",
							Type:  parsers.SelectItemTypeSubQuery,
							Value: parsers.SelectStmt{
								SelectItems: []parsers.SelectItem{
									testutils.SelectItem_Path("tag", "name"),
								},
								Table: parsers.Table{
									Value:      "tag",
									SelectItem: testutils.SelectItem_Path("c", "tags"),
									IsInSelect: true,
								},
							},
						},
					},
				},
			},
		)
	})

	t.Run("Should parse JOIN EXISTS subquery", func(t *testing.T) {
		testQueryParse(
			t,
			`SELECT c.id
			FROM c
			JOIN (
				SELECT VALUE EXISTS(SELECT tag.name FROM tag IN c.tags)
			) AS hasTags
			WHERE hasTags`,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					testutils.SelectItem_Path("c", "id"),
				},
				Table: parsers.Table{
					SelectItem: testutils.SelectItem_Path("c"),
				},
				JoinItems: []parsers.JoinItem{
					{
						Table: parsers.Table{Value: "hasTags"},
						SelectItem: parsers.SelectItem{
							Alias: "hasTags",
							Type:  parsers.SelectItemTypeSubQuery,
							Value: parsers.SelectStmt{
								SelectItems: []parsers.SelectItem{
									{
										IsTopLevel: true,
										Type:       parsers.SelectItemTypeSubQuery,
										Value: parsers.SelectStmt{
											SelectItems: []parsers.SelectItem{
												testutils.SelectItem_Path("tag", "name"),
											},
											Table: parsers.Table{
												Value:      "tag",
												SelectItem: testutils.SelectItem_Path("c", "tags"),
												IsInSelect: true,
											},
											Exists: true,
										},
									},
								},
							},
						},
					},
				},
				Filters: parsers.SelectItem{
					Path: []string{"hasTags"},
				},
			},
		)
	})
}
