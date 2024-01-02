package god

import (
	"fmt"
	"os"
	"os/signal"
	"pmon3/cli/cmd/kill"
	"pmon3/pmond"
	"pmon3/pmond/model"
	process2 "pmon3/pmond/svc/process"
	"pmon3/pmond/utils/iconv"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/goinbox/shell"
)

func NewMonitor() {
	if pmond.Config.ShouldHandleInterrupts() {
		pmond.Log.Debugf("Capturing interrupts.")
		interruptHandler()
	}
	runMonitor()
}

func interruptHandler() {
	sigc := make(chan os.Signal)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		pmond.Log.Debugf("Captured interrupt: %s \n", s)
		kill.Kill(model.StatusClosed)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
}

func runMonitor() {
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		runningTask()
	}
}

var pendingTask sync.Map

func runningTask() {
	var all []model.Process
	err := pmond.Db().Find(&all, "status in (?, ?, ?, ?)",
		model.StatusRunning,
		model.StatusFailed,
		model.StatusQueued,
		model.StatusClosed, //closes on pmond handled interrupts
	).Error
	if err != nil {
		return
	}

	for _, process := range all {
		// just check failed process
		key := "process_id:" + strconv.Itoa(int(process.ID))
		_, ok := pendingTask.LoadOrStore(key, process.ID)
		if ok {
			return
		}

		go func(p model.Process, key string) {
			var cur model.Process
			defer func() {
				pendingTask.Delete(key)
			}()
			err = pmond.Db().First(&cur, p.ID).Error
			if err != nil {
				return
			}

			if cur.Status == model.StatusRunning || cur.Status == model.StatusFailed || cur.Status == model.StatusClosed {
				//only process older than 5 seconds can be restarted
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					return
				}

				err = restartProcess(p)
				if err != nil {
					pmond.Log.Error(err)
				}
			} else if cur.Status == model.StatusQueued {
				//only process older than 1 seconds can be enqueued
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					return
				}

				err = enqueueProcess(p)
				if err != nil {
					pmond.Log.Error(err)
				}
			}

		}(process, key)
	}
}

// Detects whether a new process is created
// @TODO this probably wont work when process contain substrings of other process names
func checkFork(process model.Process) bool {
	// try to get process new pid
	rel := shell.RunCmd(fmt.Sprintf("ps -ef | grep '%s ' | grep -v grep | awk '{print $2}'", process.Name))
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := iconv.MustInt(newPidStr)
		if newPid != 0 && newPid != process.Pid {
			process.Pid = newPid
			process.Status = model.StatusRunning
			return pmond.Db().Save(&process).Error == nil
		}
	}

	return false
}

func enqueueProcess(p model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
		//fmt.Printf("Monitor: process (%d) already running \n", p.Pid)
		return nil
	}

	if os.IsNotExist(err) && p.Status == model.StatusQueued {
		if checkFork(p) {
			return nil
		}

		_, err := process2.TryStart(p, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func restartProcess(p model.Process) error {
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

		_, err := process2.TryRestart(p, "")
		if err != nil {
			return err
		}
	}

	return nil
}
