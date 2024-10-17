package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Stop(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()
	forced := cmd.GetArg2() == "force"
	return StopByParams(cmd, idOrName, forced)
}

func StopByParams(cmd *protos2.Cmd, idOrName string, forced bool) *protos2.CmdResp {
	g, err := repo.Group().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	for i := range g.Processes {
		p := g.Processes[i]
		// check process is running
		err = stop.ByProcess(p, forced)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}
	}
	return listProcesses(cmd, g)
}
