package executor

import (
	"os"
	"os/user"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/crypto"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func Exec(processFile, customLogFile, name, extArgs string, user *user.User, autoRestart bool, logDir string) (*model.Process, error) {
	logPath, err := GetLogPath(customLogFile, crypto.Crc32Hash(processFile), logDir)
	if err != nil {
		return nil, err
	}
	logOutput, err := GetLogFile(logPath)
	if err != nil {
		return nil, err
	}

	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	lastSepIdx := strings.LastIndex(processFile, string(os.PathSeparator))
	attr := &os.ProcAttr{
		Dir:   processFile[0 : lastSepIdx+1],
		Env:   os.Environ(),
		Files: []*os.File{nil, logOutput, logOutput},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
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
		Pid:         process.Pid,
		Log:         logPath,
		Name:        name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     process,
		Status:      model.StatusInit,
		Uid:         user.Uid,
		Gid:         user.Gid,
		Username:    user.Username,
		AutoRestart: autoRestart,
	}

	return &pModel, nil
}
