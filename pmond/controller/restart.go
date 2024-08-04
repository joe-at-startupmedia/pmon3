package controller

import (
	"fmt"
	"github.com/pkg/errors"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"strings"
)

func Restart(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	flags := cmd.GetArg2()
	return RestartByParams(cmd, idOrName, flags, true)
}

func RestartByParams(cmd *protos.Cmd, idOrName string, flags string, incrementCounter bool) *protos.CmdResp {
	// kill the process and insert a new record with "queued" status

	err, p := model.FindProcessByIdOrName(db.Db(), idOrName)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("could not find process: %w", err))
	}
	if process.IsRunning(p.Pid) {
		if err := process.SendOsKillSignal(p, model.StatusStopped, false); err != nil {
			return ErroredCmdResp(cmd, errors.New(err.Error()))
		}
	}
	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("could not parse flags: %w", err))
	} else {
		pmond.Log.Debugf("update as queued: %v", flags)
		err = UpdateAsQueued(p, p.ProcessFile, parsedFlags)
		newProcess := protos.Process{
			Log: p.Log,
		}
		newCmdResp := protos.CmdResp{
			Id:      cmd.GetId(),
			Name:    cmd.GetName(),
			Process: &newProcess,
		}
		if err != nil {
			newCmdResp.Error = err.Error()
		} else if incrementCounter {
			p.IncrRestartCount()
		}
		return &newCmdResp
	}
}

func UpdateAsQueued(m *model.Process, processFile string, flags *model.ExecFlags) error {
	// only stopped and failed process can be restarted
	if m.Status != model.StatusStopped && m.Status != model.StatusFailed {
		return fmt.Errorf("process already running with the name provided: %s", m.Name)
	}
	if len(flags.Log) > 0 || len(flags.LogDir) > 0 {
		logPath, err := process.GetLogPath(flags.Log, processFile, m.Name, flags.LogDir)
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

	return db.Db().Save(&m).Error
}
