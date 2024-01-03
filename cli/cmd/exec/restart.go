package exec

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/executor"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/crypto"
	"pmon3/pmond/worker"
	"strings"
)

func Restart(m *model.Process, processFile string, flags *model.ExecFlags) error {
	// only stopped and failed process can be restart
	if m.Status == model.StatusStopped || m.Status == model.StatusFailed {
		if len(flags.Log) > 0 || len(flags.LogDir) > 0 {
			logPath, err := executor.GetLogPath(flags.Log, crypto.Crc32Hash(processFile), flags.LogDir)
			if err != nil {
				return err
			}
			m.Log = logPath
		}

		if len(flags.Args) > 0 {
			var processParams = []string{flags.Name}
			processParams = append(processParams, strings.Split(flags.Args, " ")...)
			m.Args = strings.Join(processParams[1:], " ")
		}

		if len(flags.User) > 0 {
			user, err := worker.GetProcUser(flags)
			if err != nil {
				return err
			}
			m.Uid = user.Uid
			m.Gid = user.Gid
			m.Username = user.Username
		}

		if flags.NoAutoRestart {
			m.AutoRestart = !flags.NoAutoRestart
		}

		m.Status = model.StatusQueued
		m.ProcessFile = processFile

		return pmond.Db().Save(&m).Error
	}

	return fmt.Errorf("process already running with the name provided: %s", m.Name)
}
