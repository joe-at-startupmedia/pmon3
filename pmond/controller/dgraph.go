package controller

import (
	"fmt"
	model2 "pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
	"strings"
)

func Dgraph(cmd *protos2.Cmd) *protos2.CmdResp {

	var (
		nonDeptProcessNames string
		deptProcessNames    string
	)

	if cmd.GetArg1() == "process-config-only" {
		nonDependentProcesses, dependentProcesses, err := pmond.Config.ProcessConfig.ComputeDepGraph()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		nonDeptProcessNames = strings.Join(model2.ExecFlagsNames(nonDependentProcesses), "\n")
		deptProcessNames = strings.Join(model2.ExecFlagsNames(dependentProcesses), "\n")
	} else {
		all, err := repo.Process().FindAll()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		var qPs []model2.Process
		qNm := map[string]bool{}

		for _, execFlags := range pmond.Config.ProcessConfig.Processes {
			processName := execFlags.Name
			p := model2.FromExecFlags(&execFlags, "", nil, nil)
			qPs = append(qPs, *p)
			qNm[processName] = true
		}

		for _, dbPs := range all {
			processName := dbPs.Name
			if !qNm[processName] {
				qPs = append(qPs, dbPs)
				pmond.Log.Infof("append reamainder from db: pushing to stack %s", processName)
			} else {
				pmond.Log.Infof("overwritten with process config %s", processName)
			}
		}

		nonDependentProcessesDb, dependentProcessesDb, err := model2.ComputeDepGraph(&qPs)
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		}

		nonDeptProcessNames = strings.Join(model2.ProcessNames(nonDependentProcessesDb), "\n")
		deptProcessNames = strings.Join(model2.ProcessNames(dependentProcessesDb), "\n")
	}

	newCmdResp := protos2.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", nonDeptProcessNames, deptProcessNames),
	}
	return &newCmdResp
}
