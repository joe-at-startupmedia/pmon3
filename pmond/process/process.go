package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/observer"
	"pmon3/pmond/repo"
	"pmon3/pmond/shell"
	"strconv"
	"sync"
	"time"
)

var pendingTask sync.Map

func BeginPendingTask(p *model.Process, action string) (bool, error) {
	existingAction, loaded := pendingTask.LoadOrStore(p.ID, action)
	if loaded {
		err := fmt.Errorf("process.BeginPendingTask, task is already pending for process with id (%d) with action(%s)", p.ID, existingAction)
		pmond.Log.Warn(err)
		return false, err
	}
	return true, nil
}

func FinishPendingTask(p *model.Process) {
	pendingTask.Delete(p.ID)
}

func Enqueue(p *model.Process, force bool) error {
	if (!shell.ExecIsRunning(p) && p.Status == model.StatusQueued) || force {
		if updatedFromPsCmd(p) {
			return nil
		}

		_, err := proxyWorker(p, "start")
		if err != nil {
			return err
		}
	}

	return nil
}

func Restart(p *model.Process, isInitializing bool) (bool, error) {

	_, loaded := pendingTask.Load(p.ID)
	if loaded {
		pmond.Log.Warnf("process.Restart, task is already pending for process with id(%d)", p.ID)
		return false, nil
	}

	restarted := false
	if !shell.ExecIsRunning(p) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed || p.Status == model.StatusClosed) {
		if updatedFromPsCmd(p) {
			return false, nil
		}

		if !p.AutoRestart {
			if p.Status == model.StatusRunning { // but process is dead, update db state
				observer.HandleEvent(&observer.Event{
					Type:    observer.FailedEvent,
					Process: p,
				})
				repo.ProcessOf(p).UpdateStatus(model.StatusFailed)
			}
			return false, nil
		}

		if isInitializing {
			pmond.Log.Infof("(re)starting process during initialization(%t): %s", isInitializing, p.Stringify())
		} else {
			observer.HandleEvent(&observer.Event{
				Type:    observer.RestartEvent,
				Process: p,
			})
		}

		restarted = true

		_, err := proxyWorker(p, "restart")
		if err != nil {
			return restarted, err
		} else {
			if p.Status != model.StatusClosed {
				p.IncrRestartCount()
			}
		}
	}

	return restarted, nil
}

func SendOsKillSignal(p *model.Process, forced bool) error {

	if !shell.ExecIsRunning(p) {
		pmond.Log.Warnf("Cannot kill process (%s - %s) that isnt running", p.Stringify(), p.GetPidStr())
		return nil
	}

	if forced {
		return shell.ExecKillProcessForcefully(p)
	} else {
		return shell.ExecKillProcess(p)
	}
}

func SetUser(runUser string) (*user.User, []string, error) {
	var curUser *user.User
	var err error

	if len(runUser) <= 0 {
		curUser, err = user.LookupId(strconv.Itoa(os.Getuid()))
	} else {
		curUser, err = user.Lookup(runUser)
	}

	if err != nil {
		return nil, nil, err
	}

	groupIds, err := curUser.GroupIds()

	if err != nil {
		return nil, nil, err
	}

	return curUser, groupIds, nil
}

func findPidFromPsCmd(p *model.Process) uint32 {
	if len(p.Args) > 0 {
		return shell.ExecFindPidFromProcessNameAndArgs(p)
	} else {
		return shell.ExecFindPidFromProcessName(p)
	}
}

func findPpidFromPsCmd(p *model.Process) uint32 {
	if len(p.Args) > 0 {
		return shell.ExecFindPpidFromProcessNameAndArgs(p)
	} else {
		return shell.ExecFindPpidFromProcessName(p)
	}
}

// used as a last alternative if /proc/[pid]/status and golang isNotExist fail to detect running
func updatedFromPsCmd(p *model.Process) bool {

	newPid := findPidFromPsCmd(p)

	if newPid != 0 && newPid != p.Pid {
		newPpid := findPpidFromPsCmd(p)
		if newPpid == 1 {
			pmond.Log.Errorf("Detected orphan process with the same process name: %s pid: %d", p.Name, newPid)
			return true
		}
		p.Pid = newPid
		p.Status = model.StatusRunning
		err := repo.ProcessOf(p).Save()
		return err == nil
	}

	return false
}

func proxyWorker(m *model.Process, cmd string) ([]string, error) {

	var (
		tb     []string
		output string
		err    error
	)

	switch cmd {
	case "start":
		pmond.Log.Infof("starting process: %s", m.Stringify())
		output, err = workerStart(m)
	case "restart":
		output, err = workerRestart(m)
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(output), &tb)
	return tb, nil
}

func workerRestart(p *model.Process) (string, error) {
	//returns an instance of the process model
	execP, err := exec(p)
	if err != nil {
		return "", err
	}

	execP.ID = p.ID
	execP.CreatedAt = p.CreatedAt
	execP.UpdatedAt = time.Now()

	process := NewProcStat(execP).Wait()
	return repo.ProcessOf(process).FindAndSave()
}

func workerStart(p *model.Process) (string, error) {
	//returns an instance of the process model
	execP, err := exec(p)
	if err != nil {
		return "", err
	}

	execP.CreatedAt = time.Now()
	execP.UpdatedAt = time.Now()

	process := NewProcStat(execP).Wait()
	return repo.ProcessOf(process).FindAndSave()
}
