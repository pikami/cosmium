package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
	testutils "github.com/pikami/cosmium/test_utils"
)

func Test_Execute_SubQuery(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{"id": "123", "info": map[string]interface{}{"name": "row-1"}},
		map[string]interface{}{
			"id":   "456",
			"info": map[string]interface{}{"name": "row-2"},
			"tags": []map[string]interface{}{
				{"name": "tag-a"},
				{"name": "tag-b"},
			},
		},
		map[string]interface{}{
			"id":   "789",
			"info": map[string]interface{}{"name": "row-3"},
			"tags": []map[string]interface{}{
				{"name": "tag-b"},
				{"name": "tag-c"},
			},
		},
	}

	t.Run("Should execute FROM subquery", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "name"}},
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"name": "row-1"},
				map[string]interface{}{"name": "row-2"},
				map[string]interface{}{"name": "row-3"},
			},
		)
	})

	t.Run("Should execute JOIN subquery", func(t *testing.T) {
		testQueryExecute(
			t,
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "456", "name": "tag-a"},
				map[string]interface{}{"id": "456", "name": "tag-b"},
				map[string]interface{}{"id": "789", "name": "tag-b"},
				map[string]interface{}{"id": "789", "name": "tag-c"},
			},
		)
	})

	t.Run("Should execute JOIN EXISTS subquery", func(t *testing.T) {
		testQueryExecute(
			t,
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": "456"},
				map[string]interface{}{"id": "789"},
			},
		)
	})
}
