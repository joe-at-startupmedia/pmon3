package controller

import (
	"os/user"
	model2 "pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
	"strings"
	"time"
)

func Initialize(cmd *protos2.Cmd) *protos2.CmdResp {

	newCmdResp := protos2.CmdResp{
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

func getQueueableFromBoth() (*[]model2.Process, *[]model2.Process, error) {
	all, err := repo.Process().FindAll()
	if err != nil {
		return nil, nil, err
	}

	var qPs []model2.Process
	qNm := map[string]bool{}

	for _, execFlags := range pmond.Config.ProcessConfig.Processes {
		processName := execFlags.Name
		pLog, _ := getExecFlagsLogPath(&execFlags)
		pUser, _ := getExecFlagsUser(&execFlags)
		groupFlags := execFlags.Groups
		groups, _ := repo.Group().FindOrInsertByNames(groupFlags)
		p := model2.FromExecFlags(&execFlags, pLog, pUser, groups)
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

	nonDependentProcesses, dependentProcesses, err := model2.ComputeDepGraph(&qPs)
	if err != nil {
		pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		return nil, nil, err
	}

	return nonDependentProcesses, dependentProcesses, nil
}

func execFlagsEnqueueUsingDepGraphResults(nonDependentProcesses *[]model2.ExecFlags, dependentProcesses *[]model2.ExecFlags) error {

	var retErr error

	if dependentProcesses != nil {
		for _, execFlags := range *dependentProcesses {
			pmond.Log.Infof("launch dependent %s", strings.Join(model2.ExecFlagsNames(dependentProcesses), " "))
			err := EnqueueProcess(&execFlags)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	if nonDependentProcesses != nil {
		pmond.Log.Infof("launch independent %s", strings.Join(model2.ExecFlagsNames(nonDependentProcesses), " "))

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

func processEnqueueUsingDepGraphResults(nonDependentProcesses *[]model2.Process, dependentProcesses *[]model2.Process) error {

	var retErr error

	if dependentProcesses != nil {
		pmond.Log.Infof("launch dependent %s", strings.Join(model2.ProcessNames(dependentProcesses), " "))

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
		pmond.Log.Infof("launch independent %s", strings.Join(model2.ProcessNames(nonDependentProcesses), " "))

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

func getExecFlagsLogPath(execFlags *model2.ExecFlags) (string, error) {
	logPath, err := process.GetLogPath(execFlags.LogDir, execFlags.Log, execFlags.Name)
	if err != nil {
		return "", err
	}
	return logPath, nil
}

func getExecFlagsUser(execFlags *model2.ExecFlags) (*user.User, error) {
	u, _, err := process.SetUser(execFlags.User)
	if err != nil {
		return nil, err
	}
	return u, nil
}
