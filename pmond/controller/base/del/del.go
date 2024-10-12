package del

import (
	"errors"
	"os"
	"pmon3/pmond/controller/base/stop"
	"pmon3/pmond/flap_detector"
	"pmon3/pmond/model"
	"pmon3/pmond/repo"
)

func ByProcess(p *model.Process, forced bool) error {
	stopErr := stop.ByProcess(p, forced)
	delErr := repo.ProcessOf(p).Delete()
	flap_detector.Delete(p.ID)
	logErr := os.Remove(p.Log)
	return errors.Join(stopErr, delErr, logErr)
}
