package stop

import (
	"pmon3/cli/cmd/base"
	table_list "pmon3/cli/output/process/list"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
)

func Stop(idOrName string) *protos.CmdResp {

	sent := base.SendCmd("group_stop", idOrName)
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
