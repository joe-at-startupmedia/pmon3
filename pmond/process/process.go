package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"pmon3/pmond"
	"pmon3/pmond/conf"
	"pmon3/pmond/model"
	"pmon3/pmond/utils/array"
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
		//return that its exist
		return !os.IsNotExist(err)
	}

	return true
}

// used as a last alternative if /proc/[pid]/status and golang isNotExist fail to detect running
func updatedFromPsCmd(p *model.Process) bool {
	// try to get process new pid
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

func EnqueueProcess(p *model.Process) error {
	if !IsRunning(p.Pid) && p.Status == model.StatusQueued {
		if updatedFromPsCmd(p) {
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

		_, err := tryRun(p, "", "restart")
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

func GetProcUser(a *model.ExecFlags) (*user.User, error) {
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

	data, err := runProcess([]string{cmd, m.ProcessFile, flagsModel.Json()})
	if err != nil {
		return nil, err
	}

	var tb []string
	_ = json.Unmarshal(data, &tb)

	return tb, nil
}

var cmdTypes = []string{"start", "restart"}

func runProcess(args []string) ([]byte, error) {

	for _, arg := range args {
		pmond.Log.Debugf("RunProcess arg: %s\n", arg)
	}

	if len(args) <= 2 {
		return nil, errors.New("process params not valid")
	}
	err := pmond.Instance(conf.GetConfigFile())

	if err != nil {
		return nil, err
	}
	// check run type param
	typeCli := args[0]

	if !array.In(cmdTypes, typeCli) {
		return nil, errors.WithStack(err)
	}

	var output string

	flags := model.ExecFlags{}
	flagModel, err := flags.Parse(args[2])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch typeCli {
	case "start":
		output, err = workerStart(args[1], flagModel)
	case "restart":
		output, err = workerRestart(args[1], flagModel)
	}

	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte(output), nil
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
	runUser, err := GetProcUser(flags)
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
	runUser, err := GetProcUser(flags)
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
