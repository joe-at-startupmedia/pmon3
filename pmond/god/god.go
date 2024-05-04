package god

import (
	"os"
	"os/signal"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"sync"
	"syscall"
	"time"
)

func New() {
	if pmond.Config.ShouldHandleInterrupts() {
		pmond.Log.Debugf("Capturing interrupts.")
		interruptHandler()
	}

	connectResponder()

	runMonitor()
}

func handleOpenError(e error) {
	if e != nil {
		pmond.Log.Fatal("could not initialize sender: ", e.Error())
	}
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
		pmond.Log.Infof("Captured interrupt: %s", s)
		uninterrupted = false
		time.Sleep(1 * time.Second) //wait for the infinity loop to break
		emptyCmd := protos.Cmd{}
		controller.KillByParams(&emptyCmd, true, model.StatusClosed)
		err := closeResponder()
		if err != nil {
			pmond.Log.Warnf("Error closing queues: %-v", err)
		}
		time.Sleep(1 * time.Second) //wait for responder to close before exiting
		os.Exit(0)
	}()
}

func runMonitor() {

	isInitializing := true

	go func() {
		time.Sleep(30 * time.Second)
		isInitializing = false
	}()

	go monitorResponderStatus(pmond.Log)

	go processRequests(pmond.Log)

	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		runningTask(isInitializing)
	}
}

var pendingTask sync.Map

func runningTask(isInitializing bool) {

	var all []model.Process
	err := db.Db().Find(&all, "status in (?, ?, ?, ?)",
		model.StatusRunning,
		model.StatusFailed,
		model.StatusQueued,
		model.StatusClosed, //closes on pmond handled interrupts
	).Error
	if err != nil {
		return
	}

	for _, p := range all {
		key := "process_id:" + p.GetIdStr()
		_, ok := pendingTask.LoadOrStore(key, p.ID)
		if ok {
			return
		}

		go func(q model.Process, key string) {
			var cur model.Process
			defer func() {
				pendingTask.Delete(key)
			}()
			err = db.Db().First(&cur, q.ID).Error
			if err != nil {
				pmond.Log.Infof("Task monitor could not find process in database: %d", q.ID)
				return
			}

			if cur.Status == model.StatusRunning || cur.Status == model.StatusFailed || cur.Status == model.StatusClosed {
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					pmond.Log.Debugf("Only processes older than 5 seconds can be restarted: %s", q.Stringify())
					return
				}

				err = process.Restart(&cur, isInitializing)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to restart process(%s): %s", q.Stringify(), err)
				}
			} else if cur.Status == model.StatusQueued {
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					pmond.Log.Debugf("Only processes older than 1 second can be enqueued: %s", q.Stringify())
					return
				}

				err = process.Enqueue(&cur)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to enqueue process(%s): %s", q.Stringify(), err)
				}
			}

		}(p, key)
	}
}
