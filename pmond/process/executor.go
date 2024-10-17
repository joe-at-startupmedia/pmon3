package process

import (
	"github.com/pkg/errors"
	"os"
	"pmon3/model"
	"pmon3/pmond/shell"
	"pmon3/utils/conv"
	"strings"
)

func exec(p *model.Process) (*model.Process, error) {

	user, groupIds, err := SetUser(p.Username)
	if err != nil {
		return nil, err
	}
	logPath, err := GetLogPath("", p.Log, p.Name)
	if err != nil {
		return nil, err
	}
	logFile, err := GetLogFile(logPath, *user)
	if err != nil {
		return nil, err
	}

	envVars := os.Environ()

	if len(p.EnvVars) > 0 {
		envVars = append(envVars, strings.Fields(p.EnvVars)...)
	}

	process, err := shell.StartProcess(p, logFile, user, groupIds, envVars)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pModel := model.Process{
		Pid:          uint32(process.Pid),
		Log:          logPath,
		Name:         p.Name,
		ProcessFile:  p.ProcessFile,
		Args:         p.Args,
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
