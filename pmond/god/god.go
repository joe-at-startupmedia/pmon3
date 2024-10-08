package god

import (
	"context"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"pmon3/pmond"
	"pmon3/pmond/controller"
	"pmon3/pmond/model"
	"pmon3/pmond/process"
	"pmon3/pmond/protos"
	"pmon3/pmond/repo"
	"sync"
	"syscall"
	"time"
)

var xr xipc.IResponder
var pendingTask sync.Map

func New() {

	var wg sync.WaitGroup
	wg.Add(1)

	//viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWesteros), viewer.WithLinkAddr("goprofiler.test:8080"))
	//mgr := statsview.New()
	//go mgr.Start()

	ctx := interruptHandler(&wg)
	Summon(ctx)
	wg.Wait() //wait for the interrupt handler to complete
}

func Summon(ctx context.Context) {
	connectResponder()
	runMonitor(ctx)
}

func Banish() {
	processMonitorInterval := time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval)
	time.Sleep(2 * processMonitorInterval) //wait for the runMonitor loop to break
	if pmond.Config.HandleInterrupts {
		controller.KillByParams(&protos.Cmd{}, true)
	}
	err := closeResponder()
	if err != nil {
		pmond.Log.Warnf("Error closing queues: %-v", err)
	}
	time.Sleep(3 * processMonitorInterval) //wait for responder to close and requestProcessor to break before exiting
}

func interruptHandler(wg *sync.WaitGroup) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		s := <-sigc
		pmond.Log.Infof("Captured interrupt: %s", s)
		cancel() // terminate the runMonitor loop
		Banish()
		wg.Done()
	}()

	return ctx
}

func runMonitor(ctx context.Context) {

	go processRequests(ctx, pmond.Log)

	controller.StartAppsFromBoth(true)

	initCtx, cancel := context.WithTimeout(ctx, pmond.Config.GetInitializationPeriod())
	defer cancel()

	processRepo := repo.Process()

	timer := time.NewTicker(time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-initCtx.Done():
			runningTask(processRepo, false)
		default:
			runningTask(processRepo, true)
		}
		<-timer.C
	}
}

func runningTask(processRepo *repo.ProcessRepo, isInitializing bool) {

	all, err := processRepo.FindForMonitor()
	if err != nil {
		return
	}

	for _, p := range all {
		_, loaded := pendingTask.LoadOrStore(p.ID, p.ID)
		if loaded { //a goroutine is still working for this pid
			return
		}

		go func(cur *model.Process) {

			defer pendingTask.Delete(cur.ID)

			flapDetector := detectFlapping(cur)

			if cur.Status == model.StatusRunning || cur.Status == model.StatusFailed || cur.Status == model.StatusClosed {
				if time.Since(cur.UpdatedAt).Seconds() <= 5 {
					pmond.Log.Debugf("Only processes older than 5 seconds can be restarted: %s", cur.Stringify())
					return
				}
				restarted, err := process.Restart(cur, isInitializing)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to restart process(%s): %s", cur.Stringify(), err)
				} else if restarted && pmond.Config.FlapDetection.IsEnabled {
					flapDetector.RestartProcess()
				}
			} else if cur.Status == model.StatusQueued {
				if time.Since(cur.UpdatedAt).Seconds() <= 1 {
					pmond.Log.Debugf("Only processes older than 1 second can be enqueued: %s", cur.Stringify())
					return
				}

				err = process.Enqueue(cur, false)
				if err != nil {
					pmond.Log.Errorf("task monitor encountered error attempting to enqueue process(%s): %s", cur.Stringify(), err)
				}
			}

		}(&p)
	}
}

func detectFlapping(p *model.Process) *process.FlapDetector {
	var flapDetector *process.FlapDetector
	if pmond.Config.FlapDetection.IsEnabled {
		flapDetector = process.GetFlapDetectorByProcessId(p.ID, pmond.Config)

		if flapDetector.ShouldBackOff(time.Millisecond * time.Duration(pmond.Config.ProcessMonitorInterval)) {
			if p.Status != model.StatusBackoff {
				repo.ProcessOf(p).UpdateStatus(model.StatusBackoff)
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
				logger.Errorf("Error handling request: %s", err.Error())
				if err.Error() == "buffer is closed" {
					return
				}
			}
		}
	}
}

func closeResponder() error {
	return xr.CloseResponder()
}
