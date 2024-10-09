package shell

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/utils/conv"
	"pmon3/utils/shell"
	"strings"
)

func isPmondRunning(pid int) string {
	return fmt.Sprintf("ps -e -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | grep -v %d | head -n 1", pid)
}

func killProcess(p *model.Process) string {
	return fmt.Sprintf("kill %s", p.GetPidStr())
}

func killProcessForcefully(p *model.Process) string {
	return fmt.Sprintf("kill -9 %s", p.GetPidStr())
}

func execIsPmondRunning(pid int) bool {
	rel, _ := shell.GetResultWithErrorFromShellCommand(isPmondRunning(pid))
	if rel.Ok {
		pmond.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}
