package exec

import (
	"pmon3/model"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
)

func InsertAsQueued(flags *model.ExecFlags) (*model.Process, error) {

	logPath, err := process.GetLogPath(flags.LogDir, flags.Log, flags.Name)
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

	p := model.FromExecFlags(flags, logPath, user, groups)

	err = repo.ProcessOf(p).Save()

	return p, err
}
