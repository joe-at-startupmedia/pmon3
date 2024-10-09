package os_cmd

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/utils/os_cmd"
)

func ExecFindPidFromProcessNameAndArgs(p *model.Process) uint32 {
	return os_cmd.GetUint32FromShellCommand(findPidFromProcessNameAndArgs(p))
}

func ExecFindPidFromProcessName(p *model.Process) uint32 {
	return os_cmd.GetUint32FromShellCommand(findPidFromProcessName(p))
}

func ExecFindPpidFromProcessNameAndArgs(p *model.Process) uint32 {
	return os_cmd.GetUint32FromShellCommand(findPpidFromProcessNameAndArgs(p))
}

func ExecFindPpidFromProcessName(p *model.Process) uint32 {
	return os_cmd.GetUint32FromShellCommand(findPpidFromProcessName(p))
}

func ExecKillProcess(p *model.Process) error {
	return os_cmd.GetErrorFromShellCommand(killProcess(p))
}

func ExecKillProcessForcefully(p *model.Process) error {
	return os_cmd.GetErrorFromShellCommand(killProcessForcefully(p))
}

func HandleOnEventExec(cmdString string) {
	if err := os_cmd.GetErrorFromShellCommand(cmdString); err != nil {
		pmond.Log.Errorf("event executor encountered an err: %s", err)
	}
}

func ExecIsPmondRunning(pid string) bool {
	return execIsPmondRunning(pid)
}
