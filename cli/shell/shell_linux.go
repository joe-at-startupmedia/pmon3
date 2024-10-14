package shell

import (
	"strings"
)

func topCmd(pidArr []string, sortField string, refreshInterval int) []string {
	return []string{
		"-p",
		strings.Join(pidArr, ","),
		"-o",
		"%" + strings.ToUpper(sortField),
		"-b",
		"-n 1",
	}
}
