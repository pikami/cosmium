package memoryexecutor_test

import (
	"testing"

	"github.com/pikami/cosmium/parsers"
	memoryexecutor "github.com/pikami/cosmium/query_executors/memory_executor"
)

func Test_Execute_Joins(t *testing.T) {
	mockData := []memoryexecutor.RowType{
		map[string]interface{}{
			"id": 1,
			"tags": []map[string]interface{}{
				{"name": "a"},
				{"name": "b"},
			},
		},
		map[string]interface{}{
			"id": 2,
			"tags": []map[string]interface{}{
				{"name": "b"},
				{"name": "c"},
			},
		},
	}

	t.Run("Should execute JOIN on 'tags'", func(t *testing.T) {
		testQueryExecute(
			t,
			parsers.SelectStmt{
				SelectItems: []parsers.SelectItem{
					{Path: []string{"c", "id"}},
					{Path: []string{"cc", "name"}},
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"id": 1, "name": "a"},
				map[string]interface{}{"id": 1, "name": "b"},
				map[string]interface{}{"id": 2, "name": "b"},
				map[string]interface{}{"id": 2, "name": "c"},
			},
		)
	})

	t.Run("Should execute JOIN VALUE on 'tags'", func(t *testing.T) {
		testQueryExecute(
			t,
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
			mockData,
			[]memoryexecutor.RowType{
				map[string]interface{}{"name": "a"},
				map[string]interface{}{"name": "b"},
				map[string]interface{}{"name": "b"},
				map[string]interface{}{"name": "c"},
			},
		)
	})
}
