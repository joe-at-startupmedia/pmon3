package del

import (
	"errors"
	"os"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/model"
	"pmon3/pmond/repo"
)

func ByProcess(p *model.Process, forced bool) error {
	stopErr := stop.ByProcess(p, forced, model.StatusStopped)
	delErr := repo.ProcessOf(p).Delete()
	logErr := os.Remove(p.Log)
	return errors.Join(stopErr, delErr, logErr)
}
