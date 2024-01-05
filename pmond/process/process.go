package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/conv"
	"strconv"
	"strings"
	"time"

	"github.com/goinbox/shell"
	"github.com/pkg/errors"
)

func IsRunning(pid uint32) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))

	//if it doesnt exist in proc/n/status ask the OS
	if err != nil {
		//it's running if it exists
		return !os.IsNotExist(err)
	}

	return true
}

// used as a last alternative if /proc/[pid]/status and golang isNotExist fail to detect running
func updatedFromPsCmd(p *model.Process) bool {
	rel := shell.RunCmd(fmt.Sprintf("ps -ef | grep ' %s ' | grep -v grep | awk '{print $2}'", p.Name))
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := conv.StrToUint32(newPidStr)
		if newPid != 0 && newPid != p.Pid {
			p.Pid = newPid
			p.Status = model.StatusRunning
			return pmond.Db().Save(&p).Error == nil
		}
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

func Restart(p *model.Process) error {
	if !IsRunning(p.Pid) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed || p.Status == model.StatusClosed) {
		if updatedFromPsCmd(p) {
			return nil
		}

		if !p.AutoRestart {
			if p.Status == model.StatusRunning { // but process is dead, update db state
				p.Status = model.StatusFailed
				pmond.Db().Save(&p)
			}
			return nil
		}

		_, err := proxyWorker(p, "restart")
		if err != nil {
			return err
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
		pmond.Log.Fatal(err)
	}

	p.Status = status

	return pmond.Db().Save(p).Error
}

func SetUser(a *model.ExecFlags) (*user.User, error) {
	runUser := a.User
	var curUser *user.User
	var err error

	if len(runUser) <= 0 {
		curUser, err = user.LookupId(strconv.Itoa(os.Getuid()))
	} else {
		curUser, err = user.Lookup(runUser)
	}

	if err != nil {
		return nil, err
	}

	return curUser, nil
}

func proxyWorker(m *model.Process, cmd string) ([]string, error) {
	var flagsModel = model.ExecFlags{
		User:          m.Username,
		Log:           m.Log,
		NoAutoRestart: !m.AutoRestart,
		Args:          m.Args,
		Name:          m.Name,
	}

	pmond.Log.Infof("%sing process: %s %s\n", cmd, m.Name, m.ProcessFile)

	var (
		tb     []string
		output string
		err    error
	)

	switch cmd {
	case "start":
		output, err = workerStart(m.ProcessFile, &flagsModel)
	case "restart":
		output, err = workerRestart(m.ProcessFile, &flagsModel)
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(output), &tb)
	return tb, nil
}

func workerRestart(pFile string, flags *model.ExecFlags) (string, error) {
	err, m := model.FindProcessByFileAndName(pmond.Db(), pFile, flags.Name)
	if err != nil {
		return "", errors.New("Could not find process")
	}

	cstLog := flags.Log
	if len(cstLog) > 0 && cstLog != m.Log {
		m.Log = cstLog
	}

	// if reset log dir
	if len(flags.LogDir) > 0 && len(flags.Log) == 0 {
		m.Log = ""
	}

	cstName := flags.Name
	if len(cstName) > 0 && cstName != m.Name {
		m.Name = cstName
	}

	extArgs := flags.Args
	if len(extArgs) > 0 {
		m.Args = extArgs
	}

	// get run process user
	runUser, err := SetUser(flags)
	if err != nil {
		return "", err
	}

	//returns an instance of the process model
	p, err := Exec(m.ProcessFile, m.Log, m.Name, m.Args, runUser, !flags.NoAutoRestart, flags.LogDir)
	if err != nil {
		return "", err
	}

	// update process extra data
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = time.Now()

	waitData := NewProcStat(p).Wait()

	return waitData.Save(pmond.Db())
}

func workerStart(processFile string, flags *model.ExecFlags) (string, error) {
	// prepare params
	file, err := os.Stat(processFile)
	if os.IsNotExist(err) || file.IsDir() {
		return "", errors.Errorf("%s not exist", processFile)
	}

	// get run process user
	runUser, err := SetUser(flags)
	if err != nil {
		return "", err
	}

	name := flags.Name
	// get process file name
	if len(name) <= 0 {
		name = filepath.Base(processFile)
	}
	// checkout process name whether exist
	if pmond.Db().First(&model.Process{}, "name = ? AND status != ?", name, model.StatusQueued).Error == nil {
		return "", errors.Errorf("process name: %s already exist, please set other name by --name", name)
	}
	// start process
	process, err := Exec(processFile, flags.Log, name, flags.Args, runUser, !flags.NoAutoRestart, flags.LogDir)
	if err != nil {
		return "", err
	}
	process.CreatedAt = time.Now()
	process.UpdatedAt = time.Now()
	// waiting process state
	stat := NewProcStat(process).Wait()
	// return process data
	return stat.Save(pmond.Db())
}
