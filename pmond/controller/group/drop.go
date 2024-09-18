package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/del"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
)

func Drop(cmd *protos.Cmd) *protos.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return DropByParams(cmd, idOrName, forced)
}

func DropByParams(cmd *protos.Cmd, idOrName string, forced bool) *protos.CmdResp {
	g, err := repo.Group().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	for i := range g.Processes {
		p := g.Processes[i]
		// check process is running
		err = del.ByProcess(p, forced)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}
	}
	return listProcesses(cmd, g)
}
