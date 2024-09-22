package exec

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
)

func GetExecFileAbsPath(execFile string) (string, error) {
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

func InsertAsQueued(processFile string, flags *model.ExecFlags) (*model.Process, error) {

	logPath, err := process.GetLogPath(flags.LogDir, flags.Log, processFile, flags.Name)
	if err != nil {
		return nil, err
	}

	user, _, err := process.SetUser(flags.User)
	if err != nil {
		return nil, err
	}

	groups, err := repo.Group().FindOrInsertByNames(flags.Groups)
	if err != nil {
		return nil, err
	}

	p := model.FromFileAndExecFlags(processFile, flags, logPath, user, groups)

	err = repo.ProcessOf(p).Save()

	return p, err
}
