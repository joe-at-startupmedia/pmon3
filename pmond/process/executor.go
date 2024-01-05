package process

import (
	"os"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/conv"
	"pmon3/pmond/utils/crypto"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func Exec(processFile, customLogFile, name, extArgs string, username string, autoRestart bool) (*model.Process, error) {
	user, err := SetUser(username)
	if err != nil {
		return nil, err
	}
	logPath, err := GetLogPath(customLogFile, crypto.Crc32Hash(processFile), "")
	if err != nil {
		return nil, err
	}
	logOutput, err := GetLogFile(logPath)
	if err != nil {
		return nil, err
	}

	lastSepIdx := strings.LastIndex(processFile, string(os.PathSeparator))
	attr := &os.ProcAttr{
		Dir:   processFile[0 : lastSepIdx+1],
		Env:   os.Environ(),
		Files: []*os.File{nil, logOutput, logOutput},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: conv.StrToUint32(user.Uid),
				Gid: conv.StrToUint32(user.Gid),
			},
			Setsid: true,
		},
	}

	var processParams = []string{name}
	if len(extArgs) > 0 {
		processParams = append(processParams, strings.Split(extArgs, " ")...)
	}

	process, err := os.StartProcess(processFile, processParams, attr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pModel := model.Process{
		Pid:         uint32(process.Pid),
		Log:         logPath,
		Name:        name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     process,
		Status:      model.StatusInit,
		Uid:         conv.StrToUint32(user.Uid),
		Gid:         conv.StrToUint32(user.Gid),
		Username:    user.Username,
		AutoRestart: autoRestart,
	}

	return &pModel, nil
}
