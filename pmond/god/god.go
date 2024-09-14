package god

import (
	"context"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/sirupsen/logrus"
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

var xr xipc.IResponder

func New() {

	var wg sync.WaitGroup
	wg.Add(1)
	ctx := interruptHandler(pmond.Config.HandleInterrupts, &wg)
	connectResponder()
	runMonitor(ctx)
	wg.Wait() //wait for the interrupt handler to complete
}

func interruptHandler(shouldCloseOnInterrupt bool, wg *sync.WaitGroup) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	processMonitorInterval := time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval)

	go func() {
		s := <-sigc
		pmond.Log.Infof("Captured interrupt: %s, should close(%t)", s, shouldCloseOnInterrupt)
		cancel()                               // terminate the runMonitor loop
		time.Sleep(2 * processMonitorInterval) //wait for the runMonitor loop to break
		if shouldCloseOnInterrupt {
			controller.KillByParams(&protos.Cmd{}, true, model.StatusClosed)
		}
		err := closeResponder()
		if err != nil {
			pmond.Log.Warnf("Error closing queues: %-v", err)
		}
		time.Sleep(3 * processMonitorInterval) //wait for responder to close and requestProcessor to break before exiting
		wg.Done()
	}()

	return ctx
}

func runMonitor(ctx context.Context) {

	go processRequests(ctx, pmond.Log)

	controller.StartAppsFromBoth(true)

	initCtx, cancel := context.WithTimeout(ctx, pmond.Config.GetInitializationPeriod())
	defer cancel()

	timer := time.NewTicker(time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-initCtx.Done():
			runningTask(false)
		default:
			runningTask(true)
		}
		<-timer.C
	}
}

var pendingTask sync.Map

func runningTask(isInitializing bool) {

	var all []model.Process
	err := db.Db().Find(&all, "status in (?, ?, ?, ?, ?)",
		model.StatusRunning,
		model.StatusFailed,
		model.StatusQueued,
		model.StatusClosed,
		model.StatusBackoff,
	).Error
	if err != nil {
		return
	}

	for _, p := range all {
		key := p.GetIdStr()
		_, loaded := pendingTask.LoadOrStore(key, p.ID)
		if loaded { //a goroutine is still working for this pid
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

			flapDetector := detectFlapping(&cur)

			if cur.Status == model.StatusRunning || cur.Status == model.StatusFailed || cur.Status == model.StatusClosed {
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					pmond.Log.Debugf("Only processes older than 5 seconds can be restarted: %s", q.Stringify())
					return
				}
				restarted, err := process.Restart(&cur, isInitializing)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to restart process(%s): %s", q.Stringify(), err)
				} else if restarted && pmond.Config.FlapDetectionEnabled {
					flapDetector.RestartProcess()
				}
			} else if cur.Status == model.StatusQueued {
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					pmond.Log.Debugf("Only processes older than 1 second can be enqueued: %s", q.Stringify())
					return
				}

				err = process.Enqueue(&cur, false)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to enqueue process(%s): %s", q.Stringify(), err)
				}
			}

		}(p, key)
	}
}

func detectFlapping(p *model.Process) *process.FlapDetector {
	var flapDetector *process.FlapDetector
	if pmond.Config.FlapDetectionEnabled {
		flapDetector = process.GetFlapDetectorByProcessId(p.ID, pmond.Config)

		if flapDetector.ShouldBackOff(time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval)) {
			if p.Status != model.StatusBackoff {
				p.UpdateStatus(db.Db(), model.StatusBackoff)
			}
		} else if p.Status == model.StatusBackoff {
			//set it back to failed so process evaluation can resume
			p.Status = model.StatusFailed
		}
	}
	return flapDetector
}

func handleOpenError(e error) {
	if e != nil {
		pmond.Log.Fatal("could not initialize sender: ", e.Error())
	}
}

// HandleCmdRequest provides a concrete implementation of HandleRequestFromProto using the local Cmd protobuf type
func handleCmdRequest(mqr xipc.IResponder) error {
	cmd := &protos.Cmd{}
	return mqr.HandleRequestFromProto(cmd, func() (processed []byte, err error) {
		return controller.MsgHandler(cmd)
	})
}

func processRequests(ctx context.Context, logger *logrus.Logger) {
	for {
		err := handleCmdRequest(xr) //blocking
		select {
		case <-ctx.Done():
			logger.Infof("request process is terminated.")
			return
		default:
			if err != nil {
				logger.Errorf("Error handling request: %-v", err)
			}
		}
	}
}

func closeResponder() error {
	return xr.CloseResponder()
}
