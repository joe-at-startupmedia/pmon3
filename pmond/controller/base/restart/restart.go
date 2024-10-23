package restart

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/controller/base/exec"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
	"pmon3/protos"
	"pmon3/utils/conv"
	"strings"
)

func setExecFileAbsPath(execFlags *model.ExecFlags) error {
	_, err := os.Stat(execFlags.File)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %w", execFlags.File, err)
	}

	if path.IsAbs(execFlags.File) {
		return nil
	}

	absPath, err := filepath.Abs(execFlags.File)
	if err != nil {
		return fmt.Errorf("get file path error: %w", err)
	}
	execFlags.File = absPath

	return nil
}

func ByProcess(cmd *protos.Cmd, p *model.Process, idOrName string, flags string, incrementCounter bool) (*model.Process, error) {
	// kill the process and insert a new record with "queued" status

	//the process doesn't exist,  so we'll look in the AppConfig
	if p == nil {

		execFlags, err := pmond.Config.ProcessConfig.GetExecFlagsByName(idOrName)
		if err != nil {
			return nil, fmt.Errorf("command error: start process error: %w", err)
		}

		// get exec abs file path
		err = setExecFileAbsPath(&execFlags)
		if err != nil {
			return nil, fmt.Errorf("command error: file argument error: %w", err)
		}

		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		if p, err = exec.InsertAsQueued(&execFlags); err != nil {
			return nil, fmt.Errorf("could not start process: %w", err)
		}

	} else {

		if _, err := process.BeginPendingTask(p, "restart"); err != nil {
			return nil, err
		}
		defer process.FinishPendingTask(p)

		if err := repo.ProcessOf(p).UpdateStatus(model.StatusRestarting); err != nil {
			return nil, err
		}

		if err := process.SendOsKillSignal(p, true); err != nil {
			return nil, err
		}

		execFlags := model.ExecFlags{}
		parsedFlags, err := execFlags.Parse(flags)
		if err != nil {
			return nil, fmt.Errorf("could not parse flags: %w", err)
		} else {
			parsedFlags.File = p.ProcessFile
			pmond.Log.Debugf("update as queued: %v", flags)
			err = UpdateAsQueued(p, parsedFlags)
			if err != nil {
				return nil, err
			} else if incrementCounter && strings.HasSuffix(cmd.GetName(), "restart") {
				p.IncrRestartCount()
			}
		}
	}

	return p, nil
}

func UpdateAsQueued(m *model.Process, flags *model.ExecFlags) error {
	// only stopped and failed process can be restarted
	if m.Status != model.StatusStopped && m.Status != model.StatusFailed && m.Status != model.StatusRestarting {
		return fmt.Errorf("process already running with the name provided: %s", m.Name)
	}
	if len(flags.Log) > 0 || len(flags.LogDir) > 0 {
		logPath, err := process.GetLogPath(flags.LogDir, flags.Log, m.Name)
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
		user, _, err := process.SetUser(flags.User)
		if err != nil {
			return err
		}
		m.Uid = conv.StrToUint32(user.Uid)
		m.Gid = conv.StrToUint32(user.Gid)
		m.Username = user.Username
	}

	if flags.NoAutoRestart {
		m.AutoRestart = !flags.NoAutoRestart
	}

	m.Status = model.StatusQueued
	m.ProcessFile = flags.File

	//allow updating dependencies on restart
	if len(flags.Dependencies) > 0 {
		m.Dependencies = strings.Join(flags.Dependencies, " ")
	}

	return repo.ProcessOf(m).Save()
}
