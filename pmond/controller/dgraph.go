package controller

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/controller/base"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {

	var (
		nonDeptProcessNames string
		deptProcessNames    string
	)

	if cmd.GetArg1() == "process-config-only" {
		nonDependentProcesses, dependentProcesses, err := pmond.Config.ProcessConfig.ComputeDepGraph()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		nonDeptProcessNames = strings.Join(model.ExecFlagsNames(nonDependentProcesses), "\n")
		deptProcessNames = strings.Join(model.ExecFlagsNames(dependentProcesses), "\n")
	} else {
		all, err := repo.Process().FindAll()
		if err != nil {
			return base.ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
		}

		var qPs []model.Process
		qNm := map[string]bool{}

		for _, execFlags := range pmond.Config.ProcessConfig.Processes {
			processName := execFlags.Name
			p := model.FromExecFlags(&execFlags, "", nil, nil)
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

		nonDependentProcessesDb, dependentProcessesDb, err := model.ComputeDepGraph(&qPs)
		if err != nil {
			pmond.Log.Errorf("encountered error attempting to prioritize databse processes from dep graph: %s", err)
		}

		nonDeptProcessNames = strings.Join(model.ProcessNames(nonDependentProcessesDb), "\n")
		deptProcessNames = strings.Join(model.ProcessNames(dependentProcessesDb), "\n")
	}

	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", nonDeptProcessNames, deptProcessNames),
	}
	return &newCmdResp
}
