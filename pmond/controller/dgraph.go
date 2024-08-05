package controller

import (
	"fmt"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {
	nonDeptApps, deptApps, err := conf.ComputeDepGraph(pmond.Config.AppsConfig.Apps)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
	}

	nonDeptAppNames := strings.Join(conf.MapKeys(nonDeptApps), "\n")
	deptAppNames := strings.Join(conf.MapKeys(deptApps), "\n")

	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", nonDeptAppNames, deptAppNames),
	}
	return &newCmdResp
}
