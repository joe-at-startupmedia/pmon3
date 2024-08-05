package controller

import (
	"fmt"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/protos"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {
	queueOrder, dGraph, err := conf.ComputeDepGraph(pmond.Config.AppsConfig.Apps)
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
	}

	output1 := strings.Join(conf.MapKeys(queueOrder), "\n")
	output2 := strings.Join(conf.MapKeys(dGraph), "\n")
	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", output1, output2),
	}
	return &newCmdResp
}
