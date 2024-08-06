package controller

import (
	"fmt"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"slices"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {

	var (
		nonDeptAppNames string
		deptAppNames    string
	)

	if cmd.GetArg1() == "apps-config-only" {
		nonDependentApps, dependentApps, err := conf.ComputeDepGraph(pmond.Config.AppsConfig.Apps)
		if err != nil {
			return ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		nonDeptAppNames = strings.Join(conf.AppNames(nonDependentApps), "\n")
		deptAppNames = strings.Join(conf.AppNames(dependentApps), "\n")
	} else {
		var all []model.Process
		err := db.Db().Find(&all).Error
		if err != nil {
			return ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		var qPs []model.Process
		qNm := map[string]bool{}

		dbProcessNames := model.ProcessNames(&all)

		for _, appConfigApp := range pmond.Config.AppsConfig.Apps {

			processName := appConfigApp.Flags.Name

			if slices.Contains(dbProcessNames, processName) {
				p := model.FromFileAndExecFlags(appConfigApp.File, &appConfigApp.Flags, "", nil)
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

		nonDependentAppsDb, dependentAppsDb, err := model.ComputeDepGraph(&qPs)
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		}

		nonDeptAppNames = strings.Join(model.ProcessNames(nonDependentAppsDb), "\n")
		deptAppNames = strings.Join(model.ProcessNames(dependentAppsDb), "\n")
	}

	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", nonDeptAppNames, deptAppNames),
	}
	return &newCmdResp
}
