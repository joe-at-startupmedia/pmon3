package restart

import (
	"fmt"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/controller/base/exec"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"pmon3/pmond/utils/conv"
	"strings"
)

func ByProcess(cmd *protos.Cmd, p *model.Process, idOrName string, flags string, incrementCounter bool) error {
	// kill the process and insert a new record with "queued" status

	//the process doesn't exist,  so we'll look in the AppConfig
	if p == nil {

		app, err := conf.GetAppByName(idOrName, pmond.Config.AppsConfig.Apps)
		if err != nil {
			return fmt.Errorf("command error: start process error: %w", err)
		}

		// get exec abs file path
		execPath, err := exec.GetExecFileAbsPath(app.File)
		if err != nil {
			return fmt.Errorf("command error: file argument error: %w", err)
		}

		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		if p, err = exec.InsertAsQueued(execPath, &app.Flags); err != nil {
			return fmt.Errorf("could not start process: %w", err)
		}

	} else {
		if process.IsRunning(p.Pid) {
			if err := process.SendOsKillSignal(p, model.StatusStopped, false); err != nil {
				return err
			}
		}
		execFlags := model.ExecFlags{}
		parsedFlags, err := execFlags.Parse(flags)

		if err != nil {
			return fmt.Errorf("could not parse flags: %w", err)
		} else {
			pmond.Log.Debugf("update as queued: %v", flags)
			err = UpdateAsQueued(p, p.ProcessFile, parsedFlags)
			if err != nil {
				return err
			} else if incrementCounter && strings.HasSuffix(cmd.GetName(), "restart") {
				p.IncrRestartCount()
			}
		}
	}

	return nil
}

func UpdateAsQueued(m *model.Process, processFile string, flags *model.ExecFlags) error {
	// only stopped and failed process can be restarted
	if m.Status != model.StatusStopped && m.Status != model.StatusFailed {
		return fmt.Errorf("process already running with the name provided: %s", m.Name)
	}
	if len(flags.Log) > 0 || len(flags.LogDir) > 0 {
		logPath, err := process.GetLogPath(flags.LogDir, flags.Log, processFile, m.Name)
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
	m.ProcessFile = processFile

	return repo.ProcessOf(m).Save()
}
