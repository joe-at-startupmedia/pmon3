package os_cmd

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/utils/conv"
	"pmon3/utils/os_cmd"
	"strings"
)

func findPidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $2}'", p.Name, p.Args)
}

func findPidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $2}'", p.Name)
}

func findPpidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $3}'", p.Name, p.Args)
}

func findPpidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $3}'", p.Name)
}

func killProcess(p *model.Process) string {
	return fmt.Sprintf("kill %s", p.GetPidStr())
}

func killProcessForcefully(p *model.Process) string {
	return fmt.Sprintf("kill -9 %s", p.GetPidStr())
}

func isPmondRunning(pid string) string {
	return fmt.Sprintf("ps -e -H -o pid,comm | awk '$2 ~ /pmond/ { print $1}' | grep -v %s | head -n 1", pid)
}

func execIsPmondRunning(pid string) bool {
	rel, _ := os_cmd.GetResultWithErrorFromShellCommand(isPmondRunning(pid))
	if rel.Ok {
		pmond.Log.Debugf("%s", string(rel.Output))
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid != 0
	}
	return false
}
