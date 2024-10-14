package stop

import (
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/repo"
	"pmon3/pmond/shell"
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
	isRunning := shell.ExecIsRunning(p)
	p.ResetRestartCount()
	//if process is not currently running
	if !isRunning {
		if p.Status != status {
			if err := repo.ProcessOf(p).UpdateStatus(status); err != nil {
				return fmt.Errorf("stop process error: %w", err)
			}
			//we need to wait for the process to save before killing it to avoid a restart race condition
			time.Sleep(200 * time.Millisecond)
		}
	} else {
		//save it as stopped before we kill it to avoid a race condition
		if err := repo.ProcessOf(p).UpdateStatus(status); err != nil {
			return err
		}

		if err := process.SendOsKillSignal(p, forced); err != nil {
			return err
		}

		pmond.Log.Infof("stop process %s success", p.Stringify())
	}

	return nil
}
