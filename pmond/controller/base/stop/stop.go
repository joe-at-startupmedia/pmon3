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

func ByProcess(p *model.Process, forced bool, status model.ProcessStatus) error {
	// check process is running
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	//if process is not currently running
	if os.IsNotExist(err) {
		if p.Status != status {
			if err := repo.ProcessOf(p).UpdateStatus(status); err != nil {
				return fmt.Errorf("stop process error: %w", err)
			}
		}
	}

	//we need to wait for the process to save before killing it to avoid a restart race condition
	time.Sleep(200 * time.Millisecond)
	// try to kill the process
	err = process.SendOsKillSignal(p, status, forced)
	if err != nil {
		return fmt.Errorf("stop process error: %w", err)
	}

	pmond.Log.Infof("stop process %s success", p.Stringify())

	p.ResetRestartCount()

	return nil
}
