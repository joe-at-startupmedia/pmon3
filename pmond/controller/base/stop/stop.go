package stop

import (
	"fmt"
	"os"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
	"time"
)

func ByProcess(p *model.Process, forced bool) error {
	if _, err := process.BeginPendingTask(p, "stop"); err != nil {
		return err
	}
	defer process.FinishPendingTask(p)
	return ByProcessWithoutPendingCheck(p, forced)
}

func ByProcessWithoutPendingCheck(p *model.Process, forced bool) error {
	status := model.StatusStopped
	// check process is running
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))

	p.ResetRestartCount()
	//if process is not currently running
	if os.IsNotExist(err) {
		if p.Status != status {
			if err := repo.ProcessOf(p).UpdateStatus(status); err != nil {
				return fmt.Errorf("stop process error: %w", err)
			}
			//we need to wait for the process to save before killing it to avoid a restart race condition
			time.Sleep(200 * time.Millisecond)
		}
	} else {
		//save it as stopped before we kill it to avoid a race condition
		if err = repo.ProcessOf(p).UpdateStatus(status); err != nil {
			return err
		}

		if err = process.SendOsKillSignal(p, forced); err != nil {
			return err
		}

		pmond.Log.Infof("stop process %s success", p.Stringify())
	}

	return nil
}
