package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/proxy"
	"pmon3/pmond/utils/conv"
	"strings"

	"github.com/goinbox/shell"
)

func IsRunning(pid uint32) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func TryStop(p *model.Process, status model.ProcessStatus, forced bool) error {
	var cmd *exec.Cmd
	if forced {
		cmd = exec.Command("kill", "-9", p.GetPidStr())
	} else {
		cmd = exec.Command("kill", p.GetPidStr())
	}

	err := cmd.Run()
	if err != nil {
		pmond.Log.Fatal(err)
	}

	p.Status = status

	return pmond.Db().Save(p).Error
}

func EnqueueProcess(p *model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
		//fmt.Printf("Monitor: process (%d) already running \n", p.Pid)
		return nil
	}

	if os.IsNotExist(err) && p.Status == model.StatusQueued {
		if checkFork(p) {
			return nil
		}

		_, err := tryRun(p, "", "start")
		if err != nil {
			return err
		}
	}

	return nil
}

func RestartProcess(p *model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
		//fmt.Printf("Monitor: process (%d) already running \n", p.Pid)
		return nil
	}

	// proc status file not exit
	if os.IsNotExist(err) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed || p.Status == model.StatusClosed) {
		if checkFork(p) {
			return nil
		}

		// check whether set auto restart
		if !p.AutoRestart {
			if p.Status == model.StatusRunning { // but process is dead, update db state
				p.Status = model.StatusFailed
				pmond.Db().Save(&p)
			}
			return nil
		}

		_, err := tryRun(p, "", "restart")
		if err != nil {
			return err
		}
	}

	return nil
}

func tryRun(m *model.Process, flags string, cmd string) ([]string, error) {
	var flagsModel = model.ExecFlags{
		User:          m.Username,
		Log:           m.Log,
		NoAutoRestart: !m.AutoRestart,
		Args:          m.Args,
		Name:          m.Name,
	}

	pmond.Log.Infof("%sing process: %s %s\n", cmd, m.Name, m.ProcessFile)
	if len(flags) > 0 {
		pmond.Log.Debugf("with flags: %s \n", flags)
		execFlags := model.ExecFlags{}
		curFlag, err := execFlags.Parse(flags)
		if err != nil {
			return nil, err
		}

		if len(curFlag.Log) > 0 {
			flagsModel.Log = curFlag.Log
		}

		// log dir
		if len(curFlag.LogDir) > 0 && len(curFlag.Log) == 0 {
			flagsModel.LogDir = curFlag.LogDir
			flagsModel.Log = ""
		}
	}

	data, err := proxy.RunProcess([]string{cmd, m.ProcessFile, flagsModel.Json()})
	if err != nil {
		return nil, err
	}

	var tb []string
	_ = json.Unmarshal(data, &tb)

	return tb, nil
}

// Detects whether a new process is created
// @TODO this probably wont work when process contain substrings of other process names
func checkFork(process *model.Process) bool {
	// try to get process new pid
	rel := shell.RunCmd(fmt.Sprintf("ps -ef | grep '%s ' | grep -v grep | awk '{print $2}'", process.Name))
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		if newPid != 0 && newPid != process.Pid {
			process.Pid = newPid
			process.Status = model.StatusRunning
			return pmond.Db().Save(&process).Error == nil
		}
	}

	return false
}
