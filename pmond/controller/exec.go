package controller

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
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
		return ErroredCmdResp(cmd, fmt.Errorf("command error: could not parse flags: %w, flags: %s", err, flags))
	}
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	err = EnqueueProcess(execFile, parsedFlags)
	if err != nil {
		if strings.HasPrefix(err.Error(), "command error:") {
			return ErroredCmdResp(cmd, err)
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
	err, p := model.FindProcessByFileAndName(db.Db(), execPath, name)
	//if process exists
	if err == nil {
		pmond.Log.Debugf("updating as queued with flags: %v", flags)
		err = UpdateAsQueued(p, execPath, flags)
	} else {
		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		_, err = insertAsQueued(execPath, flags)
	}
	return err
}

func insertAsQueued(processFile string, flags *model.ExecFlags) (*model.Process, error) {

	logPath, err := process.GetLogPath(flags.LogDir, flags.Log, processFile, flags.Name)
	if err != nil {
		return nil, err
	}

	user, _, err := process.SetUser(flags.User)
	if err != nil {
		return nil, err
	}

	p := model.FromFileAndExecFlags(processFile, flags, logPath, user)

	err = db.Db().Save(&p).Error

	return p, err
}
