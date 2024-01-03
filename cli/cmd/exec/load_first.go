package exec

import (
	"pmon3/pmond"
	"pmon3/pmond/executor"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/crypto"
	"pmon3/pmond/worker"
	"strings"
)

func loadFirst(processFile string, flags *model.ExecFlags) error {

	logPath, err := executor.GetLogPath(flags.Log, crypto.Crc32Hash(processFile), flags.LogDir)
	if err != nil {
		return err
	}

	var processParams = []string{flags.Name}
	if len(flags.Args) > 0 {
		processParams = append(processParams, strings.Split(flags.Args, " ")...)
	}

	user, err := worker.GetProcUser(flags)
	if err != nil {
		return err
	}

	process := model.Process{
		Pid:         0,
		Log:         logPath,
		Name:        flags.Name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     nil,
		Status:      model.StatusQueued,
		Uid:         user.Uid,
		Gid:         user.Gid,
		Username:    user.Username,
		AutoRestart: !flags.NoAutoRestart,
	}

	err = pmond.Db().Save(&process).Error

	return err
}
