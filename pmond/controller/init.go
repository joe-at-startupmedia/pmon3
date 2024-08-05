package controller

import (
	"fmt"
	"github.com/joe-at-startupmedia/depgraph"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"time"
)

func Initialize(cmd *protos.Cmd) *protos.CmdResp {

	newCmdResp := protos.CmdResp{
		Id:   cmd.GetId(),
		Name: cmd.GetName(),
	}

	startedFromConfig := StartsAppsFromConfig()
	if startedFromConfig {
		return &newCmdResp
	} else {
		var all []model.Process
		err := db.Db().Find(&all, "status = ?", model.StatusStopped).Error
		if err != nil {
			return ErroredCmdResp(cmd, fmt.Errorf("error finding stopped processes: %w", err))
		} else if len(all) == 0 {
			return ErroredCmdResp(cmd, fmt.Errorf("Could not find any stopped processes"))
		}

		var (
			cr       *protos.CmdResp
			hasError bool = false
		)

		for _, process := range all {
			pmond.Log.Debugf("restart: %s", process.GetIdStr())
			cr := RestartByParams(cmd, process.GetIdStr(), "{}", false)
			if len(cr.GetError()) > 0 {
				pmond.Log.Debugf("encountered error attempting to restart: %s", cr.GetError())
				hasError = true
				break
			}
		}

		if hasError {
			newCmdResp.Error = cr.GetError()
		}
		return &newCmdResp
	}
}

func StartsAppsFromConfig() bool {

	var all []model.Process
	err := db.Db().Find(&all).Error
	if err != nil || len(all) > 0 || pmond.Config.AppsConfig == nil {
		for _, process := range all {
			pmond.Log.Debugf("config start: %s", process.Stringify())
		}
		return false
	}
	apps := pmond.Config.AppsConfig.Apps

	if len(apps) > 0 {
		g := depgraph.New()
		depAppNames := make(map[string]conf.AppsConfigApp)
		nonDepAppNames := make(map[string]conf.AppsConfigApp)
		for _, app := range apps {
			if len(app.Flags.Dependencies) > 0 {
				depAppNames[app.Flags.Name] = app
				for _, dep := range app.Flags.Dependencies {
					err = g.DependOn(app.Flags.Name, dep)
					if err != nil {
						pmond.Log.Errorf("encountered error building app dependency tree: %s", err)
					}
				}
			} else {
				nonDepAppNames[app.Flags.Name] = app
			}
		}

		if len(g.Leaves()) > 0 {

			for i, appName := range g.TopoSorted() {
				if depAppNames[appName].File != "" {
					pmond.Log.Infof("%d: %s\n", i, appName)
					app := depAppNames[appName]
					err := EnqueueProcess(app.File, &app.Flags)
					time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
					if err != nil {
						pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
					}
				} else if nonDepAppNames[appName].File != "" {
					pmond.Log.Infof("%d: %s\n", i, appName)
					app := nonDepAppNames[appName]
					err := EnqueueProcess(app.File, &app.Flags)
					time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
					if err != nil {
						pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
					}
					nonDepAppNames[appName] = conf.AppsConfigApp{}
				} else {
					pmond.Log.Warnf("dependencies: %s is not a valid app name", appName)
				}
			}

			for appName, app := range nonDepAppNames {
				err := EnqueueProcess(app.File, &app.Flags)
				if err != nil {
					pmond.Log.Errorf("encountered error attempting to enqueue process %s: %s", appName, err)
				}
			}

		} else {
			for _, app := range apps {
				err := EnqueueProcess(app.File, &app.Flags)
				if err != nil {
					pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
				}
			}
		}

	}
	return true
}
