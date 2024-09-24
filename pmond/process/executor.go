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

func Exec(p *model.Process) (*model.Process, error) {

	user, groupIds, err := SetUser(p.Username)
	if err != nil {
		return nil, err
	}
	logPath, err := GetLogPath("", p.Log, p.ProcessFile, p.Name)
	if err != nil {
		return nil, err
	}
	logOutput, err := GetLogFile(logPath, *user)
	if err != nil {
		return nil, err
	}

	Env := os.Environ()

	if len(p.EnvVars) > 0 {
		Env = append(Env, strings.Fields(p.EnvVars)...)
	}

	lastSepIdx := strings.LastIndex(p.ProcessFile, string(os.PathSeparator))
	attr := &os.ProcAttr{
		Dir:   p.ProcessFile[0 : lastSepIdx+1],
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

	var processParams = []string{p.Name}
	if len(p.Args) > 0 {
		processParams = append(processParams, strings.Split(p.Args, " ")...)
	}

	process, err := os.StartProcess(p.ProcessFile, processParams, attr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pModel := model.Process{
		Pid:          uint32(process.Pid),
		Log:          logPath,
		Name:         p.Name,
		ProcessFile:  p.ProcessFile,
		Args:         strings.Join(processParams[1:], " "),
		EnvVars:      p.EnvVars,
		Pointer:      process,
		Status:       model.StatusInit,
		Uid:          conv.StrToUint32(user.Uid),
		Gid:          conv.StrToUint32(user.Gid),
		Username:     user.Username,
		AutoRestart:  p.AutoRestart,
		Dependencies: p.Dependencies,
		Groups:       p.Groups,
	}

	return &pModel, nil
}
