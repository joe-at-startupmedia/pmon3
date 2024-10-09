package shell

import (
	"fmt"
	"os"
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/utils/array"
	"pmon3/utils/conv"
	"strings"
	"syscall"
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

func execIsRunning(p *model.Process) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%s/status", p.GetPidStr()))

	//if it doesn't exist in proc/n/status ask the OS
	if err != nil {
		//it's running if it exists
		return !os.IsNotExist(err)
	}

	return true
}

func startProcess(p *model.Process, logFile *os.File, user *user.User, groupIds []string, envVars []string) (*os.Process, error) {
	lastSepIdx := strings.LastIndex(p.ProcessFile, string(os.PathSeparator))
	attr := &os.ProcAttr{
		Dir:   p.ProcessFile[0 : lastSepIdx+1],
		Env:   envVars,
		Files: []*os.File{nil, logFile, logFile},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid:    conv.StrToUint32(user.Uid),
				Gid:    conv.StrToUint32(user.Gid),
				Groups: array.Map(groupIds, func(gid string) uint32 { return conv.StrToUint32(gid) }),
			},
			Setsid: true,
		},
	}

	var processParams = []string{p.Name}
	if len(p.Args) > 0 {
		processParams = append(processParams, strings.Split(p.Args, " ")...)
	}

	pmond.Log.Infof("os.StartProcess: %s %s %-v", p.ProcessFile, processParams, attr)
	return os.StartProcess(p.ProcessFile, processParams, attr)
}
