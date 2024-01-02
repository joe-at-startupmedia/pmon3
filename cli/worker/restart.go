package worker

import (
	"errors"
	"pmon3/cli/service"
	"pmon3/pmond"
	"pmon3/pmond/executor"
	"pmon3/pmond/model"
	"time"
)

func Restart(pFile string, flags *model.ExecFlags) (string, error) {
	err, m := model.FindByProcessFileAndName(pmond.Db(), pFile, flags.Name)
	if err != nil {
		return "", errors.New("Could not find process")
	}

	cstLog := flags.Log
	if len(cstLog) > 0 && cstLog != m.Log {
		m.Log = cstLog
	}

	// if reset log dir
	if len(flags.LogDir) > 0 && len(flags.Log) == 0 {
		m.Log = ""
	}

	cstName := flags.Name
	if len(cstName) > 0 && cstName != m.Name {
		m.Name = cstName
	}

	extArgs := flags.Args
	if len(extArgs) > 0 {
		m.Args = extArgs
	}

	// get run process user
	runUser, err := GetProcUser(flags)
	if err != nil {
		return "", err
	}

	process, err := executor.Exec(m.ProcessFile, m.Log, m.Name, m.Args, runUser, !flags.NoAutoRestart, flags.LogDir)
	if err != nil {
		return "", err
	}

	// update process extra data
	process.ID = m.ID
	process.CreatedAt = m.CreatedAt
	process.UpdatedAt = time.Now()

	waitData := service.NewProcStat(process).Wait()

	return service.AddData(waitData)
}
