package god

import (
	"os"
	"os/signal"
	"pmon3/cli/cmd/kill"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"pmon3/pmond/pmq"
	"pmon3/pmond/svc/process"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func NewMonitor() {
	if pmond.Config.ShouldHandleInterrupts() {
		pmond.Log.Debugf("Capturing interrupts.")
		interruptHandler()
	}
	pmq.New()
	runMonitor()
}

var uninterrupted bool = true

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
		uninterrupted = false
		//wait for the inifity loop to break
		time.Sleep(1 * time.Second)
		kill.Kill(model.StatusClosed)
		pmq.Close()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func runMonitor() {

	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		runningTask()
		pmq.HandleRequest()
		if !uninterrupted {
			break
		}
	}
	//wait for the interrupt handler to complete
	time.Sleep(5 * time.Second)
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

	for _, p := range all {
		key := "process_id:" + strconv.Itoa(int(p.ID))
		_, ok := pendingTask.LoadOrStore(key, p.ID)
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
				//only processes older than 5 seconds can be restarted
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					return
				}

				err = process.RestartProcess(&p)
				if err != nil {
					pmond.Log.Error(err)
				}
			} else if cur.Status == model.StatusQueued {
				//only processes older than 1 seconds can be enqueued
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					return
				}

				err = process.EnqueueProcess(&p)
				if err != nil {
					pmond.Log.Error(err)
				}
			}

		}(p, key)
	}
}
