package group

import (
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Drop(idOrName string, forceKill bool) *protos.CmdResp {

	var sent *protos.Cmd
	if forceKill {
		sent = base.SendCmdArg2("group_drop", idOrName, "force")
	} else {
		sent = base.SendCmd("group_drop", idOrName)
	}
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		all := newCmdResp.GetProcessList().GetProcesses()
		var allProcess [][]string
		for _, p := range all {
			process := model.ProcessFromProtobuf(p)
			allProcess = append(allProcess, process.RenderTable())
		}
		table_list.Render(allProcess)
	}
	return newCmdResp
}
