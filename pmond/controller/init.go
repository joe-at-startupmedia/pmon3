package controller

import (
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"strings"
	"time"
)

func Initialize(cmd *protos.Cmd) *protos.CmdResp {

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}

	blocking := cmd.GetArg2() == "blocking"

	var err error

	if cmd.GetArg1() == "process-config-only" {
		err = StartsAppsFromConfig(blocking)
	} else {
		err = StartAppsFromBoth(blocking)
	}
	if err != nil {
		newCmdResp.Error = err.Error()
	}

	return &newCmdResp
}

func StartsAppsFromConfig(blocking bool) error {

	if pmond.Config.ProcessConfig == nil || len(pmond.Config.ProcessConfig.Processes) == 0 {
		return nil
	}

	nonDependentProcesses, dependentProcesses, err := pmond.Config.ProcessConfig.ComputeDepGraph()
	if err != nil {
		return err
	}

	if blocking {
		err = execFlagsEnqueueUsingDepGraphResults(nonDependentProcesses, dependentProcesses)
	} else {
		go execFlagsEnqueueUsingDepGraphResults(nonDependentProcesses, dependentProcesses)
	}

	return err
}

func StartAppsFromBoth(blocking bool) error {
	nonDependentProcesses, dependentProcesses, err := getQueueableFromBoth()
	if err != nil {
		return err
	}

	if blocking {
		err = processEnqueueUsingDepGraphResults(nonDependentProcesses, dependentProcesses)
	} else {
		go processEnqueueUsingDepGraphResults(nonDependentProcesses, dependentProcesses)
	}

	return err
}

func getQueueableFromBoth() (*[]model.Process, *[]model.Process, error) {
	all, err := repo.Process().FindAll()
	if err != nil {
		return nil, nil, err
	}

	var qPs []model.Process
	qNm := map[string]bool{}

	for _, execFlags := range pmond.Config.ProcessConfig.Processes {
		processName := execFlags.Name
		pLog, _ := getExecFlagsLogPath(&execFlags)
		pUser, _ := getExecFlagsUser(&execFlags)
		groupFlags := execFlags.Groups
		groups, _ := repo.Group().FindOrInsertByNames(groupFlags)
		p := model.FromExecFlags(&execFlags, pLog, pUser, groups)
		qPs = append(qPs, *p)
		qNm[processName] = true
	}

	for _, dbPs := range all {
		processName := dbPs.Name
		if !qNm[processName] {
			qPs = append(qPs, dbPs)
			pmond.Log.Infof("append reamainder from db: pushing to stack %s", processName)
		} else {
			pmond.Log.Infof("overwritten with process config: %s", processName)
		}
	}

	nonDependentProcesses, dependentProcesses, err := model.ComputeDepGraph(&qPs)
	if err != nil {
		pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		return nil, nil, err
	}

	return nonDependentProcesses, dependentProcesses, nil
}

func execFlagsEnqueueUsingDepGraphResults(nonDependentProcesses *[]model.ExecFlags, dependentProcesses *[]model.ExecFlags) error {

	var retErr error

	if dependentProcesses != nil {
		for _, execFlags := range *dependentProcesses {
			pmond.Log.Infof("launch dependent %s", strings.Join(model.ExecFlagsNames(dependentProcesses), " "))
			err := EnqueueProcess(&execFlags)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	if nonDependentProcesses != nil {
		pmond.Log.Infof("launch independent %s", strings.Join(model.ExecFlagsNames(nonDependentProcesses), " "))

		for _, execFlags := range *nonDependentProcesses {
			err := EnqueueProcess(&execFlags)
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	return retErr
}

func processEnqueueUsingDepGraphResults(nonDependentProcesses *[]model.Process, dependentProcesses *[]model.Process) error {

	var retErr error

	if dependentProcesses != nil {
		pmond.Log.Infof("launch dependent %s", strings.Join(model.ProcessNames(dependentProcesses), " "))

		for _, dp := range *dependentProcesses {
			pmond.Log.Infof("enqueue dependent and wait %s %d", dp.Name, pmond.Config.GetDependentProcessEnqueuedWait())
			err := process.Enqueue(&dp, true)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	if nonDependentProcesses != nil {
		pmond.Log.Infof("launch independent %s", strings.Join(model.ProcessNames(nonDependentProcesses), " "))

		for _, ndp := range *nonDependentProcesses {
			pmond.Log.Infof("enqueue nondependent %s", ndp.Name)
			err := process.Enqueue(&ndp, true)
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	return retErr
}

func getExecFlagsLogPath(execFlags *model.ExecFlags) (string, error) {
	logPath, err := process.GetLogPath(execFlags.LogDir, execFlags.Log, execFlags.Name)
	if err != nil {
		return "", err
	}
	return logPath, nil
}

func getExecFlagsUser(execFlags *model.ExecFlags) (*user.User, error) {
	u, _, err := process.SetUser(execFlags.User)
	if err != nil {
		return nil, err
	}
	return u, nil
}
