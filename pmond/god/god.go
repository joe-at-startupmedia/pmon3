package god

import (
	"github.com/joe-at-startupmedia/pmq_responder"
	"os"
	"os/signal"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"sync"
	"syscall"
	"time"
)

var pmr *pmq_responder.MqResponder

func New() {
	if pmond.Config.ShouldHandleInterrupts() {
		pmond.Log.Debugf("Capturing interrupts.")
		interruptHandler()
	}

	queueConfig := pmq_responder.QueueConfig{
		Name:  "pmon3_mq",
		Dir:   pmond.Config.GetPosixMessageQueueDir(),
		Flags: pmq_responder.O_RDWR | pmq_responder.O_CREAT | pmq_responder.O_NONBLOCK,
	}
	ownership := pmq_responder.Ownership{
		Group:    pmond.Config.PosixMessageQueueGroup,
		Username: pmond.Config.PosixMessageQueueUser,
	}
	pmqResponder := pmq_responder.NewResponder(&queueConfig, &ownership)
	if pmqResponder.HasErrors() {
		handleOpenError(pmqResponder.ErrRqst)
		handleOpenError(pmqResponder.ErrResp)
	}
	pmr = pmqResponder
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
		pmond.Log.Debugf("Captured interrupt: %s \n", s)
		uninterrupted = false
		//wait for the inifity loop to break
		time.Sleep(1 * time.Second)
		emptyCmd := protos.Cmd{}
		controller.KillByParams(&emptyCmd, true, model.StatusClosed)
		err := pmr.UnlinkResponder()
		if err != nil {
			pmond.Log.Warnf("Error closing queues: %-v", err)
		}
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func runMonitor() {

	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		runningTask()
		err := controller.HandleCmdRequest(pmr)
		if err != nil {
			pmond.Log.Warnf("Error handling request: %-v", err)
		}
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
			err = pmond.Db().First(&cur, q.ID).Error
			if err != nil {
				return
			}

			if cur.Status == model.StatusRunning || cur.Status == model.StatusFailed || cur.Status == model.StatusClosed {
				//only processes older than 5 seconds can be restarted
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					return
				}

				err = process.Restart(&cur)
				if err != nil {
					pmond.Log.Error(err)
				}
			} else if cur.Status == model.StatusQueued {
				//only processes older than 1 seconds can be enqueued
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					return
				}

				err = process.Enqueue(&cur)
				if err != nil {
					pmond.Log.Error(err)
				}
			}

		}(p, key)
	}
}
