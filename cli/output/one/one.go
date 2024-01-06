package table_one

import table_list "pmon3/cli/output/list"

// render single one table row
func Render(data []string) {
	table_list.Render([][]string{data})
}
