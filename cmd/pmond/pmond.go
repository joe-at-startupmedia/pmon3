package main

import (
	"context"
	"os"
	"os/signal"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/god"
	"pmon3/pmond/shell"
	"sync"
	"syscall"
)

func main() {

	err := pmond.Instance(conf.GetConfigFile(), conf.GetProcessConfigFile())
	if err != nil {
		pmond.Log.Fatal(err)
	}

	if shell.ExecIsPmondRunning(os.Getpid()) {
		pmond.Log.Fatal("pmond is already running")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	//viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWesteros), viewer.WithLinkAddr("goprofiler.test:8080"))
	//mgr := statsview.New()
	//go mgr.Start()

	ctx := interruptHandler(&wg)
	god.Summon(ctx)
	wg.Wait() //wait for the interrupt handler to complete

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
		god.Banish()
		wg.Done()
	}()

	return ctx
}
