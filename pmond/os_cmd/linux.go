//go:build linux

package os_cmd

import (
	"errors"
	"fmt"
	"github.com/goinbox/shell"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/conv"
	"strings"
)

func ExecFindPidFromProcessNameAndArgs(p *model.Process) uint32 {
	return getUint32FromShellCommand(FindPidFromProcessNameAndArgs(p))
}

func FindPidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $2}'", p.Name, p.Args)
}

func ExecFindPidFromProcessName(p *model.Process) uint32 {
	return getUint32FromShellCommand(FindPidFromProcessName(p))
}

func FindPidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $2}'", p.Name)
}

func ExecFindPpidFromProcessNameAndArgs(p *model.Process) uint32 {
	return getUint32FromShellCommand(FindPpidFromProcessNameAndArgs(p))
}

func FindPpidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $3}'", p.Name, p.Args)
}

func ExecFindPpidFromProcessName(p *model.Process) uint32 {
	return getUint32FromShellCommand(FindPpidFromProcessName(p))
}

func FindPpidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $3}'", p.Name)
}

func KillProcess(p *model.Process) string {
	return fmt.Sprintf("kill %s", p.GetPidStr())
}

func ExecKillProcess(p *model.Process) error {
	return getErrorFromShellCommand(KillProcess(p))
}

func KillProcessForcefully(p *model.Process) string {
	return fmt.Sprintf("kill -9 %s", p.GetPidStr())
}

func ExecKillProcessForcefully(p *model.Process) error {
	return getErrorFromShellCommand(KillProcessForcefully(p))
}

func HandleOnEventExec(cmdString string) {
	if err := getErrorFromShellCommand(cmdString); err != nil {
		pmond.Log.Errorf("event executor encountered an err: %s", err)
	}
}

func getUint32FromShellCommand(cmdString string) uint32 {

	var rel *shell.ShellResult

	rel = shell.RunCmd(cmdString)
	if rel.Ok {
		newPpidStr := strings.TrimSpace(string(rel.Output))
		newPpid := conv.StrToUint32(newPpidStr)
		return newPpid
	}

	return 0
}

func getErrorFromShellCommand(cmdString string) error {

	var rel *shell.ShellResult

	rel = shell.RunCmd(cmdString)

	if !rel.Ok {
		return errors.New(string(rel.Output))
	}

	return nil
}
