package table_one

import (
	"pmon3/cli/output/process/list"
)

// render single one table row
func Render(data []string) {
	table_list.Render([][]string{data})
}
