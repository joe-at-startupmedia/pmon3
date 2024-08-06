package controller

import (
	"os/user"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"slices"
	"strings"
	"time"
)

func Initialize(cmd *protos.Cmd) *protos.CmdResp {

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}

	if cmd.GetArg1() == "apps-config-only" {
		startedFromConfig := StartsAppsFromConfig()
		if !startedFromConfig {
			newCmdResp.Error = "no applications were started"
		}
		return &newCmdResp
	} else {
		if err := StartApps(); err != nil {
			newCmdResp.Error = err.Error()
		}

		return &newCmdResp
	}
}

func StartsAppsFromConfig() bool {

	if pmond.Config.AppsConfig == nil || len(pmond.Config.AppsConfig.Apps) == 0 {
		return false
	}

	nonDependentApps, dependentApps, err := conf.ComputeDepGraph(pmond.Config.AppsConfig.Apps)
	if err != nil {
		return false
	}

	for _, app := range dependentApps {
		err := EnqueueProcess(app.File, &app.Flags)
		time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
		}
	}

	for _, app := range nonDependentApps {
		err := EnqueueProcess(app.File, &app.Flags)
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
		}
	}

	return true
}

func StartApps() error {
	var all []model.Process
	err := db.Db().Find(&all).Error
	if err != nil {
		return err
	}

	var qPs []model.Process
	qNm := map[string]bool{}

	dbProcessNames := model.ProcessNames(&all)

	for _, appConfigApp := range pmond.Config.AppsConfig.Apps {

		processName := appConfigApp.Flags.Name

		if slices.Contains(dbProcessNames, processName) {
			appLog, _ := getAppsConfigAppLogPath(&appConfigApp)
			appUser, _ := getAppsConfigAppUser(&appConfigApp)
			p := model.FromFileAndExecFlags(appConfigApp.File, &appConfigApp.Flags, appLog, appUser)
			qPs = append(qPs, *p)
			pmond.Log.Infof("overwrite with conf: pushing to stack %s", processName)
			qNm[processName] = true
		} else {
			dbPs, _ := model.GetProcessByName(processName, &all)
			if dbPs != nil {
				qPs = append(qPs, *dbPs)
				qNm[processName] = true
				pmond.Log.Infof("append from db: pushing to stack %s", processName)
			}
		}
	}

	for _, dbPs := range all {
		if !qNm[dbPs.Name] {
			qPs = append(qPs, dbPs)
			pmond.Log.Infof("append reamainder from db: pushing to stack %s", dbPs.Name)
		}
	}

	nonDependentApps, dependentApps, err := model.ComputeDepGraph(&qPs)
	if err != nil {
		pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		return err
	}

	if dependentApps != nil {
		pmond.Log.Infof("launch dependent %s", strings.Join(model.ProcessNames(dependentApps), " "))

		for _, app := range *dependentApps {
			pmond.Log.Infof("enqueue dependent and wait %s %d", app.Name, pmond.Config.GetDependentProcessEnqueuedWait())
			err := process.Enqueue(&app, true)
			time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
			if err != nil {
				pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				return err
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
				return err
			}
		}
	}

	return nil
}

func getAppsConfigAppLogPath(app *conf.AppsConfigApp) (string, error) {
	logPath, err := process.GetLogPath(app.Flags.LogDir, app.Flags.Log, app.File, app.Flags.Name)
	if err != nil {
		return "", err
	}
	return logPath, nil
}

func getAppsConfigAppUser(app *conf.AppsConfigApp) (*user.User, error) {
	u, _, err := process.SetUser(app.Flags.User)
	if err != nil {
		return nil, err
	}
	return u, nil
}
