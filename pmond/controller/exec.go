package controller

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/exec"
	"pmon3/pmond/controller/base/restart"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"strings"
)

func getExecFileAbsPath(execFile string) (string, error) {
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist: %w", execFile, err)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("get file path error: %w", err)
	}

	return absPath, nil
}

func Exec(cmd *protos.Cmd) *protos.CmdResp {
	execFile := cmd.GetArg1()
	flags := cmd.GetArg2()
	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not parse flags: %w, flags: %s", err, flags))
	}
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	err = EnqueueProcess(execFile, parsedFlags)
	if err != nil {
		if strings.HasPrefix(err.Error(), "command error:") {
			return base.ErroredCmdResp(cmd, err)
		} else {
			newCmdResp.Error = err.Error()
		}
	}
	return &newCmdResp
}

func EnqueueProcess(execFile string, flags *model.ExecFlags) error {
	// get exec abs file path
	execPath, err := getExecFileAbsPath(execFile)
	if err != nil {
		return fmt.Errorf("command error: file argument error: %w", err)
	}
	name := flags.Name
	// get process file name
	if len(name) <= 0 {
		//set the filename as the name
		name = filepath.Base(execFile)
		flags.Name = name
	}
	p, err := repo.Process().FindByFileAndName(execPath, name)
	//if process exists
	if err == nil {
		pmond.Log.Debugf("updating as queued with flags: %v", flags)
		err = restart.UpdateAsQueued(p, execPath, flags)
	} else {
		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		_, err = exec.InsertAsQueued(execPath, flags)
	}
	return err
}
