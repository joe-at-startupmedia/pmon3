package process

import (
	"os"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/array"
	"pmon3/pmond/utils/conv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func Exec(processFile string, customLogFile string, processName string, extArgs string, envVars string, username string, autoRestart bool, dependencies string) (*model.Process, error) {
	user, groupIds, err := SetUser(username)
	if err != nil {
		return nil, err
	}
	logPath, err := GetLogPath("", customLogFile, processFile, processName)
	if err != nil {
		return nil, err
	}
	logOutput, err := GetLogFile(logPath, *user)
	if err != nil {
		return nil, err
	}

	Env := os.Environ()

	if len(envVars) > 0 {
		Env = append(Env, strings.Fields(envVars)...)
	}

	lastSepIdx := strings.LastIndex(processFile, string(os.PathSeparator))
	attr := &os.ProcAttr{
		Dir:   processFile[0 : lastSepIdx+1],
		Env:   Env,
		Files: []*os.File{nil, logOutput, logOutput},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid:    conv.StrToUint32(user.Uid),
				Gid:    conv.StrToUint32(user.Gid),
				Groups: array.Map(groupIds, func(gid string) uint32 { return conv.StrToUint32(gid) }),
			},
			Setsid: true,
		},
	}

	var processParams = []string{processName}
	if len(extArgs) > 0 {
		processParams = append(processParams, strings.Split(extArgs, " ")...)
	}

	process, err := os.StartProcess(processFile, processParams, attr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pModel := model.Process{
		Pid:          uint32(process.Pid),
		Log:          logPath,
		Name:         processName,
		ProcessFile:  processFile,
		Args:         strings.Join(processParams[1:], " "),
		EnvVars:      envVars,
		Pointer:      process,
		Status:       model.StatusInit,
		Uid:          conv.StrToUint32(user.Uid),
		Gid:          conv.StrToUint32(user.Gid),
		Username:     user.Username,
		AutoRestart:  autoRestart,
		Dependencies: dependencies,
	}

	return &pModel, nil
}
