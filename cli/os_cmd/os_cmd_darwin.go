package os_cmd

import (
	"fmt"
)

func topCmd(pidArr []string, sortField string, refreshInterval int) []string {

	pidArrLen := len(pidArr)
	pidArgs := "-pid " + pidArr[0]
	for i := range pidArrLen {
		if i == 0 {
			continue
		}
		pidArgs = fmt.Sprintf("%s -pid %s", pidArgs, pidArr[i])
	}

	return []string{
		pidArgs,
		"-o",
		sortField,
		"-l 1",
		"-ncols 13",
	}
}
