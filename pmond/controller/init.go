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

	if cmd.GetArg1() == "apps-config-only" {
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

	if pmond.Config.AppsConfig == nil || len(pmond.Config.AppsConfig.Apps) == 0 {
		return nil
	}

	nonDependentApps, dependentApps, err := pmond.Config.AppsConfig.ComputeDepGraph()
	if err != nil {
		return err
	}

	if blocking {
		err = appConfAppEnqueueUsingDepGraphResults(nonDependentApps, dependentApps)
	} else {
		go appConfAppEnqueueUsingDepGraphResults(nonDependentApps, dependentApps)
	}

	return err
}

func StartAppsFromBoth(blocking bool) error {
	nonDependentApps, dependentApps, err := getQueueableFromBoth()
	if err != nil {
		return err
	}

	if blocking {
		err = processEnqueueUsingDepGraphResults(nonDependentApps, dependentApps)
	} else {
		go processEnqueueUsingDepGraphResults(nonDependentApps, dependentApps)
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

	for _, appConfigApp := range pmond.Config.AppsConfig.Apps {
		processName := appConfigApp.Flags.Name
		appLog, _ := getAppsConfigAppLogPath(&appConfigApp)
		appUser, _ := getAppsConfigAppUser(&appConfigApp)
		groupFlags := appConfigApp.Flags.Groups
		groups, _ := repo.Group().FindOrInsertByNames(groupFlags)
		p := model.FromFileAndExecFlags(appConfigApp.File, &appConfigApp.Flags, appLog, appUser, groups)
		qPs = append(qPs, *p)
		qNm[processName] = true
	}

	for _, dbPs := range all {
		processName := dbPs.Name
		if !qNm[processName] {
			qPs = append(qPs, dbPs)
			pmond.Log.Infof("append reamainder from db: pushing to stack %s", processName)
		} else {
			pmond.Log.Infof("overwritten with apps conf: %s", processName)
		}
	}

	nonDependentApps, dependentApps, err := model.ComputeDepGraph(&qPs)
	if err != nil {
		pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		return nil, nil, err
	}

	return nonDependentApps, dependentApps, nil
}

func appConfAppEnqueueUsingDepGraphResults(nonDependentApps *[]model.AppsConfigApp, dependentApps *[]model.AppsConfigApp) error {

	var retErr error

	if dependentApps != nil {
		for _, app := range *dependentApps {
			pmond.Log.Infof("launch dependent %s", strings.Join(model.AppsConfigAppNames(dependentApps), " "))
			err := EnqueueProcess(app.File, &app.Flags)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	if nonDependentApps != nil {
		pmond.Log.Infof("launch independent %s", strings.Join(model.AppsConfigAppNames(nonDependentApps), " "))

		for _, app := range *nonDependentApps {
			err := EnqueueProcess(app.File, &app.Flags)
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	return retErr
}

func processEnqueueUsingDepGraphResults(nonDependentApps *[]model.Process, dependentApps *[]model.Process) error {

	var retErr error

	if dependentApps != nil {
		pmond.Log.Infof("launch dependent %s", strings.Join(model.ProcessNames(dependentApps), " "))

		for _, app := range *dependentApps {
			pmond.Log.Infof("enqueue dependent and wait %s %d", app.Name, pmond.Config.GetDependentProcessEnqueuedWait())
			err := process.Enqueue(&app, true)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	if nonDependentApps != nil {
		pmond.Log.Infof("launch independent %s", strings.Join(model.ProcessNames(nonDependentApps), " "))

		for _, app := range *nonDependentApps {
			pmond.Log.Infof("enqueue nondependent %s", app.Name)
			err := process.Enqueue(&app, true)
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				retErr = err
			}
		}
	}

	return retErr
}

func getAppsConfigAppLogPath(app *model.AppsConfigApp) (string, error) {
	logPath, err := process.GetLogPath(app.Flags.LogDir, app.Flags.Log, app.File, app.Flags.Name)
	if err != nil {
		return "", err
	}
	return logPath, nil
}

func getAppsConfigAppUser(app *model.AppsConfigApp) (*user.User, error) {
	u, _, err := process.SetUser(app.Flags.User)
	if err != nil {
		return nil, err
	}
	return u, nil
}
