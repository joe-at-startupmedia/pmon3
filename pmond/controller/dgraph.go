package controller

import (
	"fmt"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {

	var (
		nonDeptAppNames string
		deptAppNames    string
	)

	if cmd.GetArg1() == "apps-config-only" {
		nonDependentApps, dependentApps, err := conf.ComputeDepGraph(&pmond.Config.AppsConfig.Apps)
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		nonDeptAppNames = strings.Join(conf.AppNames(nonDependentApps), "\n")
		deptAppNames = strings.Join(conf.AppNames(dependentApps), "\n")
	} else {
		all, err := repo.Process().FindAll()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		var qPs []model.Process
		qNm := map[string]bool{}

		for _, appConfigApp := range pmond.Config.AppsConfig.Apps {
			processName := appConfigApp.Flags.Name
			p := model.FromFileAndExecFlags(appConfigApp.File, &appConfigApp.Flags, "", nil)
			qPs = append(qPs, *p)
			qNm[processName] = true
		}

		for _, dbPs := range all {
			processName := dbPs.Name
			if !qNm[processName] {
				qPs = append(qPs, dbPs)
				pmond.Log.Infof("append reamainder from db: pushing to stack %s", processName)
			} else {
				pmond.Log.Infof("overwritten with apps conf %s", processName)
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
