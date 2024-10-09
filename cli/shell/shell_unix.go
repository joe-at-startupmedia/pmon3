package shell

import (
	"fmt"
	"os/exec"
	"pmon3/cli"
	"pmon3/utils/conv"
	"pmon3/utils/shell"
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

func isPmondRunning() string {
	return "ps -e -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | head -n 1"
}

func execTopCmd(pidArr []string, sortField string, refreshInterval int) *exec.Cmd {
	return exec.Command("top", topCmd(pidArr, sortField, refreshInterval)...)
}

func execTailLogFile(logFileName string, numLines string) *exec.Cmd {
	return shell.ExecBashCommand(tailLogFile(logFileName, numLines))
}

func execTailFLogFile(logFileName string, numLines string) *exec.Cmd {
	return shell.ExecBashCommand(tailFLogFile(logFileName, numLines))
}

func execCatArchivedLogs(logFileName string) *exec.Cmd {
	return shell.ExecBashCommand(catArchivedLogs(logFileName))
}

func execIsPmondRunning() bool {
	rel, _ := shell.GetResultWithErrorFromShellCommand(isPmondRunning())
	if rel.Ok {
		cli.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}
