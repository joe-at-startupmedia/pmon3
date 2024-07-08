package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/observer"
	"pmon3/pmond/utils/conv"
	"strconv"
	"strings"
	"time"

	"github.com/goinbox/shell"
)

func IsRunning(pid uint32) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))

	//if it doesn't exist in proc/n/status ask the OS
	if err != nil {
		//it's running if it exists
		return !os.IsNotExist(err)
	}

	return true
}

func findPidFromPsCmd(p *model.Process) uint32 {
	var rel *shell.ShellResult
	if len(p.Args) > 0 {
		rel = shell.RunCmd(fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $2}'", p.Name, p.Args))
	} else {
		rel = shell.RunCmd(fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $2}'", p.Name))
	}
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		return newPid
	}

	return 0
}

func findPpidFromPsCmd(p *model.Process) uint32 {
	var rel *shell.ShellResult
	if len(p.Args) > 0 {
		rel = shell.RunCmd(fmt.Sprintf("ps -ef | grep ' %s %s$' | grep -v grep | awk '{print $3}'", p.Name, p.Args))
	} else {
		rel = shell.RunCmd(fmt.Sprintf("ps -ef | grep ' %s$' | grep -v grep | awk '{print $3}'", p.Name))
	}
	if rel.Ok {
		newPpidStr := strings.TrimSpace(string(rel.Output))
		newPpid := conv.StrToUint32(newPpidStr)
		return newPpid
	}

	return 0
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
		return db.Db().Save(&p).Error == nil
	}

	return false
}

func Enqueue(p *model.Process) error {
	if !IsRunning(p.Pid) && p.Status == model.StatusQueued {
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

func Restart(p *model.Process, isInitializing bool) error {
	if !IsRunning(p.Pid) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed || p.Status == model.StatusClosed) {
		if updatedFromPsCmd(p) {
			return nil
		}

		if !p.AutoRestart {
			if p.Status == model.StatusRunning { // but process is dead, update db state
				observer.HandleEvent(&observer.Event{
					Type:    observer.FailedEvent,
					Process: p,
				})
				p.Status = model.StatusFailed
				db.Db().Save(&p)
			}
			return nil
		}

		if isInitializing {
			pmond.Log.Infof("(re)starting process during initialization: %s", p.Stringify())
		} else {
			observer.HandleEvent(&observer.Event{
				Type:    observer.RestartEvent,
				Process: p,
			})
		}

		_, err := proxyWorker(p, "restart")
		if err != nil {
			return err
		} else {
			if p.Status != model.StatusClosed {
				p.IncrRestartCount()
			}
		}
	}

	return nil
}

func SendOsKillSignal(p *model.Process, status model.ProcessStatus, forced bool) error {
	var cmd *exec.Cmd
	if forced {
		cmd = exec.Command("kill", "-9", p.GetPidStr())
	} else {
		cmd = exec.Command("kill", p.GetPidStr())
	}

	err := cmd.Run()
	if err != nil {
		pmond.Log.Warn(err)
	}

	p.Status = status

	return db.Db().Save(p).Error
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
	execP, err := Exec(p.ProcessFile, p.Log, p.Name, p.Args, p.Username, p.AutoRestart)
	if err != nil {
		return "", err
	}

	execP.ID = p.ID
	execP.CreatedAt = p.CreatedAt
	execP.UpdatedAt = time.Now()

	waitData := NewProcStat(execP).Wait()
	return waitData.Save(db.Db())
}

func workerStart(p *model.Process) (string, error) {
	//returns an instance of the process model
	execP, err := Exec(p.ProcessFile, p.Log, p.Name, p.Args, p.Username, p.AutoRestart)
	if err != nil {
		return "", err
	}

	execP.CreatedAt = time.Now()
	execP.UpdatedAt = time.Now()

	waitData := NewProcStat(execP).Wait()
	return waitData.Save(db.Db())
}
