package desc

import (
	"pmon3/cli/cmd/base"
	"pmon3/cli/output/process/desc"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"pmon3/utils/conv"
	"strconv"
)

func Desc(idOrName string) *protos.CmdResp {

	sent := base.SendCmd("desc", idOrName)
	newCmdResp := base.GetResponse(sent)
	process := newCmdResp.GetProcess()
	if process != nil {
		rel := [][]string{
			{"status", process.Status},
			{"id", conv.Uint32ToStr(process.Id)},
			{"name", process.Name},
			{"pid", conv.Uint32ToStr(process.Pid)},
			{"restarted", conv.Uint32ToStr(process.RestartCount)},
			{"process", process.ProcessFile},
			{"args", process.Args},
			{"env-vars", process.EnvVars},
			{"user", process.Username},
			{"log", process.Log},
			{"no-autorestart", strconv.FormatBool(!process.AutoRestart)},
			{"dependencies", process.Dependencies},
			{"groups", model.GetGroupString(process)},
			{"created_at", process.CreatedAt},
			{"updated_at", process.UpdatedAt},
		}
		table_desc.Render(rel)
	}
	return newCmdResp
}
