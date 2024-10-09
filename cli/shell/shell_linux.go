package shell

import (
	"strconv"
	"strings"
)

func topCmd(pidArr []string, sortField string, refreshInterval int) []string {
	return []string{
		"-p",
		strings.Join(pidArr, ","),
		"-o",
		"%" + strings.ToUpper(sortField),
		"-d",
		strconv.Itoa(refreshInterval),
		"-b",
	}
}
