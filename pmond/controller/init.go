package controller

import (
	"fmt"
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
			hasError = false
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

	nonDeptApps, deptApps, err := conf.ComputeDepGraph(pmond.Config.AppsConfig.Apps)
	if err != nil {
		return false
	}

	for _, app := range deptApps {
		err := EnqueueProcess(app.File, &app.Flags)
		time.Sleep(pmond.Config.GetDependentProcessEnqueuedWait())
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
		}
	}

	for _, app := range nonDeptApps {
		err := EnqueueProcess(app.File, &app.Flags)
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to enqueue process: %s", err)
		}
	}

	return true
}
