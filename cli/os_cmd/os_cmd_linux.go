package os_cmd

import (
	"fmt"
	"os/exec"
	"pmon3/cli"
	"pmon3/utils/conv"
	"pmon3/utils/os_cmd"
	"strconv"
	"strings"
)

func tailLogFile(logFileName string, numLines string) string {
	return fmt.Sprintf("tail %s -n %s", logFileName, numLines)
}

func tailFLogFile(logFileName string, numLines string) string {
	return fmt.Sprintf("tail -f %s -n %s", logFileName, numLines)
}

func catArchivedLogs(logFileName string) string {
	return fmt.Sprintf("zcat -v %s*.gz", logFileName)
}

func topCmd(pidArr []string, sortField string, refreshInterval int) []string {
	return []string{
		"-p",
		strings.Join(pidArr, ","),
		"-o",
		sortField,
		"-d",
		strconv.Itoa(refreshInterval),
		"-b",
	}
}

func isPmondRunning() string {
	return "ps -e -H -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | head -n 1"
}

func execTopCmd(pidArr []string, sortField string, refreshInterval int) *exec.Cmd {
	return exec.Command("top", topCmd(pidArr, sortField, refreshInterval)...)
}

func execTailLogFile(logFileName string, numLines string) *exec.Cmd {
	return os_cmd.ExecBashCommand(tailLogFile(logFileName, numLines))
}

func execTailFLogFile(logFileName string, numLines string) *exec.Cmd {
	return os_cmd.ExecBashCommand(tailFLogFile(logFileName, numLines))
}

func execCatArchivedLogs(logFileName string) *exec.Cmd {
	return os_cmd.ExecBashCommand(catArchivedLogs(logFileName))
}

func execIsPmondRunning() bool {
	rel, _ := os_cmd.GetResultWithErrorFromShellCommand(isPmondRunning())
	if rel.Ok {
		cli.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}
