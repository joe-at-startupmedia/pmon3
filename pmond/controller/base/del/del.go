package del

import (
	"os"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/model"
	"pmon3/pmond/repo"
)

func ByProcess(p *model.Process, forced bool) error {
	err := stop.ByProcess(p, forced, model.StatusStopped)
	if err != nil {
		return err
	}
	err = repo.ProcessOf(p).Delete()
	_ = os.Remove(p.Log)
	return err
}
