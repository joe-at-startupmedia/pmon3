package shell

import (
	"os"
	"os/user"
	"pmon3/model"
	"pmon3/pmond"
	"pmon3/utils/shell"
)

func ExecFindPidFromProcessNameAndArgs(p *model.Process) uint32 {
	return shell.GetUint32FromShellCommand(findPidFromProcessNameAndArgs(p))
}

func ExecFindPidFromProcessName(p *model.Process) uint32 {
	return shell.GetUint32FromShellCommand(findPidFromProcessName(p))
}

func ExecFindPpidFromProcessNameAndArgs(p *model.Process) uint32 {
	return shell.GetUint32FromShellCommand(findPpidFromProcessNameAndArgs(p))
}

func ExecFindPpidFromProcessName(p *model.Process) uint32 {
	return shell.GetUint32FromShellCommand(findPpidFromProcessName(p))
}

func ExecKillProcess(p *model.Process) error {
	return shell.GetErrorFromShellCommand(killProcess(p))
}

func ExecKillProcessForcefully(p *model.Process) error {
	return shell.GetErrorFromShellCommand(killProcessForcefully(p))
}

func HandleOnEventExec(cmdString string) {
	if err := shell.GetErrorFromShellCommand(cmdString); err != nil {
		pmond.Log.Errorf("event executor encountered an err: %s", err)
	}
}

func ExecIsPmondRunning(pid int) bool {
	return execIsPmondRunning(pid)
}

func ExecIsRunning(p *model.Process) bool {
	return execIsRunning(p)
}

func StartProcess(p *model.Process, logFile *os.File, user *user.User, groupIds []string, envVars []string) (*os.Process, error) {
	return startProcess(p, logFile, user, groupIds, envVars)
}
