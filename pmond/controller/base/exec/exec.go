package exec

import (
	model2 "pmon3/model"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
)

func InsertAsQueued(flags *model2.ExecFlags) (*model2.Process, error) {

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

	p := model2.FromExecFlags(flags, logPath, user, groups)

	err = repo.ProcessOf(p).Save()

	return p, err
}
