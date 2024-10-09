package process

import (
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/shell"
	"time"
)

type ProcStat struct {
	Done    chan int
	Process *model.Process
}

// NewProcStat kicks off two goroutines:
// the first routine calls os.Process.Wait() which returns when the process exits
// When the process exits it will return the status code
// if the first routine doesn't finish (process did not exit) in 200ms
// the next goroutine will return the result of Process.IsRunning
// which utilizes proc/%d/stat os.Stat and !os.IsNotExist
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
	go r.processExistCheck(r.Process)
}

func (r *ProcStat) processWait(process *model.Process) {
	processState, err := process.Pointer.Wait()
	if err != nil {
		pmond.Log.Warnf("ProcStat processWait error: %v", err)
		r.Done <- 1
		return
	}
	pmond.Log.Infof("ProcStat processWait: %v", processState)
	if processState.Exited() {
		r.Done <- processState.ExitCode()
	}
}

func (r *ProcStat) processExistCheck(p *model.Process) {
	timer := time.NewTicker(time.Millisecond * 200)
	defer timer.Stop()
	for {
		select {
		case existCode := <-r.Done:
			pmond.Log.Warnf("ProcStat exitCode: %d", existCode)
			if existCode != 0 { // process exist exception
				r.Done <- 1
				return
			}
		case <-timer.C: // check process status by proc file
			isRunning := shell.ExecIsRunning(p)
			if isRunning {
				r.Done <- 0
			} else {
				pmond.Log.Warnf("ProcStat timer.C IsRunning: %t", isRunning)
				r.Done <- 1
			}
			return
		}
	}
}
