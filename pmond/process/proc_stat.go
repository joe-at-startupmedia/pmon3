package process

import (
	"fmt"
	"os"
	"pmon3/pmond/model"
	"time"
)

type ProcStat struct {
	Done    chan int
	Process *model.Process
}

func NewProcStat(p *model.Process) *ProcStat {
	stat := &ProcStat{
		Done:    make(chan int),
		Process: p,
	}

	stat.run()

	return stat
}

func (r *ProcStat) Wait() *model.Process {
	// waiting result
	status := <-r.Done

	if status > 0 {
		r.Process.Status = model.StatusFailed
	} else {
		r.Process.Status = model.StatusRunning
	}

	return r.Process
}

func (r *ProcStat) run() {
	go r.processWait(r.Process)
	go r.processExistCheck(int(r.Process.Pid))
}

func (r *ProcStat) processWait(process *model.Process) {
	processState, err := process.Pointer.Wait()
	if err != nil {
		r.Done <- 1
		return
	}

	if processState.Exited() {
		r.Done <- processState.ExitCode()
	}
}

func (r *ProcStat) processExistCheck(pid int) {
	timer := time.NewTicker(time.Millisecond * 200)
	defer timer.Stop()
	for {
		select {
		case existCode := <-r.Done:
			if existCode != 0 { // process exist exception
				r.Done <- 1
				return
			}
		case <-timer.C: // check process status by proc file
			_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))
			if !os.IsNotExist(err) {
				r.Done <- 0
				return
			}
		}
	}
}
