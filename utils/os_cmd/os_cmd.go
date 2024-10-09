package os_cmd

import (
	"errors"
	"github.com/goinbox/shell"
	"os/exec"
	"pmon3/pmond"
	"pmon3/utils/conv"
	"strings"
)

func GetUint32FromShellCommand(cmdString string) uint32 {
	rel, _ := GetResultWithErrorFromShellCommand(cmdString)
	if rel.Ok {
		newPpidStr := strings.TrimSpace(string(rel.Output))
		newPpid := conv.StrToUint32(newPpidStr)
		return newPpid
	}
	return 0
}

func GetErrorFromShellCommand(cmdString string) error {
	_, err := GetResultWithErrorFromShellCommand(cmdString)
	return err
}

func GetResultWithErrorFromShellCommand(cmdString string) (*shell.ShellResult, error) {

	var rel *shell.ShellResult

	rel = shell.RunCmd(cmdString)

	if !rel.Ok {
		errString := string(rel.Output)
		pmond.Log.Warnf("%s errored with: %s", cmdString, errString)
		return nil, errors.New(errString)
	}

	return rel, nil
}

func ExecBashCommand(commandString string) *exec.Cmd {
	return exec.Command("bash", "-c", commandString)
}
