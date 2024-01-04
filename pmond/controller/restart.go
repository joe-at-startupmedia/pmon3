package controller

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/executor"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"pmon3/pmond/utils/crypto"
	"strings"
)

func Restart(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	flags := cmd.GetArg2()
	return RestartByParams(cmd, idOrName, flags)
}

func RestartByParams(cmd *protos.Cmd, idOrName string, flags string) *protos.CmdResp {
	err, p := model.FindProcessByIdOrName(pmond.Db(), idOrName)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("could not find process: %+v", err))
	}
	if process.IsRunning(p.Pid) {
		if err := process.TryStop(p, model.StatusStopped, false); err != nil {
			return ErroredCmdResp(cmd, fmt.Sprintf("restart error: %s", err.Error()))
		}
	}
	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("could not parse flags: %+v", err))
	} else {
		pmond.Log.Debugf("restart process: %v", flags)
		err = ExecRestart(p, p.ProcessFile, parsedFlags)
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
		}
		return &newCmdResp
	}
}

func ExecRestart(m *model.Process, processFile string, flags *model.ExecFlags) error {
	// only stopped and failed process can be restarted
	if m.Status != model.StatusStopped && m.Status != model.StatusFailed {
		return fmt.Errorf("process already running with the name provided: %s", m.Name)
	}
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
		user, err := process.GetProcUser(flags)
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

	return pmond.Db().Save(&m).Error
}
