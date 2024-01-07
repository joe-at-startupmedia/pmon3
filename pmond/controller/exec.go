package controller

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"pmon3/pmond/utils/crypto"
	"strings"
)

func getExecFileAbsPath(execFile string) (string, error) {
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist", execFile)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("get file path error: %v", err.Error())
	}

	return absPath, nil
}

func Exec(cmd *protos.Cmd) *protos.CmdResp {
	execFile := cmd.GetArg1()
	flags := cmd.GetArg2()
	// get exec abs file path
	execPath, err := getExecFileAbsPath(execFile)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("file arguments error: %+v", err))
	}
	execflags := model.ExecFlags{}
	parsedFlags, err := execflags.Parse(flags)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Sprintf("could not parse flags: %+v", err))
	}
	name := parsedFlags.Name
	// get process file name
	if len(name) <= 0 {
		//set the filename as the name
		name = filepath.Base(execFile)
		parsedFlags.Name = name
	}
	err, p := model.FindProcessByFileAndName(pmond.Db(), execPath, name)
	//if process exists
	if err == nil {
		pmond.Log.Debugf("updating as queued with flags: %v", flags)
		err = UpdateAsQueued(p, execPath, parsedFlags)
	} else {
		pmond.Log.Debugf("inserting as queued with flags: %v", flags)
		err = insertAsQueued(execPath, parsedFlags)
	}
	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}
	if err != nil {
		newCmdResp.Error = err.Error()
	}
	return &newCmdResp
}

func insertAsQueued(processFile string, flags *model.ExecFlags) error {

	var processParams = []string{flags.Name}
	if len(flags.Args) > 0 {
		processParams = append(processParams, strings.Split(flags.Args, " ")...)
	}

	logPath, err := process.GetLogPath(flags.Log, crypto.Crc32Hash(processFile+flags.Name), flags.LogDir)
	if err != nil {
		return err
	}

	user, err := process.SetUser(flags.User)
	if err != nil {
		return err
	}

	p := model.Process{
		Pid:         0,
		Log:         logPath,
		Name:        flags.Name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     nil,
		Status:      model.StatusQueued,
		Uid:         conv.StrToUint32(user.Uid),
		Gid:         conv.StrToUint32(user.Gid),
		Username:    user.Username,
		AutoRestart: !flags.NoAutoRestart,
	}

	err = pmond.Db().Save(&p).Error

	return err
}
