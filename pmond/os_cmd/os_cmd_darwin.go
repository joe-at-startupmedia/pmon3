package os_cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/utils/conv"
	"pmon3/utils/os_cmd"
	"strings"
	"syscall"
)

func findPidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s %s$' | grep -v grep | awk '{print $2}'", p.ProcessFile, p.Name, p.Args)
}

func findPidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $2}'", p.ProcessFile, p.Name)
}

func findPpidFromProcessNameAndArgs(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s %s$' | grep -v grep | awk '{print $3}'", p.ProcessFile, p.Name, p.Args)
}

func findPpidFromProcessName(p *model.Process) string {
	return fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $3}'", p.ProcessFile, p.Name)
}

func execIsRunning(p *model.Process) bool {
	var pid uint32
	if len(p.Args) > 0 {
		pid = os_cmd.GetUint32FromShellCommand(findPidFromProcessNameAndArgs(p))
	} else {
		pid = os_cmd.GetUint32FromShellCommand(findPidFromProcessName(p))
	}

	return pid > 0
}

func startProcess(p *model.Process, logFile *os.File, user *user.User, groupIds []string, envVars []string) (*os.Process, error) {

	attr := &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: conv.StrToUint32(user.Uid),
			Gid: conv.StrToUint32(user.Gid),
			//Groups: array.Map(groupIds, func(gid string) uint32 { return conv.StrToUint32(gid) }),//this errors
		},
		Setsid: true,
	}

	var processParams = []string{p.Name}
	if len(p.Args) > 0 {
		processParams = append(processParams, strings.Split(p.Args, " ")...)
	}

	pmond.Log.Infof("execCommand: %s %s %-v", p.ProcessFile, processParams, attr)

	cmd := exec.Command(p.ProcessFile, processParams...)

	cmd.SysProcAttr = attr
	cmd.Env = envVars
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the child process
	err := cmd.Start()

	return cmd.Process, err
}
