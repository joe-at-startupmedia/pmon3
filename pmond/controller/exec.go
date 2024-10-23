package controller

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/exec"
	"pmon3/pmond/controller/base/restart"
	"pmon3/pmond/repo"
	"pmon3/protos"
)

func setExecFileAbsPath(execFlags *model.ExecFlags) error {

	if path.IsAbs(execFlags.File) {
		_, err := os.Stat(execFlags.File)
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist: %w", execFlags.File, err)
		}

		return nil
	} else {
		absPath, err := filepath.Abs(execFlags.File)
		_, err = os.Stat(absPath)
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist: %w", execFlags.File, err)
		}
		execFlags.File = absPath
	}

	return nil
}

func Exec(cmd *protos.Cmd) *protos.CmdResp {
	flags := cmd.GetArg1()
	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not parse flags: %w, flags: %s", err, flags))
	}
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	err = EnqueueProcess(parsedFlags)
	if err != nil {
		newCmdResp.Error = err.Error()
	}
	return &newCmdResp
}

func EnqueueProcess(flags *model.ExecFlags) error {
	// get exec abs file path
	err := setExecFileAbsPath(flags)
	if err != nil {
		return fmt.Errorf("command error: file argument error: %w", err)
	}
	name := flags.Name
	// get process file name
	if len(name) == 0 {
		//set the filename as the name
		name = filepath.Base(flags.File)
		flags.Name = name
	}
	p, err := repo.Process().FindByFileAndName(flags.File, name)
	//if process exists
	if err == nil {
		pmond.Log.Debugf("updating as queued with flags: %v", flags)
		err = restart.UpdateAsQueued(p, flags)
	} else {
		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		_, err = exec.InsertAsQueued(flags)
	}
	return err
}
