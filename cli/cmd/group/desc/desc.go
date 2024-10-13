package desc

import (
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/desc"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/utils/conv"
)

func Desc(idOrName string) *protos.CmdResp {
	sent := base.SendCmd("group_desc", idOrName)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) == 0 {
		group := newCmdResp.GetGroup()
		rel := [][]string{
			{"id", conv.Uint32ToStr(group.Id)},
			{"name", group.Name},
		}
		all := newCmdResp.GetProcessList().GetProcesses()
		var allProcess [][]string
		for _, p := range all {
			process := model.ProcessFromProtobuf(p)
			allProcess = append(allProcess, process.RenderTable())
		}
		table_desc.Render(rel)
		table_list.Render(allProcess)
	}
	return newCmdResp
}
