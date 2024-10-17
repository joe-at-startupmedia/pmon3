package group

import (
	"pmon3/pmond/controller/base"
	"pmon3/pmond/controller/base/restart"
	"pmon3/pmond/repo"
	protos2 "pmon3/protos"
)

func Restart(cmd *protos2.Cmd) *protos2.CmdResp {
	idOrName := cmd.GetArg1()
	flags := cmd.GetArg2()
	return RestartByParams(cmd, idOrName, flags, true)
}

func RestartByParams(cmd *protos2.Cmd, idOrName string, flags string, incrementCounter bool) *protos2.CmdResp {
	g, err := repo.Group().FindByIdOrName(idOrName)
	if err != nil {
		return base.ErroredCmdResp(cmd, err)
	}

	for i := range g.Processes {
		p := g.Processes[i]
		_, err = restart.ByProcess(cmd, p, p.GetIdStr(), flags, incrementCounter)
		if err != nil {
			return base.ErroredCmdResp(cmd, err)
		}
	}
	return listProcesses(cmd, g)
}
